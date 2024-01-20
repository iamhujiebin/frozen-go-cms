package product_price_r

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/domain/model/product_price_m"
	"frozen-go-cms/myerr/bizerr"
	"frozen-go-cms/req"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
)

type MaterialPriceGetReq struct {
	PageIndex int `form:"page_index" binding:"required"`
	PageSize  int `form:"page_size" binding:"required"`
}

type MaterialPrice struct {
	Id           uint64  `json:"id"`            //  id
	MaterialName string  `json:"material_name"` //  材料名称
	MaterialCode string  `json:"material_code"` //  材料代号
	MaterialGram int64   `json:"material_gram"` //  克重
	LangC        int64   `json:"lang_c"`        //  厚度c(0.01mm)
	TonPrice     float64 `json:"ton_price"`     //  吨价
	LowPrice     float64 `json:"low_price"`     //  报价系数
	PageCover    int64   `json:"page_cover"`    //  封面
	PageInner    int64   `json:"page_inner"`    //  内页
	PageTag      int64   `json:"page_tag"`      //  tag页
	Card         int64   `json:"card"`          //  卡片
	Box          int64   `json:"box"`           //  盒子
	Index        int64   `json:"index"`         //  序号
	CreatedTime  string  `json:"created_time"`  //  创建时间
	UpdatedTime  string  `json:"updated_time"`  //  更新时间
}

// @Tags 报价系统
// @Summary 获取材料价格
// @Param Authorization header string true "token"
// @Param page_index query int false "页码"
// @Param page_size query int false "页数"
// @Success 200 {object} []MaterialPrice
// @Router /v1_0/productPrice/material [get]
func MaterialPriceGet(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param MaterialPriceGetReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	offset, limit := (param.PageIndex-1)*param.PageSize, param.PageSize
	materials, total := product_price_m.PageMaterialPrice(model, offset, limit)
	var response []MaterialPrice
	for _, v := range materials {
		var material MaterialPrice
		copier.Copy(&material, v)
		material.CreatedTime = v.CreatedTime.Format("2006-01-02 15:04:05")
		material.UpdatedTime = v.UpdatedTime.Format("2006-01-02 15:04:05")
		response = append(response, material)
	}
	resp.ResponsePageOk(c, response, total)
	return myCtx, nil
}

type MaterialPriceAdd struct {
	MaterialName string  `json:"material_name"` //  材料名称
	MaterialCode string  `json:"material_code"` //  材料代号
	MaterialGram int64   `json:"material_gram"` //  克重
	LangC        int64   `json:"lang_c"`        //  厚度c(0.01mm)
	TonPrice     float64 `json:"ton_price"`     //  吨价
	LowPrice     float64 `json:"low_price"`     //  报价系数
	PageCover    int64   `json:"page_cover"`    //  封面
	PageInner    int64   `json:"page_inner"`    //  内页
	PageTag      int64   `json:"page_tag"`      //  tag页
	Card         int64   `json:"card"`          //  卡片
	Box          int64   `json:"box"`           //  盒子
	Index        int64   `json:"index"`         //  序号
}

// @Tags 报价系统
// @Summary 新增材料价格
// @Param Authorization header string true "token"
// @Param MaterialPriceAdd body MaterialPriceAdd true "请求体"
// @Success 200
// @Router /v1_0/productPrice/material [POST]
func MaterialPricePost(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param MaterialPriceAdd
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	var materialPrice product_price_m.MaterialPrice
	copier.Copy(&materialPrice, param)
	if err := product_price_m.CreateMaterialPrice(model, materialPrice); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

type MaterialPriceUpdate struct {
	MaterialName *string  `json:"material_name"` //  材料名称
	MaterialCode *string  `json:"material_code"` //  材料代号
	MaterialGram *int64   `json:"material_gram"` //  克重
	LangC        *int64   `json:"lang_c"`        //  厚度c(0.01mm)
	TonPrice     *float64 `json:"ton_price"`     //  吨价
	LowPrice     *float64 `json:"low_price"`     //  报价系数
	PageCover    *int64   `json:"page_cover"`    //  封面
	PageInner    *int64   `json:"page_inner"`    //  内页
	PageTag      *int64   `json:"page_tag"`      //  tag页
	Card         *int64   `json:"card"`          //  卡片
	Box          *int64   `json:"box"`           //  盒子
	Index        *int64   `json:"index"`         //  序号
}

// @Tags 报价系统
// @Summary 更新材料价格
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Param MaterialPriceUpdate body MaterialPriceUpdate true "请求体"
// @Success 200
// @Router /v1_0/productPrice/material/:id [put]
func MaterialPricePut(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param MaterialPriceUpdate
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	updates := req.GetNonEmptyFields(param, "json")
	if err := product_price_m.UpdateMaterialPrice(model, id, updates); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

// @Tags 报价系统
// @Summary 删除材料价格
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Success 200
// @Router /v1_0/productPrice/material/:id [delete]
func MaterialPriceDelete(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	if err := product_price_m.DeleteMaterialPrice(model, id); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}
