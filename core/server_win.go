// go build windows

package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 初始化一个标准的 http 服务器
func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,          // 服务器监听地址
		Handler:        router,           // 设置路由
		ReadTimeout:    10 * time.Minute, // 请求读取超时时间
		WriteTimeout:   10 * time.Minute, // 响应写入超时时间
		MaxHeaderBytes: 1 << 20,          // 最大请求头的大小
	}
}
