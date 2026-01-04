package main

import (
	"server/core"
	"server/flag"
	"server/global"
	"server/initialize"
)

func main() {
	// 加载 yaml 配置文件
	global.Config = core.InitConf()
	// 初始化 Logger
	global.Log = core.InitLogger()
	// 初始化其他配置（本地黑名单缓存的过期时间）
	initialize.OtherInit()

	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	defer global.Redis.Close()
	global.ESClient = initialize.ConnectES()
	flag.InitFlag()
	initialize.InitCron()

	core.RunServer()
}
