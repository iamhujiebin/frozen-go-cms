package chatgpt_r

import (
	"encoding/json"
	"frozen-go-cms/domain/model/chatgpt_m"
	"frozen-go-cms/req"
	"frozen-go-cms/resp"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Tags Chatgpt
// @Summary 会话列表
// @Param Authorization header string true "token"
// @Success 200 {object} []uint64
// @Router /v1_0/chatgpt/session/list [get]
func SessionList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	var model = domain.CreateModelContext(myCtx)
	sessions, err := chatgpt_m.GetUserSessionsInit(model, userId)
	if err != nil {
		return myCtx, err
	}
	var ids []uint64
	for _, v := range sessions {
		ids = append(ids, v.SessionId)
	}
	resp.ResponseOk(c, ids)
	return myCtx, nil
}

// @Tags Chatgpt
// @Summary 会话明细
// @Param Authorization header string true "token"
// @Param id path integer true "session id"
// @Success 200 {object} ProcessReq
// @Router /v1_0/chatgpt/session/detail/:id [get]
func SessionDetail(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	sessionId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var model = domain.CreateModelContext(myCtx)
	session, err := chatgpt_m.GetUserSession(model, userId, sessionId)
	if err != nil {
		return myCtx, err
	}
	var data ProcessReq
	if len(session.Message) > 0 {
		_ = json.Unmarshal([]byte(session.Message), &data)
	}
	resp.ResponseOk(c, data)
	return myCtx, nil
}

// @Tags Chatgpt
// @Summary 添加一个会话
// @Param Authorization header string true "token"
// @Success 200 {object} uint64
// @Router /v1_0/chatgpt/session/add [post]
func SessionAdd(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	var model = domain.CreateModelContext(myCtx)
	id, err := chatgpt_m.CreateSessionInit(model, userId)
	if err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, id)
	return myCtx, nil
}
