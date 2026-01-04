package utils

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"golang.org/x/net/context"
	"server/global"
	"server/model/other"
)

func MySQLPagination[T any](model *T, option other.MySQLOption) (list []T, total int64, err error) {
	// 设置分页默认值
	if option.Page < 1 {
		option.Page = 1 // 页码默认为 1
	}
	if option.PageSize < 1 {
		option.PageSize = 10 // 每页记录数默认为 10
	}
	if option.Order == "" {
		option.Order = "id desc"
	}

	// 创建查询
	query := global.DB.Model(model)

	// 如果传入了额外的 WHERE 条件，则应用这些条件
	if option.Where != nil {
		query = query.Where(option.Where)
	}

	// 计算符合条件的记录总数
	if err = query.Count(&total).Error; err != nil {
		// 查询总数失败，则返回错误
		return nil, 0, err
	}

	// 预加载关联模型
	for _, preload := range option.PreLoad {
		query = query.Preload(preload)
	}

	// 应用分页查询
	err = query.Order(option.Order).
		Limit(option.PageSize).
		Offset(option.PageSize * (option.Page - 1)).
		Find(&list).Error

	return list, total, err
}

// EsPagination 实现 Elasticsearch 数据分页查询
func EsPagination(ctx context.Context,
	option other.EsOption) (list []types.Hit, total int64, err error) {

	// 设置分页默认值
	if option.Page < 1 {
		option.Page = 1
	}
	if option.PageSize < 1 {
		option.PageSize = 10
	}

	// 设置 Es 查询的分页值
	from := (option.Page - 1) * option.PageSize // 起始位置
	option.Request.Size = &option.PageSize      // 每页记录数
	option.Request.From = &from                 // 起始位置

	// 执行 Es 分页查询
	res, err := global.ESClient.Search().
		Index(option.Index).                       // 指定索引
		Request(option.Request).                   // 应用查询 request
		SourceIncludes_(option.SourceIncludes...). // 需要包含的字段
		Do(ctx)                                    // 执行

	if err != nil {
		return nil, 0, err
	}

	// 提取查询结果
	list = res.Hits.Hits         // 查询结果中的文档
	total = res.Hits.Total.Value // 文档总数
	return list, total, nil
}
