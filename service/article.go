package service

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/elasticsearch"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"strconv"
	"time"
)

type ArticleService struct{}

// ArticleInfoByID 根据文章 id 获取文章
func (articleService *ArticleService) ArticleInfoByID(id string) (elasticsearch.Article, error) {
	// 异步增加文章浏览量
	go func() {
		// 将浏览量存在 redis 里，用一个定时线程存到 es
		articleView := articleService.NewArticleView()
		_ = articleView.Set(id)
	}()
	// 返回文章
	return articleService.Get(id)
}

// ArticleSearch 文章搜索
func (articleService *ArticleService) ArticleSearch(
	info request.ArticleSearch) (interface{}, int64, error) {
	// req 是一个 Elasticsearch 查询请求
	req := search.Request{Query: &types.Query{}}
	// 构造布尔查询
	boolQuery := &types.BoolQuery{}

	// 根据查询字段查询
	// 采用 should 语句，相当于 OR 逻辑，即匹配任意一个字段即可
	if info.Query != "" {
		boolQuery.Should = []types.Query{
			{Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
		}
	}

	// 根据标签字段查询
	// must 语句相当于 AND 逻辑，即必须满足标签匹配
	if info.Tag != "" {
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}},
		}
	}

	// 根据类别筛选
	// filter 不会影响相关性评分，用于精准筛选
	if info.Category != "" {
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: info.Category}}},
		}
	}

	// 如果有查询条件，则使用 Bool 查询，否则使用 MatchAll 查询
	if boolQuery.Should != nil || boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
	}

	// 设置排序字段
	if info.Sort != "" {
		var sortField string
		switch info.Sort {
		case "time":
			sortField = "created_at"
		case "view":
			sortField = "views"
		case "comment":
			sortField = "comments"
		case "like":
			sortField = "likes"
		default:
			sortField = "created_at"
		}
		var order sortorder.SortOrder
		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}

		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					sortField: {Order: &order},
				},
			},
		}
	}

	// 组装搜索请求
	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  &req,
		SourceIncludes: []string{
			"created_at",
			"cover",
			"title",
			"abstract",
			"category",
			"tags",
			"views",
			"comments",
			"likes"},
	}
	return utils.EsPagination(context.TODO(), option)
}

// ArticleCategory 获取所有文章类别及数量
func (articleService *ArticleService) ArticleCategory() ([]database.ArticleCategory, error) {
	var category []database.ArticleCategory
	if err := global.DB.Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// ArticleTags 获取所有文章标签及数量
func (articleService *ArticleService) ArticleTags() ([]database.ArticleTag, error) {
	var tags []database.ArticleTag
	if err := global.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// ArticleLike 收藏/取消收藏文章
func (articleService *ArticleService) ArticleLike(req request.ArticleLike) error {
	// 使用事务执行
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var al database.ArticleLike
		var num int

		// 查询是否已收藏过该文章，若没有则收藏
		if errors.Is(tx.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).
			First(&al).Error, gorm.ErrRecordNotFound) {

			if err := tx.Create(&database.ArticleLike{UserID: req.UserID, ArticleID: req.ArticleID}).Error; err != nil {
				return err
			}
			num = 1
		} else {
			// 如果用户已收藏过该文章，则取消收藏
			if err := tx.Delete(&al).Error; err != nil {
				return err
			}
			num = -1
		}

		// 更新 Es
		source := "ctx._source.likes += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())
		return err
	})
}

// ArticleIsLike 判断用户是否收藏了某个文章
func (articleService *ArticleService) ArticleIsLike(req request.ArticleLike) (bool, error) {
	return !errors.Is(global.DB.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).
		First(&database.ArticleLike{}).Error, gorm.ErrRecordNotFound), nil
}

// ArticleLikesList 获取用户的收藏文章列表
// 先查询数据库收藏表，再查询 Es
func (articleService *ArticleService) ArticleLikesList(info request.ArticleLikesList) (interface{}, int64, error) {
	db := global.DB.Where("user_id = ?", info.UserID)
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	// 分页查询数据库
	l, total, err := utils.MySQLPagination(&database.ArticleLike{}, option)
	if err != nil {
		return nil, 0, err
	}

	var list []struct {
		Id_     string                `json:"_id"`
		Source_ elasticsearch.Article `json:"_source"`
	}

	// 查询 Es
	for _, articleLike := range l {
		article, err := articleService.Get(articleLike.ArticleID)
		if err != nil {
			return nil, 0, err
		}

		// 将不需要的内容置为空
		article.UpdatedAt = ""
		article.Keyword = ""
		article.Content = ""

		list = append(list, struct {
			Id_     string                `json:"_id"`
			Source_ elasticsearch.Article `json:"_source"`
		}{Id_: articleLike.ArticleID, Source_: article})
	}

	return list, total, nil
}

