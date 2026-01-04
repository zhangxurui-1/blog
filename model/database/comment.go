package database

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/elasticsearch"
)

// Comment 评论表
type Comment struct {
	global.MODEL
	ArticleID string    `json:"article_id"`
	PID       *uint     `json:"p_id"`
	PComment  *Comment  `json:"-" gorm:"foreignKey:PID"`        // 父评论
	Children  []Comment `json:"children" gorm:"foreignKey:PID"` // 子评论
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`
	User      User      `json:"user" gorm:"foreignKey:UserUUID;references:UUID"`
	Content   string    `json:"content"`
}

// AfterCreate 实现 callbacks.AfterCreateInterface 接口
// gorm 支持钩子函数，这里用钩子函数实现在创建评论时更新文章的评论量
func (c *Comment) AfterCreate(_ *gorm.DB) error {
	source := "ctx._source.comments += 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).
		Script(&script).Do(context.TODO())

	return err
}

// BeforeDelete 实现 callback.BeforeDeleteInterface 接口
func (c *Comment) BeforeDelete(_ *gorm.DB) error {
	// 在 GORM 的 BeforeDelete 钩子执行时，
	// c 仅仅是一个包含主键（如 ID）的实例，GORM 不会自动填充所有字段,
	// 因此要先通过 select 拿到 article_id
	var articleID string
	if err := global.DB.Model(c).Select("article_id").First(&articleID).Error; err != nil {
		return err
	}
	source := "ctx._source.comments -= 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), articleID).
		Script(&script).Do(context.TODO())

	return err
}
