package product_price_r

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/domain/model/product_price_m"
	"frozen-go-cms/myerr/bizerr"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
)

type OrderGetReq struct {
	Search    string `form:"search"`
	PageIndex int    `form:"page_index" binding:"required"`
	PageSize  int    `form:"page_size" binding:"required"`
}

type Order struct {
	Id          uint64 `json:"id"`           //  id
	ProductName string `json:"product_name"` //  产品名称
	ClientName  string `json:"client_name"`  //  客户名称
	File        string `json:"file"`         //  文件路径
	CreatedTime string `json:"created_time"` //  创建时间
	UpdatedTime string `json:"updated_time"` //  更新时间
}

// @Tags 报价系统
// @Summary 订单列表
// @Param Authorization header string true "token"
// @Param page_index query int false "页码"
// @Param page_size query int false "页数"
// @Param search query string false "搜索词"
// @Success 200 {object} []Order
// @Router /v1_0/productPrice/orders [get]
func OrderGet(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param OrderGetReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	offset, limit := (param.PageIndex-1)*param.PageSize, param.PageSize
	orders, total := product_price_m.PageOrderGenerate(model, param.Search, offset, limit)
	var response []Order
	for _, v := range orders {
		var order Order
		copier.Copy(&order, v)
		order.File = "https://api.frozenhu.cn/" + order.File
		order.CreatedTime = v.CreatedTime.Format("2006-01-02 15:04:05")
		order.UpdatedTime = v.UpdatedTime.Format("2006-01-02 15:04:05")
		response = append(response, order)
	}
	resp.ResponsePageOk(c, response, total)
	return myCtx, nil
}

// @Tags 报价系统
// @Summary 删除历史订单
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Success 200
// @Router /v1_0/productPrice/order/:id [delete]
func OrderDelete(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	if err := product_price_m.DeleteOrderGenerate(model, id); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}
