package channel_m

import (
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/resource/mysql"
	"gorm.io/gorm/clause"
)

type Channel struct {
	mysql.Entity
	Name string
}

// 添加channel
func AddChannel(model *domain.Model, channel Channel) error {
	return model.DB().Clauses(clause.OnConflict{DoNothing: true}).Create(&channel).Error
}
