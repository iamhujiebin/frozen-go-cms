package article_r

import (
	"encoding/json"
	"fmt"
	"frozen-go-cms/domain/model/article_m"
	"frozen-go-cms/domain/model/channel_m"
	"frozen-go-cms/myerr/bizerr"
	"frozen-go-cms/resp"
	"git.hilo.cn/hilo-common/domain"
	"git.hilo.cn/hilo-common/mycontext"
	"git.hilo.cn/hilo-common/resource/mysql"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ArticleCover struct {
	Type   int      `json:"type"`
	Images []string `json:"images"`
}

type GetArticleResp struct {
	Id        string       `json:"id"`
	Title     string       `json:"title"`
	ChannelId uint64       `json:"channel_id"`
	Content   string       `json:"content"`
	Cover     ArticleCover `json:"cover"`
	PubDate   string       `json:"pub_date"`
}

// @Tags 文章
// @Summary 详情
// @Param Authorization header string true "请求体"
// @Param id path string true "文章id"
// @Success 200 {object} GetArticleResp
// @Router /v1_0/mp/articles/:id [get]
func GetArticle(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	_id := c.Param("id")
	id, _ := strconv.ParseUint(_id, 10, 64)
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}
	model := domain.CreateModelContext(myCtx)
	article := article_m.GetArticle(model, id)
	response := GetArticleResp{
		Id:        fmt.Sprintf("%d", article.ID),
		Title:     article.Title,
		ChannelId: article.ChannelId,
		Content:   article.Content,
		Cover:     ArticleCover{Type: article.CoverType},
		PubDate:   article.PubDate.Format("2006-01-02 15:04:05"),
	}
	_ = json.Unmarshal([]byte(article.CoverImages), &response.Cover.Images)
	resp.ResponseOk(c, response)
	return myCtx, nil
}

type PostArticleReq struct {
	ChannelName string       `json:"channel_id" binding:"required"`
	Content     string       `json:"content" binding:"required"`
	Title       string       `json:"title" binding:"required"`
	Type        int          `json:"type" binding:"required"`
	Cover       ArticleCover `json:"cover"`
}

type PostPutArticleResp struct {
	Id string `json:"id"`
}

// @Tags 文章
// @Summary 发布
// @Param Authorization header string true "请求体"
// @Param draft query string false "是否草稿"
// @Param PostArticleReq body PostArticleReq true "请求体"
// @Success 200 {object} PostPutArticleResp
// @Router /v1_0/mp/articles [post]
func PostArticle(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param PostArticleReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	images, _ := json.Marshal(param.Cover.Images)
	channelId, err := channel_m.GetChannelIdByName(model, param.ChannelName)
	if err != nil {
		return myCtx, err
	}
	articleId, err := article_m.AddArticle(model, article_m.Article{
		ChannelId:    channelId,
		Title:        param.Title,
		Status:       2, // todo 审核状态默认是2
		CommentCount: uint32(rand.Intn(100)),
		LikeCount:    uint32(rand.Intn(100)),
		ReadCount:    uint32(rand.Intn(100)),
		PubDate:      time.Now(),
		CoverType:    param.Type,
		CoverImages:  string(images),
		Content:      param.Content,
	})
	if err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, PostPutArticleResp{Id: fmt.Sprintf("%d", articleId)})
	return myCtx, nil
}

type PutArticleReq struct {
	ChannelId uint64       `json:"channel_id" binding:"required"`
	Content   string       `json:"content" binding:"required"`
	Title     string       `json:"title" binding:"required"`
	Type      int          `json:"type" binding:"required"`
	Cover     ArticleCover `json:"cover"`
}

