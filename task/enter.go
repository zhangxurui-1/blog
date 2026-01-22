package task

import (
	"server/global"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// RegisterScheduledTasks 定时任务
func RegisterScheduledTasks(c *cron.Cron) error {
	// 定时任务：同步文章的浏览量到 es
	if _, err := c.AddFunc("@hourly", func() {
		if err := UpdateArticleViewsSyncTask(); err != nil {
			global.Log.Error("Failed to update article views", zap.Error(err))
		}
	}); err != nil {
		return err
	}

	return nil
}
