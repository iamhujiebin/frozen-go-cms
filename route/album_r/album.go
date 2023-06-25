package album_r

import (
	"encoding/json"
	"frozen-go-cms/domain/model/album_m"
	"frozen-go-cms/hilo-common/domain"
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/req"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Tags 相册
// @Summary 列表
// @Param Authorization header string true "token"
// @Success 200 {object} []uint64
// @Router /v1_0/album/list [get]
func AlbumList(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	var model = domain.CreateModelContext(myCtx)
	Albums, err := album_m.GetUserAlbumsInit(model, userId)
	if err != nil {
		return myCtx, err
	}
	var ids []uint64
	for _, v := range Albums {
		ids = append(ids, v.AlbumId)
	}
	resp.ResponseOk(c, ids)
	return myCtx, nil
}

// @Tags 相册
// @Summary 明细
// @Param Authorization header string true "token"
// @Param id path integer true "Album id"
// @Success 200 {object} []string
// @Router /v1_0/album/detail/:id [get]
func AlbumDetail(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	AlbumId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var model = domain.CreateModelContext(myCtx)
	Album, err := album_m.GetUserAlbum(model, userId, AlbumId)
	if err != nil {
		return myCtx, err
	}
	var data []string
	if len(Album.Content) > 0 {
		_ = json.Unmarshal([]byte(Album.Content), &data)
	}
	resp.ResponseOk(c, data)
	return myCtx, nil
}

// @Tags 相册
// @Summary 添加
// @Param Authorization header string true "token"
// @Success 200 {object} uint64
// @Router /v1_0/album/add [post]
func AlbumAdd(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	var model = domain.CreateModelContext(myCtx)
	id, err := album_m.CreateAlbumInit(model, userId)
	if err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, id)
	return myCtx, nil
}

// @Tags 相册
// @Summary 删除
// @Param Authorization header string true "token"
// @Success 200 {object} uint64
// @Router /v1_0/album/del/:id [delete]
func AlbumDel(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var model = domain.CreateModelContext(myCtx)
	err = album_m.DeleteAlbum(model, userId, id)
	if err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}

type ProcessReq struct {
	AlbumId uint64   `json:"albumId"`
	Images  []string `json:"images"`
}

// @Tags 相册
// @Summary 明细
// @Param Authorization header string true "token"
// @Param id path integer true "Album id"
// @Param ProcessReq body ProcessReq true "请求体"
// @Success 200
// @Router /v1_0/album/detail [POST]
func AddAlbumDetail(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	userId, err := req.GetUserId(c)
	if err != nil {
		return myCtx, err
	}
	var param ProcessReq
	if err := c.ShouldBind(&param); err != nil {
		return myCtx, err
	}
	var model = domain.CreateModelContext(myCtx)
	content, _ := json.Marshal(param.Images)
	album := album_m.Album{
		UserId:  userId,
		AlbumId: param.AlbumId,
		Content: string(content),
	}
	err = album_m.UpdateAlbumInit(model, album)
	if err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}
