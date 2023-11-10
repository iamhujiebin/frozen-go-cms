package ai_m

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/resource/mysql"
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

type AiPromptData struct {
	Code string `json:"code"`
	Name struct {
		En string `json:"en"`
		Zh string `json:"zh"`
	} `json:"name"`
	SubTabs []AiSAubTab `json:"subTabs"`
}

type AiSAubTab struct {
	Code string `json:"code"`
	Name struct {
		En string `json:"en"`
		Zh string `json:"zh"`
	} `json:"name"`
	Prompts []AiPrompts `json:"prompts"`
}

type AiPrompts struct {
	En   string `json:"en"`
	Zh   string `json:"zh"`
	Code string `json:"code"`
}

func AddAiPrompt(model *domain.Model, prompt AiPrompt) error {
	return model.DB().Create(&prompt).Error
}

func GetIdByCode(model *domain.Model) map[string]uint64 {
	var res = make(map[string]uint64)
	var rows []AiPrompt
	if err := model.DB().Find(&rows).Error; err != nil {
		model.Log.Error("GetIdByCode fail:%v", err)
	}
	for i := range rows {
		res[rows[i].Code] = rows[i].ID
	}
	return res
}

// 获取所有格式化提示词
// level 2-1-0去封装
func GetAllPrompts(model *domain.Model) []AiPromptData {
	var res []AiPromptData
	var rows []AiPrompt
	if err := model.DB().Find(&rows).Error; err != nil {
		model.Log.Error("GetIdByCode fail:%v", err)
		return res
	}
	var level0Level1 = make(map[uint64][]AiPrompt)
	var level1Leafs = make(map[uint64][]AiPrompt)
	for i, v := range rows {
		if v.Level == 2 {
			level1Leafs[v.ParentId] = append(level1Leafs[v.ParentId], rows[i])
		}
	}
	for i, v := range rows {
		if v.Level == 1 {
			level0Level1[v.ParentId] = append(level0Level1[v.ParentId], rows[i])
		}
	}
	for _, v := range rows {
		if v.Level == 0 {
			tmp := AiPromptData{
				Code: v.Code,
				Name: struct {
					En string `json:"en"`
					Zh string `json:"zh"`
				}{
					En: v.En,
					Zh: v.Zh,
				},
				SubTabs: nil,
			}
			for _, v2 := range level0Level1[v.ID] {
				tmp2 := AiSAubTab{
					Code: v2.Code,
					Name: struct {
						En string `json:"en"`
						Zh string `json:"zh"`
					}{
						En: v2.En,
						Zh: v2.Zh,
					},
					Prompts: nil,
				}
				for _, v3 := range level1Leafs[v2.ID] {
					tmp2.Prompts = append(tmp2.Prompts, AiPrompts{
						En:   v3.En,
						Zh:   v3.Zh,
						Code: v3.Code,
					})
				}
				tmp.SubTabs = append(tmp.SubTabs, tmp2)
			}
			res = append(res, tmp)
		}
	}
	return res
}
