package initialize

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"os"
	"server/global"
	"server/task"
)

// ZapLogger 实现 cron.Logger 接口
// Info 和 Error 方法用于接收 cron 包生成的日志并使用 zap 进行记录
type ZapLogger struct {
	logger *zap.Logger
}

func (z *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	z.logger.Info(msg, zap.Any("keysAndValues", keysAndValues))
}

func (z *ZapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	z.logger.Error(msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}

func NewZapLogger() *ZapLogger {
	return &ZapLogger{logger: global.Log}
}

// InitCron 初始化 使用 cron 完成定时任务
func InitCron() {
	// 将 cron 包的日志记录转发到 zap 日志库中，实现统一的日志管理和记录
	c := cron.New(cron.WithLogger(NewZapLogger()))
	if err := task.RegisterScheduledTasks(c); err != nil {
		global.Log.Error("Error scheduling cron task", zap.Error(err))
		os.Exit(1)
	}
	c.Start()
}
