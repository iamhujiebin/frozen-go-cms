package mgr_r

import (
	"frozen-go-cms/_const/enum/user_e"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/common/utils"
	"frozen-go-cms/domain/model/user_m"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
)

type MgrUser struct {
	Mobile string            `json:"mobile"`
	Name   string            `json:"name"`
	Gender user_e.UserGender `json:"gender"`
	Status int               `json:"status"`
}

// @Tags 管理员
// @Summary 获取用户列表
// @Param Authorization header string true "请求体"
// @Success 200 {object} []MgrUser
// @Router /v1_0/mgr/user/list [get]
func MgrUserList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	users := user_m.ListUser(model)
	var response []MgrUser
	for _, v := range users {
		response = append(response, MgrUser{
			Mobile: v.Mobile,
			Name:   v.Name,
			Gender: v.Gender,
			Status: v.Status,
		})
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

type ChangePwdReq struct {
	UserId uint64 `json:"user_id" binding:"required"`
	Pwd    string `json:"pwd" binding:"required"`
}

// @Tags 管理员
// @Summary 修改用户密码
// @Param Authorization header string true "请求体"
// @Param ChangePwdReq body ChangePwdReq true "请求体"
// @Router /v1_0/mgr/user/changePwd [post]
func MgrUserChangePwd(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param ChangePwdReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	if err := user_m.ChangeUserPwd(model, param.UserId, utils.GetMD5Str(param.Pwd)); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}
