package ai_m

import (
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/resource/mysql"
)

// AiPrompt  咒语提示器
type AiPrompt struct {
	mysql.Entity
	Zh       string `gorm:"column:zh"`        //  中文
	En       string `gorm:"column:en"`        //  英文
	Level    uint64 `gorm:"column:level"`     //  层级 0:顶级 1:一级 2: 二级
	ParentId uint64 `gorm:"column:parent_id"` //  父id
	Code     string // 同步时候中间字段code
}

func (AiPrompt) TableName() string {
	return "ai_prompt"
}

func AddAiPrompt(model *domain.Model, prompt AiPrompt) error {
	return model.DB().Create(&prompt).Error
}
