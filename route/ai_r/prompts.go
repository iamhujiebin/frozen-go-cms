package ai_r

import (
	"frozen-go-cms/domain/model/ai_m"
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
)

// @Tags AIPrompts
// @Summary 获取所有提示词
// @Param Authorization header string true "token"
// @Success 200 {object} []ai_m.AiPromptData
// @Router /v1_0/ai/prompts [get]
func Prompts(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	res := ai_m.GetAllPrompts(model)
	resp.ResponseOk(c, res)
	return myCtx, nil
}
