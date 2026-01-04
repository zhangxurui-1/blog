package initialize

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"os"
	"server/global"
)

// ConnectRedis 初始化 redis 连接
func ConnectRedis() redis.Client {
	redisCfg := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	// 测试连通性
	if _, err := client.Ping().Result(); err != nil {
		global.Log.Error("Fail to connect redis", zap.Error(err))
		os.Exit(1)
	}

	return *client
}
