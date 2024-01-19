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
