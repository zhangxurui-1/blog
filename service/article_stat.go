package service

import (
	"server/global"
	"strconv"
)

// CountDB 用于浏览量计数
type CountDB struct {
	Index string
}

func (articleService *ArticleService) NewArticleView() CountDB {
	return CountDB{
		Index: "article_view",
	}
}

// Set 将 Redis 缓存中的文章浏览量 +1
func (c CountDB) Set(id string) error {
	num, _ := global.Redis.HGet(c.Index, id).Int()
	num++
	return global.Redis.HSet(c.Index, id, num).Err()
}

// GetInfo 获取所有文章的浏览量
func (c CountDB) GetInfo() map[string]int {
	var Info = map[string]int{}
	maps := global.Redis.HGetAll(c.Index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		Info[id] = num
	}
	return Info
}

// Clear 清除数据
func (c CountDB) Clear() {
	global.Redis.Del(c.Index)
}

func (c CountDB) Delete(id string) error {
	return global.Redis.HDel(c.Index, id).Err()
}
