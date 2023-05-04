package user_r

import (
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
)

// @Tags 用户背包
// @Summary 获取用户的背包
// @Param resType path int true "类型：1 礼物"
// @Success 200
// @Router /v1/user/bag/{resType} [get]
func UserBag(c *gin.Context) error {
	resp.ResponseOk(c, "")
	return nil
}
