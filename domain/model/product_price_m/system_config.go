package product_price_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// SystemConfig  系统配置
type SystemConfig struct {
	mysql.Entity
	DollarRate   float64 `gorm:"column:dollar_rate"`   //  美元汇率
	RmbRate      float64 `gorm:"column:rmb_rate"`      //  人民币汇率
	SizeFactor   float64 `gorm:"column:size_factor"`   //  体积系数
	WeightFactor float64 `gorm:"column:weight_factor"` //  重量系数
	SizeFactor2  float64 `gorm:"column:size_factor2"`  //  体积系数2
	PriceFactor  float64 `gorm:"column:price_factor"`  //  报价系数
	CoverFactor  float64 `gorm:"column:cover_factor"`  //  封面覆膜系数
}

func (SystemConfig) TableName() string {
	return "system_config"
}

// 更新系统配置
func UpdateSystemConfig(model *domain.Model, updates map[string]interface{}) error {
	return model.DB().Table(SystemConfig{}.TableName()).Where("id = 1").Updates(updates).Error
}

// 获取系统配置
func GetSystemConfig(model *domain.Model) SystemConfig {
	var conf SystemConfig
	if err := model.DB().Model(SystemConfig{}).First(&conf).Error; err != nil {
		model.Log.Errorf("GetSystemConfig fail:%v", err)
	}
	return conf
}
