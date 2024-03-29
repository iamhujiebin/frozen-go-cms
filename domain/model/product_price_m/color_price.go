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
	PageCover           int     `gorm:"column:page_cover"`             //  封面
	PageInner           int     `gorm:"column:page_inner"`             //  内页
	PageTag             int     `gorm:"column:page_tag"`               //  tag页
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

// 分页获取印刷价格
func PageColorPrice(model *domain.Model, search string, offset, limit int) ([]ColorPrice, int64) {
	var res []ColorPrice
	var total int64
	db := model.DB().Model(ColorPrice{}).Where("status = 1")
	if len(search) > 0 {
		args := "%" + search + "%"
		db = db.Where("color_name like ? or color_code like ?", args, args)
	}
	if err := db.Count(&total).Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		model.Log.Errorf("PageColorPrice fail:%v", err)
	}
	return res, total
}

// 更新印刷价格
func UpdateColorPrice(model *domain.Model, id mysql.ID, updates map[string]interface{}) error {
	return model.DB().Table(ColorPrice{}.TableName()).Where("id = ?", id).Updates(updates).Error
}

// 删除印刷价格
func DeleteColorPrice(model *domain.Model, id mysql.ID) error {
	return model.DB().Table(ColorPrice{}.TableName()).Where("id = ?", id).UpdateColumn("status", 0).Error
}

// 根据ids获取颜色
func GetColorPriceByIds(model *domain.Model, ids []mysql.ID) []ColorPrice {
	var res []ColorPrice
	if err := model.DB().Model(ColorPrice{}).Where("id in ?", ids).Find(&res).Error; err != nil {
		model.Log.Errorf("GetColorPriceByIds fail:%v", err)
	}
	return res
}

// 根据id获取颜色
func GetColorPriceById(model *domain.Model, id mysql.ID) ColorPrice {
	var res ColorPrice
	if err := model.DB().Model(ColorPrice{}).Where("id = ?", id).First(&res).Error; err != nil {
		model.Log.Errorf("GetColorPriceById fail:%v", err)
	}
	return res
}
