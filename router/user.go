package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
	"server/middleware"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, adminGroup *gin.RouterGroup) {
	// privateGroup 使用 JWTAuth 中间件，adminGroup 使用 AdminAuth 中间件
	userRouter := privateGroup.Group("user")
	userPublicRouter := publicGroup.Group("user")
	userLoginRouter := publicGroup.Group("user").Use(middleware.LoginRecord())
	userAdminRouter := adminGroup.Group("user")
	userApi := api.ApiGroupApp.UserApi

	// userRouter 用户登录状态下的路由
	{
		userRouter.POST("logout", userApi.Logout)
		userRouter.PUT("resetPassword", userApi.UserResetPassword)
		userRouter.GET("info", userApi.UserInfo)
		userRouter.PUT("changeInfo", userApi.UserChangeInfo)
		userRouter.GET("weather", userApi.UserWeather)
		userRouter.GET("chart", userApi.UserChart)
	}
	// userPublicRouter 游客状态下的路由
	{
		userPublicRouter.POST("forgotPassword", userApi.ForgotPassword)
		userPublicRouter.GET("card", userApi.UserCard)
	}
	// 针对登录注册的路由
	{
		userLoginRouter.POST("register", userApi.Register)
		userLoginRouter.POST("login", userApi.Login)
	}
	// 管理员路由
	{
		userAdminRouter.GET("list", userApi.UserList)
		userAdminRouter.PUT("freeze", userApi.UserFreeze)
		userAdminRouter.PUT("unfreeze", userApi.UserUnfreeze)
		userAdminRouter.GET("loginList", userApi.UserLoginList)
	}
}
