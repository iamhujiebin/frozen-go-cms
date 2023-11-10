package ai_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
)

// AiImage  ai生成的图片
type AiImage struct {
	mysql.Entity
	Prompt  string `gorm:"column:prompt" json:"prompt"` //  提示词
	Image1  string `gorm:"column:image1" json:"image1"` //  生成的图片1
	Image2  string `gorm:"column:image2" json:"image2"` //  生成的图片2
	Status  int    `gorm:"column:status" json:"status"` //  状态 1:正常
	Like    int    `gorm:"column:like" json:"like"`     //  是否喜欢 1:喜欢
	Payload string `gorm:"column:payload" json:"-"`
}

func (AiImage) TableName() string {
	return "ai_image"
}

func GetAllImages(model *domain.Model) []AiImage {
	var res []AiImage
	if err := model.DB().Model(AiImage{}).Where("status = 1").Order("id DESC").Find(&res).Error; err != nil {
		model.Log.Error("GetAllImages fail:%v", err)
	}
	return res
}

func AddAiImage(model *domain.Model, image *AiImage) error {
	return model.DB().Create(image).Error
}

func UpdateAiImage12(model *domain.Model, id mysql.ID, image1, image2 string) error {
	updates := map[string]interface{}{
		"image1": image1,
		"image2": image2,
	}
	return model.DB().Model(AiImage{}).Where("id = ?", id).Updates(updates).Error
}

func DelImage(model *domain.Model, id mysql.ID) error {
	return model.DB().Model(AiImage{}).Where("id = ?", id).UpdateColumn("status", 0).Error
}

func LikeUnLikeImage(model *domain.Model, id mysql.ID, like int) error {
	return model.DB().Model(AiImage{}).Where("id = ?", id).UpdateColumn("like", like).Error
}
