package channel_m

import (
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/resource/mysql"
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

func GetChannels(model *domain.Model) []Channel {
	var channels []Channel
	if err := model.DB().Model(Channel{}).Find(&channels).Error; err != nil {
		model.Log.Errorf("GetChannels fail:%v", err)
	}
	return channels
}

func GetChannelIdByName(model *domain.Model, name string) (mysql.ID, error) {
	var id mysql.ID
	if err := model.DB().Model(Channel{}).Where("name = ?", name).Select("id").First(&id).Error; err != nil {
		model.Log.Errorf("GetChannelIdByName fail:%v", err)
		return id, err
	}
	return id, nil
}
