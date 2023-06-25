package user_r

import (
	"fmt"
	"frozen-go-cms/_const/enum/user_e"
	"frozen-go-cms/domain/model/user_m"
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/hilo-common/resource/config"
	"frozen-go-cms/myerr/bizerr"
	"frozen-go-cms/req"
	"frozen-go-cms/req/jwt"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"time"
)

type UserAuthReq struct {
	Mobile string `json:"mobile" binding:"required"` // 手机号
	Code   string `json:"code" binding:"required"`   // 验证码
}

type UserAuthResp struct {
	Token        string `json:"token"`         // token
	RefreshToken string `json:"refresh_token"` // token
}

// @Tags 用户
// @Summary 登录
// @Param UserAuthReq body UserAuthReq true "请求体"
// @Success 200 {object} UserAuthResp
// @Router /v1_0/authorizations [post]
func UserAuth(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param UserAuthReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	if len(param.Code) < 6 || len(param.Mobile) < 6 || param.Mobile[len(param.Mobile)-6:len(param.Mobile)] != param.Code {
		return myCtx, bizerr.AuthFail
	}
	model := domain.CreateModelContext(myCtx)
	user, err := user_m.GetUserOrCreate(model, param.Mobile)
	if err != nil {
		return myCtx, err
	}
	token, err := jwt.GenerateToken(user.ID, param.Mobile, config.GetConfigJWT().ISSUER_API)
	if err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, UserAuthResp{
		Token:        token,
		RefreshToken: token,
	})
	return myCtx, nil
}

type UserProfileResp struct {
	Id       string            `json:"id"`
	Photo    string            `json:"photo"`
	Name     string            `json:"name"`
	Mobile   string            `json:"mobile"`
	Gender   user_e.UserGender `json:"gender"`
	Birthday string            `json:"birthday"`
}

// @Tags 用户
// @Summary 资料
// @Param Authorization header string true "请求体"
// @Success 200 {object} UserProfileResp
// @Router /v1_0/user/profile [get]
func UserProfile(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	user := user_m.GetUser(model, userId)
	name := fmt.Sprintf("CMS_%s", user.Mobile)
	if len(user.Name) > 0 {
		name = user.Name
	}
	resp.ResponseOk(c, UserProfileResp{
		Id:       fmt.Sprintf("%d", user.ID),
		Photo:    "",
		Name:     name,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Birthday: time.Now().Format("2006-01-02"),
	})
	return myCtx, nil
}

type PutUserProfileReq struct {
	Name   string            `json:"name"`
	Gender user_e.UserGender `json:"gender"`
}

// @Tags 用户
// @Summary 更新资料
// @Param Authorization header string true "请求体"
// @Param PutUserProfileReq body PutUserProfileReq true "请求体"
// @Success 200 {object} UserProfileResp
// @Router /v1_0/user/profile [put]
func PutUserProfile(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	var param PutUserProfileReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	if err := user_m.UpdateUser(model, userId, param.Name, param.Gender); err != nil {
		return myCtx, err
	}
	user := user_m.GetUser(model, userId)
	name := fmt.Sprintf("CMS_%s", user.Mobile)
	if len(user.Name) > 0 {
		name = user.Name
	}
	resp.ResponseOk(c, UserProfileResp{
		Id:       fmt.Sprintf("%d", user.ID),
		Photo:    "",
		Name:     name,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Birthday: time.Now().Format("2006-01-02"),
	})
	return myCtx, nil
}
