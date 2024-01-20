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

type CraftPriceGetReq struct {
	PageIndex int `form:"page_index" binding:"required"`
	PageSize  int `form:"page_size" binding:"required"`
}

type CraftPrice struct {
	Id              uint64  `json:"id"`                //  id
	CraftName       string  `json:"craft_name"`        //  名称
	CraftCode       string  `json:"craft_code"`        //  代号
	CraftBodyName   string  `json:"craft_body_name"`   //  子名称
	CraftBodyCode   string  `json:"craft_body_code"`   //  子代号
	CraftPrice      float64 `json:"craft_price"`       //  p单价
	BookPrice       float64 `json:"book_price"`        //  本单价
	CraftUnit       string  `json:"craft_unit"`        //  单位
	MinSumPrice     float64 `json:"min_sum_price"`     //  最小总价
	CraftModelPrice float64 `json:"craft_model_price"` //  模费单价
	TaskName        string  `json:"task_name"`         //  任务
	Index           int64   `json:"index"`             //  序号
	CreateIp        string  `json:"create_ip"`         //  创建用户ip
	CreateUser      string  `json:"create_user"`       //  创建用户
	UpdateUser      string  `json:"update_user"`       //  更新用户
	CreatedTime     string  `json:"created_time"`      //  创建时间
	UpdatedTime     string  `json:"updated_time"`      //  更新时间
}

// @Tags 报价系统
// @Summary 获取工艺价格
// @Param Authorization header string true "token"
// @Param page_index query int false "页码"
// @Param page_size query int false "页数"
// @Success 200 {object} []CraftPrice
// @Router /v1_0/productPrice/craft [get]
func CraftPriceGet(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param CraftPriceGetReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	offset, limit := (param.PageIndex-1)*param.PageSize, param.PageSize
	crafts, total := product_price_m.PageCraftPrice(model, offset, limit)
	var response []CraftPrice
	for _, v := range crafts {
		var craft CraftPrice
		copier.Copy(&craft, v)
		craft.CreatedTime = v.CreatedTime.Format("2006-01-02 15:04:05")
		craft.UpdatedTime = v.UpdatedTime.Format("2006-01-02 15:04:05")
		response = append(response, craft)
	}
	resp.ResponsePageOk(c, response, total)
	return myCtx, nil
}

type CraftPriceAdd struct {
	CraftName       string  `json:"craft_name"`        //  名称
	CraftCode       string  `json:"craft_code"`        //  代号
	CraftBodyName   string  `json:"craft_body_name"`   //  子名称
	CraftBodyCode   string  `json:"craft_body_code"`   //  子代号
	CraftPrice      float64 `json:"craft_price"`       //  p单价
	BookPrice       float64 `json:"book_price"`        //  本单价
	CraftUnit       string  `json:"craft_unit"`        //  单位
	MinSumPrice     float64 `json:"min_sum_price"`     //  最小总价
	CraftModelPrice float64 `json:"craft_model_price"` //  模费单价
	TaskName        string  `json:"task_name"`         //  任务
	Index           int64   `json:"index"`             //  序号
	CreateIp        string  `json:"create_ip"`         //  创建用户ip
	CreateUser      string  `json:"create_user"`       //  创建用户
	UpdateUser      string  `json:"update_user"`       //  更新用户
}

// @Tags 报价系统
// @Summary 新增工艺价格
// @Param Authorization header string true "token"
// @Param CraftPriceAdd body CraftPriceAdd true "请求体"
// @Success 200
// @Router /v1_0/productPrice/craft [POST]
func CraftPricePost(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param CraftPriceAdd
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	var craftPrice product_price_m.CraftPrice
	copier.Copy(&craftPrice, param)
	if err := product_price_m.CreateCraftPrice(model, craftPrice); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

type CraftPriceUpdate struct {
	CraftName       *string  `json:"craft_name"`        //  名称
	CraftCode       *string  `json:"craft_code"`        //  代号
	CraftBodyName   *string  `json:"craft_body_name"`   //  子名称
	CraftBodyCode   *string  `json:"craft_body_code"`   //  子代号
	CraftPrice      *float64 `json:"craft_price"`       //  p单价
	BookPrice       *float64 `json:"book_price"`        //  本单价
	CraftUnit       *string  `json:"craft_unit"`        //  单位
	MinSumPrice     *float64 `json:"min_sum_price"`     //  最小总价
	CraftModelPrice *float64 `json:"craft_model_price"` //  模费单价
	TaskName        *string  `json:"task_name"`         //  任务
	Index           *int64   `json:"index"`             //  序号
	CreateIp        *string  `json:"create_ip"`         //  创建用户ip
	CreateUser      *string  `json:"create_user"`       //  创建用户
	UpdateUser      *string  `json:"update_user"`       //  更新用户
}

// @Tags 报价系统
// @Summary 更新工艺价格
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Param CraftPriceUpdate body CraftPriceUpdate true "请求体"
// @Success 200
// @Router /v1_0/productPrice/craft/:id [put]
func CraftPricePut(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param CraftPriceUpdate
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	updates := req.GetNonEmptyFields(param, "json")
	if err := product_price_m.UpdateCraftPrice(model, id, updates); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

// @Tags 报价系统
// @Summary 删除工艺价格
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Success 200
// @Router /v1_0/productPrice/craft/:id [delete]
func CraftPriceDelete(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	if err := product_price_m.DeleteCraftPrice(model, id); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}
