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

type ColorPriceGetReq struct {
	PageIndex int `form:"page_index" binding:"required"`
	PageSize  int `form:"page_size" binding:"required"`
}

type ColorPrice struct {
	ColorName           string  `json:"color_name"`             //  名称
	ColorCode           string  `json:"color_code"`             //  代号
	PrintStartNum       int64   `json:"print_start_num"`        //  印刷开始数量
	PrintStartPrice     float64 `json:"print_start_price"`      //  印刷开始价格
	PrintBaseNum        int64   `json:"print_base_num"`         //  印刷基数
	PrintBasePrice      float64 `json:"print_base_price"`       //  印刷基数价格
	PrintLowSumNum      int64   `json:"print_low_sum_num"`      //  印刷损耗开始数量
	PrintLowNum         int64   `json:"print_low_num"`          //  印刷损耗开始增加数量
	PrintBetweenNum     int64   `json:"print_between_num"`      //  印刷损耗间隔数量
	PrintBetweenAddNum  int64   `json:"print_between_add_num"`  //  印刷损耗间隔增加数量
	PrintStartNum2      int64   `json:"print_start_num2"`       //  相同印刷开始数量
	PrintStartPrice2    float64 `json:"print_start_price2"`     //  相同印刷开始价格
	PrintBaseNum2       int64   `json:"print_base_num2"`        //  相同印刷基数
	PrintBasePrice2     float64 `json:"print_base_price2"`      //  相同印刷基数价格
	PrintLowSumNum2     int64   `json:"print_low_sum_num2"`     //  相同印刷损耗开始数量
	PrintLowNum2        int64   `json:"print_low_num2"`         //  相同印刷损耗开始增加数量
	PrintBetweenNum2    int64   `json:"print_between_num2"`     //  相同印刷损耗间隔数量
	PrintBetweenAddNum2 int64   `json:"print_between_add_num2"` //  相同印刷损耗间隔增加数量
	PageCover           int64   `json:"page_cover"`             //  封面
	PageInner           int64   `json:"page_inner"`             //  内页
	PageTag             int64   `json:"page_tag"`               //  tag页
	Card                int64   `json:"card"`                   //  卡片
	Box                 int64   `json:"box"`                    //  盒子
	Index               int64   `json:"index"`                  //  序号
	CreateIp            string  `json:"create_ip"`              //  创建用户ip
	CreateUser          string  `json:"create_user"`            //  创建用户
	UpdateUser          string  `json:"update_user"`            //  更新用户
	CreatedTime         string  `json:"created_time"`           //  创建时间
	UpdatedTime         string  `json:"updated_time"`           //  更新时间
}

// @Tags 报价系统
// @Summary 获取印刷价格
// @Param Authorization header string true "token"
// @Param page_index query int false "页码"
// @Param page_size query int false "页数"
// @Success 200 {object} []ColorPrice
// @Router /v1_0/productPrice/color [get]
func ColorPriceGet(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param ColorPriceGetReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	offset, limit := (param.PageIndex-1)*param.PageSize, param.PageSize
	colors, total := product_price_m.PageColorPrice(model, offset, limit)
	var response []ColorPrice
	for _, v := range colors {
		var color ColorPrice
		copier.Copy(&color, v)
		color.CreatedTime = v.CreatedTime.Format("2006-01-02 15:04:05")
		color.UpdatedTime = v.UpdatedTime.Format("2006-01-02 15:04:05")
		response = append(response, color)
	}
	resp.ResponsePageOk(c, response, total)
	return myCtx, nil
}

