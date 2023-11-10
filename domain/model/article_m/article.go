package article_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
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
func AddArticle(model *domain.Model, article Article) (uint64, error) {
	res := model.DB().Clauses(clause.OnConflict{DoNothing: true}).Create(&article)
	return article.ID, res.Error
}

func GetArticle(model *domain.Model, id mysql.ID) Article {
	var article Article
	if err := model.DB().Model(Article{}).Where("id = ?", id).First(&article); err != nil {
		model.Log.Errorf("GetArticle fail:%v", err)
	}
	return article
}

func UpdateArticle(model *domain.Model, article Article) (uint64, error) {
	res := model.DB().Save(&article)
	return article.ID, res.Error
}

func DeleteArticle(model *domain.Model, id mysql.ID) error {
	return model.DB().Model(Article{}).Where("id = ?", id).Delete(&Article{}).Error
}

func PageArticle(model *domain.Model, channelId *mysql.ID, beginDate, endDate string, offset, limit int) ([]Article, int64) {
	var article []Article
	var total int64
	db := model.DB().Model(Article{})
	if channelId != nil {
		db = db.Where("channel_id = ?", *channelId)
	}
	if len(beginDate) > 0 && len(endDate) > 0 {
		db = db.Where("DATE(pub_date) >= ? AND DATE(pub_date) <= ?", beginDate, endDate)
	}
	if err := db.Count(&total).Order("id DESC").Offset(offset).Limit(limit).Find(&article).Error; err != nil {
		model.Log.Errorf("PageArticle fail:%v", err)
	}
	return article, total
}
