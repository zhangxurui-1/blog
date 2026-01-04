package service

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/utils"
)

// JwtService 提供与 JWT 相关的服务
type JwtService struct {
}

// SetRedisJWT 将 JWT 存储到 Redis
func (jwtService *JwtService) SetRedisJWT(jwt string, uuid uuid.UUID) error {
	// 解析 jwt 中的过期时间
	dr, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		return err
	}
	return global.Redis.Set(uuid.String(), jwt, dr).Err()
}

// GetRedisJWT 从 redis 获取 jwt
func (jwtService *JwtService) GetRedisJWT(uuid uuid.UUID) (string, error) {
	return global.Redis.Get(uuid.String()).Result()
}

// JoinInBlackList 将 jwt 添加到黑名单
func (jwtService *JwtService) JoinInBlackList(jwtList database.JwtBlacklist) error {
	// 将 jwt 记录插入到数据库
	if err := global.DB.Create(&jwtList).Error; err != nil {
		return err
	}
	// 将 jwt 添加到内存中的黑名单缓存
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return nil
}

// IsInBlackList 检查 jwt 是否在黑名单
func (jwtService *JwtService) IsInBlackList(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
}

// LoadAll 加载数据库中的 jwt 黑名单到内存缓存
func LoadAll() {
	var data []string
	// 读取数据库，Pluck 方法用于从数据库中提取单个列的值，并将其存储到切片（slice）中
	if err := global.DB.Model(&database.JwtBlacklist{}).Pluck("jwt", &data).Error; err != nil {
		global.Log.Error("Failed to load JWT blacklist from database", zap.Error(err))
		return
	}
	// 将所有 jwt 添加到 BlackCache
	for _, d := range data {
		global.BlackCache.SetDefault(d, struct{}{})
	}
}