// ArticleCreate 发布文章
func (articleService *ArticleService) ArticleCreate(req request.ArticleCreate) error {
	// 根据标题判断文章是否已存在
	if exist, err := articleService.Exists(req.Title); err != nil {
		return err
	} else if exist {
		return errors.New("article already exists")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	articleToCreate := elasticsearch.Article{
		CreatedAt: now,
		UpdatedAt: now,
		Cover:     req.Cover, // 封面是一张图片的 url
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}

	// 一个文章的创建会涉及多个表的更新
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 更新文章类别表
		if err := articleService.UpdateCategoryCount(tx, "", articleToCreate.Category); err != nil {
			return err
		}

		// 更新文章标签表
		if err := articleService.UpdateTagsCount(tx, []string{}, articleToCreate.Tags); err != nil {
			return err
		}

		// 更新封面图片资源所属的类别（更改为 "Cover"）
		if err := utils.ChangeImagesCategory(tx, []string{articleToCreate.Cover}, appTypes.Cover); err != nil {
			return err
		}

		// 管理文章内的插图
		illustrations, err := utils.FindIllustrations(articleToCreate.Content)
		if err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, illustrations, appTypes.Illustration); err != nil {
			return err
		}

		// 创建文章
		return articleService.Create(&articleToCreate)
	})
}

// ArticleDelete 删除文章
func (articleService *ArticleService) ArticleDelete(req request.ArticleDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}

	// 一个文章的删除会涉及多个表的更新
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			// 查询该文章
			articleToDelete, err := articleService.Get(id)
			if err != nil {
				return err
			}
			// 更新文章类别表
			if err2 := articleService.UpdateCategoryCount(tx, articleToDelete.Category, ""); err2 != nil {
				return err2
			}
			// 更新文章标签表
			if err2 := articleService.UpdateTagsCount(tx, articleToDelete.Tags, []string{}); err2 != nil {
				return err2
			}

			// 更新封面图片资源所属的类别（更改为 "Null"）
			if err2 := utils.InitImagesCategory(tx, []string{articleToDelete.Cover}); err2 != nil {
				return err2
			}
			// 更新文章内的插图（插图的类别更改为 "Null"）
			illustrations, err2 := utils.FindIllustrations(articleToDelete.Content)
			if err2 != nil {
				return err2
			}
			if err := utils.InitImagesCategory(tx, illustrations); err != nil {
				return err
			}
			// 删除该文章下的所有评论
			comments, err := ServiceGroupApp.CommentService.
				CommentInfoByArticleID(request.CommentInfoByArticleID{ArticleID: id})
			for _, comment := range comments {
				if err := ServiceGroupApp.CommentService.DeleteCommentAndChildren(tx, comment.ID); err != nil {
					return err
				}
			}
		}
		// 删除文章
		return articleService.Delete(req.IDs)
	})
}

// ArticleUpdate 更新文章
func (articleService *ArticleService) ArticleUpdate(req request.ArticleUpdate) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	// 使用一个匿名结构体
	articleToUpdate := struct {
		UpdatedAt string   `json:"updated_at"`
		Cover     string   `json:"cover"`   // 文章封面
		Title     string   `json:"title"`   // 文章标题
		Keyword   string   `json:"keyword"` // 文章标题-关键字
		Category  string   `json:"category"`
		Tags      []string `json:"tags"`
		Abstract  string   `json:"abstract"`
		Content   string   `json:"content"`
	}{
		UpdatedAt: now,
		Cover:     req.Cover, // 封面是一张图片的 url
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		oldArticle, err := articleService.Get(req.ID)
		if err != nil {
			return err
		}
		// 更新文章类别表
		if err := articleService.UpdateCategoryCount(tx, oldArticle.Category, articleToUpdate.Category); err != nil {
			return err
		}
		// 更新标签表
		if err := articleService.UpdateTagsCount(tx, oldArticle.Tags, articleToUpdate.Tags); err != nil {
			return err
		}
		// 更新封面图片资源所属的类别
		if oldArticle.Cover != articleToUpdate.Cover {
			if err := utils.InitImagesCategory(tx, []string{oldArticle.Cover}); err != nil {
				return err
			}
			if err := utils.ChangeImagesCategory(tx, []string{articleToUpdate.Cover}, appTypes.Cover); err != nil {
				return err
			}
		}
		// 更新插图资源
		oldIllustrations, err := utils.FindIllustrations(oldArticle.Content)
		if err != nil {
			return err
		}
		newIllustrations, err := utils.FindIllustrations(articleToUpdate.Content)
		if err != nil {
			return err
		}
		addedIllustrations, removedIllustrations := utils.DiffArrays(oldIllustrations, newIllustrations)
		if err := utils.InitImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, addedIllustrations, appTypes.Illustration); err != nil {
			return err
		}

		// 更新 Es 中的内容
		return articleService.Update(req.ID, articleToUpdate)
	})
}

// ArticleList 获取文章列表
func (articleService *ArticleService) ArticleList(info request.ArticleList) (list interface{}, total int64, err error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}
	// 根据标题查询
	if info.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{
			Match: map[string]types.MatchQuery{"title": {Query: *info.Title}},
		})
	}
	// 根据 abstract 查询
	if info.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{
			Match: map[string]types.MatchQuery{"abstract": {Query: *info.Abstract}},
		})
	}
	// 根据类别筛选
	if info.Category != nil {
		boolQuery.Filter = []types.Query{
			{
				Term: map[string]types.TermQuery{"category": {Value: *info.Category}},
			},
		}
	}
	// 根据条件执行查询
	if boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
		// 默认按文章发布时间排序
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"created_at": {Order: &sortorder.Desc},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
	}

	return utils.EsPagination(context.TODO(), option)
}
