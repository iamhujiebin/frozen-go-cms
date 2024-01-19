package product_price_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// ColorPrice  印刷价格
type ColorPrice struct {
	mysql.Entity
	ColorName           string  `gorm:"column:color_name"`             //  名称
	ColorCode           string  `gorm:"column:color_code"`             //  代号
	PrintStartNum       int64   `gorm:"column:print_start_num"`        //  印刷开始数量
	PrintStartPrice     float64 `gorm:"column:print_start_price"`      //  印刷开始价格
	PrintBaseNum        int64   `gorm:"column:print_base_num"`         //  印刷基数
	PrintBasePrice      float64 `gorm:"column:print_base_price"`       //  印刷基数价格
	PrintLowSumNum      int64   `gorm:"column:print_low_sum_num"`      //  印刷损耗开始数量
	PrintLowNum         int64   `gorm:"column:print_low_num"`          //  印刷损耗开始增加数量
	PrintBetweenNum     int64   `gorm:"column:print_between_num"`      //  印刷损耗间隔数量
	PrintBetweenAddNum  int64   `gorm:"column:print_between_add_num"`  //  印刷损耗间隔增加数量
	PrintStartNum2      int64   `gorm:"column:print_start_num2"`       //  相同印刷开始数量
	PrintStartPrice2    float64 `gorm:"column:print_start_price2"`     //  相同印刷开始价格
	PrintBaseNum2       int64   `gorm:"column:print_base_num2"`        //  相同印刷基数
	PrintBasePrice2     float64 `gorm:"column:print_base_price2"`      //  相同印刷基数价格
	PrintLowSumNum2     int64   `gorm:"column:print_low_sum_num2"`     //  相同印刷损耗开始数量
	PrintLowNum2        int64   `gorm:"column:print_low_num2"`         //  相同印刷损耗开始增加数量
	PrintBetweenNum2    int64   `gorm:"column:print_between_num2"`     //  相同印刷损耗间隔数量
	PrintBetweenAddNum2 int64   `gorm:"column:print_between_add_num2"` //  相同印刷损耗间隔增加数量
	PageCover           int64   `gorm:"column:page_cover"`             //  封面
	PageInner           int64   `gorm:"column:page_inner"`             //  内页
	PageTag             int64   `gorm:"column:page_tag"`               //  tag页
	Card                int64   `gorm:"column:card"`                   //  卡片
	Box                 int64   `gorm:"column:box"`                    //  盒子
	Index               int64   `gorm:"column:index"`                  //  序号
	CreateIp            string  `gorm:"column:create_ip"`              //  创建用户ip
	CreateUser          string  `gorm:"column:create_user"`            //  创建用户
	UpdateUser          string  `gorm:"column:update_user"`            //  更新用户
}

func (ColorPrice) TableName() string {
	return "color_price"
}

func CreateColorPrice(model *domain.Model, price ColorPrice) error {
	return model.DB().Create(&price).Error
}
