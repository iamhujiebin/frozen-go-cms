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

// 分页获取材料价格
func PageMaterialPrice(model *domain.Model, search string, offset, limit int) ([]MaterialPrice, int64) {
	var res []MaterialPrice
	var total int64
	db := model.DB().Model(MaterialPrice{}).Where("status = 1")
	if len(search) > 0 {
		args := "%" + search + "%"
		db = db.Where("material_name like ? or material_code like ?", args, args)
	}
	if err := db.Count(&total).Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		model.Log.Errorf("PageMaterialPrice fail:%v", err)
	}
	return res, total
}

// 更新材料价格
func UpdateMaterialPrice(model *domain.Model, id mysql.ID, updates map[string]interface{}) error {
	return model.DB().Table(MaterialPrice{}.TableName()).Where("id = ?", id).Updates(updates).Error
}

// 删除材料价格
func DeleteMaterialPrice(model *domain.Model, id mysql.ID) error {
	return model.DB().Table(MaterialPrice{}.TableName()).Where("id = ?", id).UpdateColumn("status", 0).Error
}

// 根据名称获取材料
func GetMaterialPriceByName(model *domain.Model, names []string) []MaterialPrice {
	var rows []MaterialPrice
	if err := model.DB().Model(MaterialPrice{}).Where("`material_name` in ?", names).Find(&rows).Error; err != nil {
		model.Log.Errorf("GetMaterialPriceByName fail:%v", err)
	}
	return rows
}

// 根据名称和克重获取材料
func GetMaterialByNameGram(model *domain.Model, name string, gram int) MaterialPrice {
	var res MaterialPrice
	if err := model.DB().Model(MaterialPrice{}).Where("material_name = ? AND material_gram = ?", name, gram).First(&res).Error; err != nil {
		model.Log.Errorf("GetMaterialByNameGram fail:%v", err)
	}
	return res
}
