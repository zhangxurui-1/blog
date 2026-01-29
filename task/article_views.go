package task

import (
	"context"
	"server/global"
	"server/model/elasticsearch"
	"server/service"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
)

// UpdateArticleViewsSyncTask 同步文章浏览量到 Elasticsearch
func UpdateArticleViewsSyncTask() error {
	// 从 redis 获取浏览量
	articleView := service.ServiceGroupApp.ArticleService.NewArticleView()
	viewsInfo := articleView.GetInfo()

	// 获取文章 id 和 浏览量的增量
	for id, num := range viewsInfo {
		// num 是增量数据
		if num == 0 {
			continue
		}

		// 构造 Es 的 Painless 脚本 (Es 内置的安全脚本语言)，用于原子更新 views 字段
		// 让 views 字段 累加 num
		source := "ctx._source.views += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), id).
			Script(&script).Do(context.TODO())

		if err != nil {
			return err
		}
	}

	articleView.Clear()
	return nil
}
