package product_price_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// OrderGenerate  订单生成历史
type OrderGenerate struct {
	mysql.Entity
	ProductName string `gorm:"column:product_name"` //  产品名称
	ClientName  string `gorm:"column:client_name"`  //  客户名称
	File        string `gorm:"column:file"`         //  文件下载路径
	Status      int    `gorm:"column:status"`
}

func (OrderGenerate) TableName() string {
	return "order_generate"
}

func CreateOrderGenerate(model *domain.Model, generate OrderGenerate) error {
	return model.DB().Model(OrderGenerate{}).Create(&generate).Error
}

// 分页获取
func PageOrderGenerate(model *domain.Model, search string, offset, limit int) ([]OrderGenerate, int64) {
	var res []OrderGenerate
	var total int64
	db := model.DB().Model(OrderGenerate{}).Where("status = 1")
	if len(search) > 0 {
		args := "%" + search + "%"
		db = db.Where("product_name like ? or client_name like ?", args, args)
	}
	if err := db.Count(&total).Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		model.Log.Errorf("PageMaterialPrice fail:%v", err)
	}
	return res, total
}

// 删除规格尺寸
func DeleteOrderGenerate(model *domain.Model, id mysql.ID) error {
	return model.DB().Table(OrderGenerate{}.TableName()).Where("id = ?", id).UpdateColumn("status", 0).Error
}