// @Tags 文章
// @Summary 修改
// @Param Authorization header string true "请求体"
// @Param draft query string false "是否草稿"
// @Param PutArticleReq body PutArticleReq true "请求体"
// @Param id path string true "修改id"
// @Success 200 {object} PostPutArticleResp
// @Router /v1_0/mp/articles/:id [put]
func PutArticle(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param PutArticleReq
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	_id := c.Param("id")
	id, _ := strconv.ParseUint(_id, 10, 64)
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}
	model := domain.CreateModelContext(myCtx)
	images, _ := json.Marshal(param.Cover.Images)
	old := article_m.GetArticle(model, id)
	articleId, err := article_m.UpdateArticle(model, article_m.Article{
		Entity: mysql.Entity{
			ID: id,
		},
		ChannelId:    param.ChannelId,
		Title:        param.Title,
		Status:       2, // todo 审核状态默认是2
		CommentCount: uint32(rand.Intn(100)),
		LikeCount:    uint32(rand.Intn(100)),
		ReadCount:    uint32(rand.Intn(100)),
		PubDate:      old.PubDate,
		CoverType:    param.Type,
		CoverImages:  string(images),
		Content:      param.Content,
	})
	if err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, PostPutArticleResp{Id: fmt.Sprintf("%d", articleId)})
	return myCtx, nil
}

// @Tags 文章
// @Summary 删除
// @Param Authorization header string true "请求体"
// @Param id path string true "文章id"
// @Success 200
// @Router /v1_0/mp/articles/:id [delete]
func DeleteArticle(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	_id := c.Param("id")
	id, _ := strconv.ParseUint(_id, 10, 64)
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}
	model := domain.CreateModelContext(myCtx)
	if err := article_m.DeleteArticle(model, id); err != nil {
		return myCtx, err
	}
	resp.ResponseOk(c, "")
	return myCtx, nil
}

type ListArticleReq struct {
	Page         int    `form:"page" binding:"required"`
	PerPage      int    `form:"per_page" binding:"required"`
	ChannelName  string `form:"channel_id"`
	BeginPubDate string `form:"begin_pubdate"` // 2006-01-02
	EndPubDate   string `form:"end_pubdate"`
}

type ListArticleResp struct {
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalCount int64         `json:"total_count"`
	Results    []ListArticle `json:"results"`
}

type ListArticle struct {
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	Cover        ArticleCover `json:"cover"`
	PubDate      string       `json:"pubdate"`
	CommentCount uint32       `json:"comment_count"`
	LikeCount    uint32       `json:"like_count"`
	ReadCount    uint32       `json:"read_count"`
	Status       int          `json:"status"`
}

// @Tags 文章
// @Summary 详情
// @Param Authorization header string true "请求体"
// @Param page query integer true "页码"
// @Param per_page query integer true "页数"
// @Param channel_id query string false "频道名称"
// @Param begin_pubdate query string true "发布开始日期:2006-01-02"
// @Param end_pubdate query string true "发布结束日期:2006-01-02"
// @Success 200 {object} ListArticleResp
// @Router /v1_0/mp/articles [get]
func PageArticle(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	var param ListArticleReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	model := domain.CreateModelContext(myCtx)
	var cid *mysql.ID
	if len(param.ChannelName) > 0 {
		channelId, err := channel_m.GetChannelIdByName(model, param.ChannelName)
		if err != nil {
			return myCtx, err
		} else {
			cid = &channelId
		}
	}
	offset, limit := (param.Page-1)*param.PerPage, param.PerPage
	article, total := article_m.PageArticle(model, cid, param.BeginPubDate, param.EndPubDate, offset, limit)
	response := ListArticleResp{
		Page:       param.Page,
		PerPage:    param.PerPage,
		TotalCount: total,
	}
	for _, a := range article {
		cva := ListArticle{
			Id:           fmt.Sprintf("%d", a.ID),
			Title:        a.Title,
			Cover:        ArticleCover{Type: a.CoverType},
			PubDate:      a.PubDate.Format("2006-01-02 15:04:05"),
			CommentCount: a.CommentCount,
			LikeCount:    a.LikeCount,
			ReadCount:    a.ReadCount,
			Status:       a.Status,
		}
		_ = json.Unmarshal([]byte(a.CoverImages), &cva.Cover.Images)
		response.Results = append(response.Results, cva)
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}
