package casbin_r

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/domain/model/casbin_m"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
)

type CasbinInfo struct {
	Path   string `json:"path" form:"path"`
	Method string `json:"method" form:"method"`
}
type CasbinCreateRemoveRequest struct {
	UserId      uint64       `json:"user_id" binding:"required"`
	CasbinInfos []CasbinInfo `json:"casbin_infos" description:"权限模型列表"`
}

type CasbinUpdateRequest struct {
	OldPath    string     `json:"old_path"`
	OldMethod  string     `json:"old_method"`
	CasbinInfo CasbinInfo `json:"casbin_info" description:"权限模型列表"`
}

type CasbinListRequest struct {
	UserId uint64 `json:"user_id" binding:"required"`
}

type CasbinListResponse struct {
	List []CasbinInfo `json:"list" form:"list"`
}

// @Tags 权限管理
// @Summary 创建权限
// @Param CasbinCreateRemoveRequest body CasbinCreateRemoveRequest true "请求体"
// @Success 200
// @Router /v1_0/casbin/create [post]
func CasbinCreate(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param CasbinCreateRemoveRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	for _, v := range param.CasbinInfos {
		err := casbin_m.CasbinCreate(model, param.UserId, v.Path, v.Method)
		if err != nil {
			return myCtx, err
		}
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}

// @Tags 权限管理
// @Summary 删除权限
// @Param CasbinCreateRemoveRequest body CasbinCreateRemoveRequest true "请求体"
// @Success 200
// @Router /v1_0/casbin/remove [post]
func CasbinRemove(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param CasbinCreateRemoveRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	for _, v := range param.CasbinInfos {
		err := casbin_m.CasbinRemove(model, param.UserId, v.Path, v.Method)
		if err != nil {
			return myCtx, err
		}
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}

// @Tags 权限管理
// @Summary 权限列表
// @Param CasbinListRequest body CasbinListRequest true "请求体"
// @Success 200 {object} CasbinListResponse
// @Router /v1_0/casbin/list [get]
func CasbinList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param CasbinListRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	casbinList := casbin_m.CasbinList(model, param.UserId)
	var respList []CasbinInfo
	for _, host := range casbinList {
		respList = append(respList, CasbinInfo{
			Path:   host[1],
			Method: host[2],
		})
	}
	resp.ResponseOk(c, respList)
	return myCtx, nil
}

// @Tags 权限管理
// @Summary 权限测试
// @Param CasbinListRequest body CasbinListRequest true "请求体"
// @Success 200 {object} CasbinListResponse
// @Router /v1_0/casbin/test [get]
func CasbinTest(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	resp.ResponseOk(c, "success")
	return myCtx, nil
}