type ColorPriceAdd struct {
	ColorName           string  `json:"color_name"`             //  名称
	ColorCode           string  `json:"color_code"`             //  代号
	PrintStartNum       int64   `json:"print_start_num"`        //  印刷开始数量
	PrintStartPrice     float64 `json:"print_start_price"`      //  印刷开始价格
	PrintBaseNum        int64   `json:"print_base_num"`         //  印刷基数
	PrintBasePrice      float64 `json:"print_base_price"`       //  印刷基数价格
	PrintLowSumNum      int64   `json:"print_low_sum_num"`      //  印刷损耗开始数量
	PrintLowNum         int64   `json:"print_low_num"`          //  印刷损耗开始增加数量
	PrintBetweenNum     int64   `json:"print_between_num"`      //  印刷损耗间隔数量
	PrintBetweenAddNum  int64   `json:"print_between_add_num"`  //  印刷损耗间隔增加数量
	PrintStartNum2      int64   `json:"print_start_num2"`       //  相同印刷开始数量
	PrintStartPrice2    float64 `json:"print_start_price2"`     //  相同印刷开始价格
	PrintBaseNum2       int64   `json:"print_base_num2"`        //  相同印刷基数
	PrintBasePrice2     float64 `json:"print_base_price2"`      //  相同印刷基数价格
	PrintLowSumNum2     int64   `json:"print_low_sum_num2"`     //  相同印刷损耗开始数量
	PrintLowNum2        int64   `json:"print_low_num2"`         //  相同印刷损耗开始增加数量
	PrintBetweenNum2    int64   `json:"print_between_num2"`     //  相同印刷损耗间隔数量
	PrintBetweenAddNum2 int64   `json:"print_between_add_num2"` //  相同印刷损耗间隔增加数量
	PageCover           int64   `json:"page_cover"`             //  封面
	PageInner           int64   `json:"page_inner"`             //  内页
	PageTag             int64   `json:"page_tag"`               //  tag页
	Card                int64   `json:"card"`                   //  卡片
	Box                 int64   `json:"box"`                    //  盒子
}

// @Tags 报价系统
// @Summary 新增印刷价格
// @Param Authorization header string true "token"
// @Param ColorPriceAdd body ColorPriceAdd true "请求体"
// @Success 200
// @Router /v1_0/productPrice/color [POST]
func ColorPricePost(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param ColorPriceAdd
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	var colorPrice product_price_m.ColorPrice
	copier.Copy(&colorPrice, param)
	if err := product_price_m.CreateColorPrice(model, colorPrice); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

type ColorPriceUpdate struct {
	ColorName           *string  `json:"color_name"`             //  名称
	ColorCode           *string  `json:"color_code"`             //  代号
	PrintStartNum       *int64   `json:"print_start_num"`        //  印刷开始数量
	PrintStartPrice     *float64 `json:"print_start_price"`      //  印刷开始价格
	PrintBaseNum        *int64   `json:"print_base_num"`         //  印刷基数
	PrintBasePrice      *float64 `json:"print_base_price"`       //  印刷基数价格
	PrintLowSumNum      *int64   `json:"print_low_sum_num"`      //  印刷损耗开始数量
	PrintLowNum         *int64   `json:"print_low_num"`          //  印刷损耗开始增加数量
	PrintBetweenNum     *int64   `json:"print_between_num"`      //  印刷损耗间隔数量
	PrintBetweenAddNum  *int64   `json:"print_between_add_num"`  //  印刷损耗间隔增加数量
	PrintStartNum2      *int64   `json:"print_start_num2"`       //  相同印刷开始数量
	PrintStartPrice2    *float64 `json:"print_start_price2"`     //  相同印刷开始价格
	PrintBaseNum2       *int64   `json:"print_base_num2"`        //  相同印刷基数
	PrintBasePrice2     *float64 `json:"print_base_price2"`      //  相同印刷基数价格
	PrintLowSumNum2     *int64   `json:"print_low_sum_num2"`     //  相同印刷损耗开始数量
	PrintLowNum2        *int64   `json:"print_low_num2"`         //  相同印刷损耗开始增加数量
	PrintBetweenNum2    *int64   `json:"print_between_num2"`     //  相同印刷损耗间隔数量
	PrintBetweenAddNum2 *int64   `json:"print_between_add_num2"` //  相同印刷损耗间隔增加数量
	PageCover           *int64   `json:"page_cover"`             //  封面
	PageInner           *int64   `json:"page_inner"`             //  内页
	PageTag             *int64   `json:"page_tag"`               //  tag页
	Card                *int64   `json:"card"`                   //  卡片
	Box                 *int64   `json:"box"`                    //  盒子
}

// @Tags 报价系统
// @Summary 更新印刷价格
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Param ColorPriceUpdate body ColorPriceUpdate true "请求体"
// @Success 200
// @Router /v1_0/productPrice/color/:id [put]
func ColorPricePut(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param ColorPriceUpdate
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	updates := req.GetNonEmptyFields(param, "json")
	if err := product_price_m.UpdateColorPrice(model, id, updates); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

// @Tags 报价系统
// @Summary 删除印刷价格
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Success 200
// @Router /v1_0/productPrice/color/:id [delete]
func ColorPriceDelete(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	if err := product_price_m.DeleteColorPrice(model, id); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}
