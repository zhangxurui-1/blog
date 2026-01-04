package service

import (
	"encoding/json"
	"server/global"
	"server/model/other"
	"server/utils/hotSearch"
	"time"
)

type HotSearchService struct {
}

// GetHotSearchDataBySource 获取热搜数据
func (hostSearchService *HotSearchService) GetHotSearchDataBySource(sourceStr string) (other.HotSearchData, error) {
	// 优先从缓存拿
	result, err := global.Redis.Get(sourceStr).Result()
	if err != nil {
		// 缓存拿不到则重新请求
		source := hotSearch.NewSource(sourceStr)
		hotSearchData, err := source.GetHotSearchData(30)
		if err != nil {
			return other.HotSearchData{}, err
		}

		// 将新数据存入 redis
		bytes, err := json.Marshal(hotSearchData)
		if err != nil {
			return other.HotSearchData{}, err
		}
		if err := global.Redis.Set(sourceStr, bytes, time.Hour).Err(); err != nil {
			return other.HotSearchData{}, err
		}
		return hotSearchData, nil
	}
	// 若 redis 中有数据，则取出并反序列化
	var hotSearchData other.HotSearchData
	if err := json.Unmarshal([]byte(result), &hotSearchData); err != nil {
		return other.HotSearchData{}, err
	}
	return hotSearchData, nil
}
