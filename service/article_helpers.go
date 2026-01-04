package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"gorm.io/gorm"
	"server/global"
	"server/model/database"
	"server/model/elasticsearch"
	"server/utils"
)

// Create 创建文章
// 这段代码等效于：
//
// POST /article_index/_doc/?refresh=true
// Content-Type: application/json
//
//	{
//		"title": "xxx",
//		"content": "xxx",
//		"author": "xxx"
//	}
func (articleService *ArticleService) Create(a *elasticsearch.Article) error {
	_, err := global.ESClient.
		Index(elasticsearch.ArticleIndex()). // 指定索引
		Request(a).                          // 指定要插入的文档内容
		Refresh(refresh.True).               // 强制刷新索引
		Do(context.TODO())                   // 执行请求
	return err
}

// Delete 删除文章
// 这段代码等效于：
//
// POST /article_index/_bulk?refresh=true
// Content-Type: application/json
//
// { "delete": { "_id": "xxx" } }
// { "delete": { "_id": "xxx" } }
// { "delete": { "_id": "xxx" } }
func (articleService *ArticleService) Delete(ids []string) error {
	// 构造批量操作请求
	var request bulk.Request
	for _, id := range ids {
		request = append(request, types.OperationContainer{Delete: &types.DeleteOperation{Id_: &id}})
	}

	_, err := global.ESClient.Bulk().Request(&request).Index(elasticsearch.ArticleIndex()).Refresh(refresh.True).Do(context.TODO())
	return err
}

// Get 根据 id 查询 es 文档
func (articleService *ArticleService) Get(id string) (elasticsearch.Article, error) {
	var a elasticsearch.Article
	// 到 es 中查询
	res, err := global.ESClient.Get(elasticsearch.ArticleIndex(), id).Do(context.TODO())
	if err != nil {
		return elasticsearch.Article{}, err
	}
	// 未查询到
	if !res.Found {
		return elasticsearch.Article{}, errors.New("document not found")
	}

	// 反序列化
	err = json.Unmarshal(res.Source_, &a)
	return a, err
}

// Update 更新文档
// 相当于：
//
// POST /article_index/_update/{articleID}?refresh=true
// Content-Type: application/json
//
//	{
//		"doc": {
//		"title": "xxx",
//		"views": xxx
//		}
//	}
func (articleService *ArticleService) Update(articleID string, v any) error {
	// 序列化
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// 执行更新请求
	_, err = global.ESClient.Update(elasticsearch.ArticleIndex(), articleID).
		Request(&update.Request{Doc: bytes}).
		Refresh(refresh.True).
		Do(context.TODO())
	return err
}

// Exists 根据标题检查文章是否存在
// 这段代码等效于：
// POST /article_index/_search
// Content-Type: application/json
//
//	{
//		"query": {
//			"match": {
//				"keyword": {
//					"query": "xxx"
//				}
//			}
//		},
//		"size": 1
//	}
func (articleService *ArticleService) Exists(title string) (bool, error) {
	// 创建 match 查询
	req := &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{"keyword": {Query: title}},
		},
	}

	// 查询
	res, err := global.ESClient.Search().Index(elasticsearch.ArticleIndex()).Request(req).Size(1).Do(context.TODO())
	if err != nil {
		return false, err
	}

	return res.Hits.Total.Value > 0, nil
}

// UpdateCategoryCount （在文章增/删/改时）更新文章类别的计数（+1 或 -1）
func (articleService *ArticleService) UpdateCategoryCount(tx *gorm.DB, oldCategory, newCategory string) error {
	if oldCategory == newCategory {
		return nil
	}

	// 新类别计数自增
	if newCategory != "" {
		var newArticleCategory database.ArticleCategory
		// 类别不存在，新建一个类别
		if errors.Is(tx.Where("category = ?", newCategory).First(&newArticleCategory).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleCategory{Category: newCategory, Number: 1}).Error; err != nil {
				return err
			}
		} else {
			// 类别存在，计数自增
			if err := tx.Model(&newArticleCategory).Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
				return err
			}
		}
	}

	// 旧类别计数自减
	if oldCategory != "" {
		var oldArticleCategory database.ArticleCategory
		// 更新计数
		if err := tx.Where("category = ?", oldCategory).First(&oldArticleCategory).
			Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}
		// 若更新为 0 则删除
		if oldArticleCategory.Number <= 1 {
			if err := tx.Delete(&oldArticleCategory).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateTagsCount 更新标签计数
// 比如文章的原标签是 (1,2,3)，更新后是 (2,3,4)，则需要将标签 1 的计数 -1，标签 4 的计数 +1
func (articleService *ArticleService) UpdateTagsCount(tx *gorm.DB, oldTags, newTags []string) error {
	// 获取新增和移除的标签
	addedTags, removedTags := utils.DiffArrays(oldTags, newTags)

	// 处理新增的标签
	for _, addedTag := range addedTags {
		var t database.ArticleTag
		// 查询标签，不存在则创建
		if errors.Is(tx.Where("tag = ?", addedTag).First(&t).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleTag{Tag: addedTag, Number: 1}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&t).Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
				return err
			}
		}
	}

	// 处理移除的标签
	for _, removedTag := range removedTags {
		var t database.ArticleTag
		if err := tx.Where("tag = ?", removedTag).First(&t).
			Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}

		if t.Number <= 1 {
			if err := tx.Delete(&t).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
