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

type SizeConfigGetReq struct {
	PageIndex int `form:"page_index" binding:"required"`
	PageSize  int `form:"page_size" binding:"required"`
}

type SizeConfig struct {
	Id                uint64  `json:"id"`                  //  id
	SizeName          string  `json:"size_name"`           //  名称
	SizeCode          string  `json:"size_code"`           //  代号
	Type              int64   `json:"type"`                //  所属 1:书本 2:天地盖盒子 3:卡片
	SizeWidth         int64   `json:"size_width"`          //  宽
	SizeWidthMax      int64   `json:"size_width_max"`      //  最大宽
	SizeWidthMin      int64   `json:"size_width_min"`      //  最小宽
	SizeHeight        int64   `json:"size_height"`         //  高
	SizeHeightMax     int64   `json:"size_height_max"`     //  最大高
	SizeHeightMin     int64   `json:"size_height_min"`     //  最小高
	PerSqmX           float64 `json:"per_sqm_x"`           //  每平方米x
	PerSqmY           float64 `json:"per_sqm_y"`           //  每平方米y
	DeviceWidth       int64   `json:"device_width"`        //  上机尺寸宽
	DeviceHeight      int64   `json:"device_height"`       //  上机尺寸高
	DeviceAddBase     int64   `json:"device_add_base"`     //  上机尺寸增加基数
	DeviceAddPosition int64   `json:"device_add_position"` //  上机尺寸增加位置
	SizeOpenNum       int64   `json:"size_open_num"`       //  开数
	Index             int64   `json:"index"`               //  序号
	CreatedTime       string  `json:"created_time"`        //  创建时间
	UpdatedTime       string  `json:"updated_time"`        //  更新时间
}

// @Tags 报价系统
// @Summary 获取规格尺寸
// @Param Authorization header string true "token"
// @Param page_index query int false "页码"
// @Param page_size query int false "页数"
// @Success 200 {object} []SizeConfig
// @Router /v1_0/productPrice/size [get]
func SizeConfigGet(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param SizeConfigGetReq
	if err := c.ShouldBindQuery(&param); err != nil {
		return myCtx, err
	}
	offset, limit := (param.PageIndex-1)*param.PageSize, param.PageSize
	sizes, total := product_price_m.PageSizeConfig(model, offset, limit)
	var response []SizeConfig
	for _, v := range sizes {
		var size SizeConfig
		copier.Copy(&size, v)
		size.CreatedTime = v.CreatedTime.Format("2006-01-02 15:04:05")
		size.UpdatedTime = v.UpdatedTime.Format("2006-01-02 15:04:05")
		response = append(response, size)
	}
	resp.ResponsePageOk(c, response, total)
	return myCtx, nil
}

type SizeConfigAdd struct {
	SizeName          string  `json:"size_name"`           //  名称
	SizeCode          string  `json:"size_code"`           //  代号
	Type              int64   `json:"type"`                //  所属 1:书本 2:天地盖盒子 3:卡片
	SizeWidth         int64   `json:"size_width"`          //  宽
	SizeWidthMax      int64   `json:"size_width_max"`      //  最大宽
	SizeWidthMin      int64   `json:"size_width_min"`      //  最小宽
	SizeHeight        int64   `json:"size_height"`         //  高
	SizeHeightMax     int64   `json:"size_height_max"`     //  最大高
	SizeHeightMin     int64   `json:"size_height_min"`     //  最小高
	PerSqmX           float64 `json:"per_sqm_x"`           //  每平方米x
	PerSqmY           float64 `json:"per_sqm_y"`           //  每平方米y
	DeviceWidth       int64   `json:"device_width"`        //  上机尺寸宽
	DeviceHeight      int64   `json:"device_height"`       //  上机尺寸高
	DeviceAddBase     int64   `json:"device_add_base"`     //  上机尺寸增加基数
	DeviceAddPosition int64   `json:"device_add_position"` //  上机尺寸增加位置
	SizeOpenNum       int64   `json:"size_open_num"`       //  开数
	Index             int64   `json:"index"`               //  序号
}

// @Tags 报价系统
// @Summary 新增规格尺寸
// @Param Authorization header string true "token"
// @Param SizeConfigAdd body SizeConfigAdd true "请求体"
// @Success 200
// @Router /v1_0/productPrice/size [POST]
func SizeConfigPost(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param SizeConfigAdd
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	var sizePrice product_price_m.SizeConfig
	copier.Copy(&sizePrice, param)
	if err := product_price_m.CreateSizeConfig(model, sizePrice); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

type SizeConfigUpdate struct {
	SizeName          *string  `json:"size_name"`           //  名称
	SizeCode          *string  `json:"size_code"`           //  代号
	Type              *int64   `json:"type"`                //  所属 1:书本 2:天地盖盒子 3:卡片
	SizeWidth         *int64   `json:"size_width"`          //  宽
	SizeWidthMax      *int64   `json:"size_width_max"`      //  最大宽
	SizeWidthMin      *int64   `json:"size_width_min"`      //  最小宽
	SizeHeight        *int64   `json:"size_height"`         //  高
	SizeHeightMax     *int64   `json:"size_height_max"`     //  最大高
	SizeHeightMin     *int64   `json:"size_height_min"`     //  最小高
	PerSqmX           *float64 `json:"per_sqm_x"`           //  每平方米x
	PerSqmY           *float64 `json:"per_sqm_y"`           //  每平方米y
	DeviceWidth       *int64   `json:"device_width"`        //  上机尺寸宽
	DeviceHeight      *int64   `json:"device_height"`       //  上机尺寸高
	DeviceAddBase     *int64   `json:"device_add_base"`     //  上机尺寸增加基数
	DeviceAddPosition *int64   `json:"device_add_position"` //  上机尺寸增加位置
	SizeOpenNum       *int64   `json:"size_open_num"`       //  开数
	Index             *int64   `json:"index"`               //  序号
}

// @Tags 报价系统
// @Summary 更新规格尺寸
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Param SizeConfigUpdate body SizeConfigUpdate true "请求体"
// @Success 200
// @Router /v1_0/productPrice/size/:id [put]
func SizeConfigPut(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param SizeConfigUpdate
	if err := c.ShouldBindJSON(&param); err != nil {
		return myCtx, err
	}
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	updates := req.GetNonEmptyFields(param, "json")
	if err := product_price_m.UpdateSizeConfig(model, id, updates); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

// @Tags 报价系统
// @Summary 删除规格尺寸
// @Param Authorization header string true "token"
// @Param id path int true "记录id"
// @Success 200
// @Router /v1_0/productPrice/size/:id [delete]
func SizeConfigDelete(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	id := cast.ToUint64(c.Param("id"))
	if id <= 0 {
		return myCtx, bizerr.ParaMissing
	}

	if err := product_price_m.DeleteSizeConfig(model, id); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}
