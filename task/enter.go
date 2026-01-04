package task

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"server/global"
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

	// 定时任务：获取热搜信息
	if _, err := c.AddFunc("@hourly", func() {
		if err2 := GetHotListSyncTask(); err2 != nil {
			global.Log.Error("Failed to get hot list", zap.Error(err2))
		}
	}); err != nil {
		return err
	}

	// 定时任务：获取当天的日历信息
	if _, err := c.AddFunc("@daily", func() {
		if err2 := GetCalendarSyncTask(); err2 != nil {
			global.Log.Error("Failed to get calendar", zap.Error(err2))
		}
	}); err != nil {
		return err
	}
	return nil
}
