package product_price_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// MaterialPrice  材料价格
type MaterialPrice struct {
	mysql.Entity
	MaterialName string  `gorm:"column:material_name"` //  材料名称
	MaterialCode string  `gorm:"column:material_code"` //  材料代号
	MaterialGram int64   `gorm:"column:material_gram"` //  克重
	LangC        int64   `gorm:"column:lang_c"`        //  厚度c(0.01mm)
	TonPrice     float64 `gorm:"column:ton_price"`     //  吨价
	LowPrice     float64 `gorm:"column:low_price"`     //  报价系数
	PageCover    int64   `gorm:"column:page_cover"`    //  封面
	PageInner    int64   `gorm:"column:page_inner"`    //  内页
	PageTag      int64   `gorm:"column:page_tag"`      //  tag页
	Card         int64   `gorm:"column:card"`          //  卡片
	Box          int64   `gorm:"column:box"`           //  盒子
	Index        int64   `gorm:"column:index"`         //  序号
	CreateIp     string  `gorm:"column:create_ip"`     //  创建用户ip
	CreateUser   string  `gorm:"column:create_user"`   //  创建用户
	UpdateUser   string  `gorm:"column:update_user"`   //  更新用户
}

func (MaterialPrice) TableName() string {
	return "material_price"
}

func CreateMaterialPrice(model *domain.Model, material MaterialPrice) error {
	return model.DB().Create(&material).Error
}