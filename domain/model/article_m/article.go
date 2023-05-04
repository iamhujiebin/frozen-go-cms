package article_m

import (
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
	"gorm.io/gorm/clause"
	"time"
)

type Article struct {
	mysql.Entity
	ChannelId    mysql.ID
	Title        string
	Status       int
	CommentCount mysql.Num
	LikeCount    mysql.Num
	ReadCount    mysql.Num
	PubDate      time.Time
	CoverType    int
	CoverImages  string
	Content      string
}

// 添加文章
func AddArticle(model *domain.Model, article Article) error {
	return model.DB().Clauses(clause.OnConflict{DoNothing: true}).Create(&article).Error
}
