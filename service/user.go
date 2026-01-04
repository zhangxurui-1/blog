package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"
)

type UserService struct{}

// Register 用户注册
func (userService *UserService) Register(u database.User) (database.User, error) {
	// 验证数据库中是否已经有这条记录
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in")
	}
	// 填充对象
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "/image/avatar.jpg"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email

	// 添加到数据库
	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}
	return u, nil
}

// EmailLogin 用户邮箱登录
func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	// 如果在数据库找到该用户
	if err == nil {
		if !utils.BcryptCheck(u.Password, user.Password) {
			return database.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return database.User{}, err
}

// QQLogin 用户 QQ 登录
func (userService *UserService) QQLogin(accessTokenResponse other.AccessTokenResponse) (database.User, error) {
	var user database.User

	// 尝试查找用户
	err := global.DB.Where("openid = ?", accessTokenResponse.OpenId).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.User{}, err
	}

	// 如果用户不存在，则创建新用户
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userInfoResponse, err2 := ServiceGroupApp.QQService.GetUserInfoByAccessTokenAndOpenid(
			accessTokenResponse.AccessToken, accessTokenResponse.OpenId)
		if err2 != nil {
			return database.User{}, err2
		}

		// 填充新 user 对象
		user.UUID = uuid.Must(uuid.NewV4())
		user.Username = userInfoResponse.Nickname
		user.Openid = accessTokenResponse.OpenId
		user.Avatar = userInfoResponse.FigureurlQQ2
		user.RoleID = appTypes.User
		user.Register = appTypes.QQ
		// 向数据库插入新记录
		if err := global.DB.Create(&user).Error; err != nil {
			return database.User{}, err
		}
	}
	return user, nil
}

// Logout 用户登出
func (userService *UserService) Logout(c *gin.Context) {
	// jwt 中间件已经把 access token 的 claims 放入上下文
	uuid := utils.GetUUID(c)

	jwtStr := utils.GetRefreshToken(c)
	// 让用户端清除 refresh token
	utils.ClearRefreshToken(c)
	global.Redis.Del(uuid.String())
	// 把 refresh token 加入黑名单
	_ = ServiceGroupApp.JwtService.JoinInBlackList(database.JwtBlacklist{Jwt: jwtStr})
}

// ForgotPassword 用户忘记密码
func (userService *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

// UserCard 查找用户卡片
func (userService *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid = ?", req.UUID).
		Select("uuid", "username", "avatar", "address", "signature").First(&user).Error; err != nil {

		return response.UserCard{}, err
	}

	return response.UserCard{
		UUID:      user.UUID,
		UserName:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

// UserResetPassword 用户重置密码
func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	// 检查原密码是否正确
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("original password does not match the current account")
	}
	// 更新密码
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

// UserInfo 获取用户信息
func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}

// UserChangeInfo 用户修改信息
func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error
}

// UserWeather 用户获取天气
func (userService *UserService) UserWeather(ip string) (string, error) {
	// 尝试从 redis 中获取
	result, err := global.Redis.Get("weather-" + ip).Result()
	// redis 中没有找到
	if err != nil {
		// 获取地理位置
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err2 := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)
		if err2 != nil {
			return "", err2
		}

		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature +
			"°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"

		// 将天气存入 redis
		if err := global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	// 如果 redis 找到了，则直接返回
	return result, nil
}
func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	// 构建查询条件
	where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))

	var res response.UserChart
	// 构建日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	// 获取登录数据
	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	// 获取注册数据
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)

	for _, date := range res.DateList {
		loginCount := loginCounts[date]
		registerCount := registerCounts[date]
		res.LoginData = append(res.LoginData, loginCount)
		res.RegisterData = append(res.RegisterData, registerCount)
	}
	return res, nil
}

// UserList 获取用户列表
func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB
	if info.UUID != nil {
		db = db.Where("uuid = ?", info.UUID)
	}
	if info.Register != nil {
		db = db.Where("register = ?", info.Register)
	}

	// 配置选项
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	// 分页查询
	return utils.MySQLPagination(&database.User{}, option)
}

// UserFreeze 冻结用户
func (userService *UserService) UserFreeze(req request.UserOperation) error {
	var user database.User
	if err := global.DB.Take(&user, req.ID).Update("freeze", true).Error; err != nil {
		return err
	}

	// 将 jwt 加入黑名单
	jwtStr, _ := ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)
	if jwtStr != "" {
		_ = ServiceGroupApp.JwtService.JoinInBlackList(database.JwtBlacklist{Jwt: jwtStr})
	}
	return nil
}

// UserUnfreeze 解冻用户
func (userService *UserService) UserUnfreeze(req request.UserOperation) error {
	return global.DB.Take(&database.User{}, req.ID).Update("freeze", false).Error
}

// UserLoginList 获取用户登录日志
func (userService *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		var userID uint
		if err := global.DB.Model(&database.User{}).Where("uuid = ?", *info.UUID).Pluck("id", &userID).Error; err != nil {
			return nil, 0, nil
		}

		db = db.Where("user_id = ?", userID)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
		PreLoad:  []string{"User"},
	}

	return utils.MySQLPagination(&database.Login{}, option)
}
