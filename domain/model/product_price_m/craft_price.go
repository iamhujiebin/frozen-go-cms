package product_price_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// CraftPrice  工艺价格
type CraftPrice struct {
	mysql.Entity
	CraftName       string  `gorm:"column:craft_name"`        //  名称
	CraftCode       string  `gorm:"column:craft_code"`        //  代号
	CraftBodyName   string  `gorm:"column:craft_body_name"`   //  子名称
	CraftBodyCode   string  `gorm:"column:craft_body_code"`   //  子代号
	CraftPrice      float64 `gorm:"column:craft_price"`       //  p单价
	BookPrice       float64 `gorm:"column:book_price"`        //  本单价
	CraftUnit       string  `gorm:"column:craft_unit"`        //  单位
	MinSumPrice     float64 `gorm:"column:min_sum_price"`     //  最小总价
	CraftModelPrice float64 `gorm:"column:craft_model_price"` //  模费单价
	TaskName        string  `gorm:"column:task_name"`         //  任务
	Index           int64   `gorm:"column:index"`             //  序号
	CreateIp        string  `gorm:"column:create_ip"`         //  创建用户ip
	CreateUser      string  `gorm:"column:create_user"`       //  创建用户
	UpdateUser      string  `gorm:"column:update_user"`       //  更新用户
}

func (CraftPrice) TableName() string {
	return "craft_price"
}

func CreateCraftPrice(model *domain.Model, color CraftPrice) error {
	return model.DB().Create(&color).Error
}

// 分页获取工艺价格
func PageCraftPrice(model *domain.Model, offset, limit int) ([]CraftPrice, int64) {
	var res []CraftPrice
	var total int64
	if err := model.DB().Model(CraftPrice{}).Where("status = 1").
		Count(&total).
		Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		model.Log.Errorf("PageCraftPrice fail:%v", err)
	}
	return res, total
}

// 更新工艺价格
func UpdateCraftPrice(model *domain.Model, id mysql.ID, updates map[string]interface{}) error {
	return model.DB().Table(CraftPrice{}.TableName()).Where("id = ?", id).Updates(updates).Error
}

// 删除工艺价格
func DeleteCraftPrice(model *domain.Model, id mysql.ID) error {
	return model.DB().Table(CraftPrice{}.TableName()).Where("id = ?", id).UpdateColumn("status", 0).Error
}

// 根据CraftBodyName找数据
func GetCraftByCraftBodyName(model *domain.Model, names []string) []CraftPrice {
	var rows []CraftPrice
	if err := model.DB().Model(CraftPrice{}).Where("craft_body_name in ?", names).Find(&rows).Error; err != nil {
		model.Log.Errorf("GetCraftByCraftBodyName fail:%v", err)
	}
	return rows
}

// 根据CraftName找数据
func GetCraftByCraftName(model *domain.Model, names []string) []CraftPrice {
	var rows []CraftPrice
	if err := model.DB().Model(CraftPrice{}).Where("craft_name in ?", names).Find(&rows).Error; err != nil {
		model.Log.Errorf("GetCraftByCraftName fail:%v", err)
	}
	return rows
}

// 获取所有包装要求
func GetPackageCrafts(model *domain.Model) []CraftPrice {
	var rows []CraftPrice
	if err := model.DB().Model(CraftPrice{}).Where("craft_name = '包装要求'").Find(&rows).Error; err != nil {
		model.Log.Errorf("GetPackageCrafts fail:%v", err)
	}
	return rows
}

// 根据ids获取工艺
func GetCraftByIds(model *domain.Model, ids []mysql.ID) []CraftPrice {
	var res []CraftPrice
	if err := model.DB().Model(CraftPrice{}).Where("id in ?", ids).Find(&res).Error; err != nil {
		model.Log.Errorf("GetCraftByIds fail:%v", err)
	}
	return res
}
