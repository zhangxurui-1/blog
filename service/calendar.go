package service

import (
	"encoding/json"
	"server/global"
	"server/model/other"
	"server/utils"
	"time"
)

type CalendarService struct {
}

// GetCalendarByDate 获取日历
func (calendarService *CalendarService) GetCalendarByDate(dateStr string) (other.Calendar, error) {
	// 优先从 redis 缓存中取
	result, err := global.Redis.Get("calendar-" + dateStr).Result()
	if err != nil {
		// 如果缓存里没有，则通过 http 请求获取第三方数据
		calendar, err := utils.GetCalendar(dateStr)
		if err != nil {
			return other.Calendar{}, err
		}
		// 序列化后存入 redis
		data, err := json.Marshal(calendar)
		if err != nil {
			return other.Calendar{}, err
		}
		if err := global.Redis.Set("calendar-"+dateStr, data, time.Hour*24).Err(); err != nil {
			return other.Calendar{}, err
		}
		return calendar, nil
	}
	// 如果 redis 中有，则反序列化后返回
	var calendar other.Calendar
	if err := json.Unmarshal([]byte(result), &calendar); err != nil {
		return other.Calendar{}, err
	}
	return calendar, nil
}
