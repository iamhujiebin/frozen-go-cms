package product_price_m

import (
	"frozen-go-cms/_const/enum/product_price_e"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// SizeConfig  规格尺寸
type SizeConfig struct {
	mysql.Entity
	SizeName          string  `gorm:"column:size_name"`           //  名称
	SizeCode          string  `gorm:"column:size_code"`           //  代号
	Type              int64   `gorm:"column:type"`                //  所属 1:书本 2:天地盖盒子 3:卡片
	SizeWidth         int64   `gorm:"column:size_width"`          //  宽
	SizeWidthMax      int64   `gorm:"column:size_width_max"`      //  最大宽
	SizeWidthMin      int64   `gorm:"column:size_width_min"`      //  最小宽
	SizeHeight        int64   `gorm:"column:size_height"`         //  高
	SizeHeightMax     int64   `gorm:"column:size_height_max"`     //  最大高
	SizeHeightMin     int64   `gorm:"column:size_height_min"`     //  最小高
	PerSqmX           float64 `gorm:"column:per_sqm_x"`           //  每平方米x
	PerSqmY           float64 `gorm:"column:per_sqm_y"`           //  每平方米y
	DeviceWidth       int64   `gorm:"column:device_width"`        //  上机尺寸宽
	DeviceHeight      int64   `gorm:"column:device_height"`       //  上机尺寸高
	DeviceAddBase     int64   `gorm:"column:device_add_base"`     //  上机尺寸增加基数
	DeviceAddPosition int64   `gorm:"column:device_add_position"` //  上机尺寸增加位置
	SizeOpenNum       int     `gorm:"column:size_open_num"`       //  开数
	Index             int64   `gorm:"column:index"`               //  序号
	CreateIp          string  `gorm:"column:create_ip"`           //  创建用户ip
	CreateUser        string  `gorm:"column:create_user"`         //  创建用户
	UpdateUser        string  `gorm:"column:update_user"`         //  更新用户
}

func (SizeConfig) TableName() string {
	return "size_config"
}

func CreateSizeConfig(model *domain.Model, size SizeConfig) error {
	return model.DB().Create(&size).Error
}

// 分页获取规格尺寸
func PageSizeConfig(model *domain.Model, search string, _type product_price_e.SizeConfigType, offset, limit int) ([]SizeConfig, int64) {
	var res []SizeConfig
	var total int64
	db := model.DB().Model(SizeConfig{}).Where("status = 1")
	if _type > 0 {
		db = db.Where("type = ?", _type)
	}
	if len(search) > 0 {
		args := "%" + search + "%"
		db = db.Where("size_name like ? or size_code like ?", args, args)
	}
	if err := db.Count(&total).Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		model.Log.Errorf("PageSizeConfig fail:%v", err)
	}
	return res, total
}

// 更新规格尺寸
func UpdateSizeConfig(model *domain.Model, id mysql.ID, updates map[string]interface{}) error {
	return model.DB().Table(SizeConfig{}.TableName()).Where("id = ?", id).Updates(updates).Error
}

// 删除规格尺寸
func DeleteSizeConfig(model *domain.Model, id mysql.ID) error {
	return model.DB().Table(SizeConfig{}.TableName()).Where("id = ?", id).UpdateColumn("status", 0).Error
}

// 根据类型获取规格尺寸
func GetSizeConfigByType(model *domain.Model, _type product_price_e.SizeConfigType) []SizeConfig {
	var res []SizeConfig
	if err := model.DB().Model(SizeConfig{}).Where("type = ?", _type).Find(&res).Error; err != nil {
		model.Log.Errorf("GetSizeConfigByType fail:%v", err)
	}
	return res
}

// 根据id获取规格尺寸
func GetSizeConfigById(model *domain.Model, id mysql.ID) SizeConfig {
	var size SizeConfig
	if err := model.DB().Model(SizeConfig{}).Where("id = ?", id).First(&size).Error; err != nil {
		model.Log.Errorf("GetSizeConfigById fail:%v", err)
	}
	return size
}
