package product_price_r

import (
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/domain/model/product_price_m"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"reflect"
)

// SystemConfig  系统配置
type SystemConfig struct {
	DollarRate   float64 `json:"dollar_rate"`   //  美元汇率
	RmbRate      float64 `json:"rmb_rate"`      //  人民币汇率
	SizeFactor   float64 `json:"size_factor"`   //  体积系数
	WeightFactor float64 `json:"weight_factor"` //  重量系数
	SizeFactor2  float64 `json:"size_factor2"`  //  体积系数2
	PriceFactor  float64 `json:"price_factor"`  //  报价系数
	CoverFactor  float64 `json:"cover_factor"`  //  封面覆膜系数
	UpdateTime   string  `json:"update_time"`   //  更新时间
}

// @Tags 报价系统
// @Summary 获取系统配置
// @Param Authorization header string true "token"
// @Success 200 {object} SystemConfig
// @Router /v1_0/productPrice/system/config [get]
func SystemConfigGet(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	conf := product_price_m.GetSystemConfig(model)
	response := SystemConfig{
		DollarRate:   conf.DollarRate,
		RmbRate:      conf.RmbRate,
		SizeFactor:   conf.SizeFactor,
		WeightFactor: conf.WeightFactor,
		SizeFactor2:  conf.SizeFactor2,
		PriceFactor:  conf.PriceFactor,
		CoverFactor:  conf.CoverFactor,
		UpdateTime:   conf.UpdatedTime.Format("2006-01-02 15:04:05"),
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

// SystemConfig  系统配置
type SystemConfigUpdate struct {
	DollarRate   *float64 `form:"dollar_rate"`   //  美元汇率
	RmbRate      *float64 `form:"rmb_rate"`      //  人民币汇率
	SizeFactor   *float64 `form:"size_factor"`   //  体积系数
	WeightFactor *float64 `form:"weight_factor"` //  重量系数
	SizeFactor2  *float64 `form:"size_factor2"`  //  体积系数2
	PriceFactor  *float64 `form:"price_factor"`  //  报价系数
	CoverFactor  *float64 `form:"cover_factor"`  //  封面覆膜系数
}

// @Tags 报价系统
// @Summary 更新系统配置
// @Param Authorization header string true "token"
// @Success 200
// @Router /v1_0/productPrice/system/config [put]
func SystemConfigPut(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var param SystemConfigUpdate
	if err := c.ShouldBind(&param); err != nil {
		return myCtx, err
	}
	updates := GetNonEmptyFields(param)
	if err := product_price_m.UpdateSystemConfig(model, updates); err != nil {
		return myCtx, err
	}

	resp.ResponseOk(c, "success")
	return myCtx, nil
}

func GetNonEmptyFields(config SystemConfigUpdate) map[string]interface{} {
	result := make(map[string]interface{})

	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 判断字段是否为空
		if value.IsNil() {
			continue
		}

		// 获取tag为form的注释
		tag := field.Tag.Get("form")
		if tag != "" {
			result[tag] = value.Interface()
		}
	}

	return result
}
