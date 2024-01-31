package product_price_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// GenerateLog  报价日志
type GenerateLog struct {
	mysql.Entity
	ProductName string `gorm:"column:product_name"` //  产品名称
	ClientName  string `gorm:"column:client_name"`  //  客户名称
	Req         string `gorm:"column:req"`          //  请求
	Resp        string `gorm:"column:resp"`         //  返回
}

func (GenerateLog) TableName() string {
	return "generate_log"
}

func CreateGenerateLog(model *domain.Model, log GenerateLog) error {
	return model.DB().Model(GenerateLog{}).Create(&log).Error
}
