package product_price_r

import (
	"encoding/json"
	"errors"
	"fmt"
	"frozen-go-cms/_const/enum/product_price_e"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/common/resource/mysql"
	"frozen-go-cms/domain/model/product_price_m"
	"frozen-go-cms/resp"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 工艺
type Craft struct {
	Id              uint64                        `json:"id"`                // 记录id
	CraftName       string                        `json:"craft_name"`        // 名称
	CraftBodyName   string                        `json:"craft_body_name"`   // body名称
	MinSumPrice     float64                       `json:"min_sum_price"`     // 最低价格
	CraftUnitType   product_price_e.CraftUnitType `json:"craft_unit_type"`   // 0:无 1:元/张 2:元/本 3:元/次 4:元/m² 5:元/cm²
	CraftUnitPrice  float64                       `json:"craft_unit_price"`  // 单价
	CraftUnitName   string                        `json:"craft_unit_name"`   // 单位名称
	CraftNums       []int                         `json:"craft_nums"`        // 数量(比如灰板厚度用)
	DefaultCraftNum int                           `json:"default_craft_num"` // 默认数量(比如次数)
}

// 成品尺寸
type Size struct {
	Id                uint64   `json:"id"`                  // 记录id
	SizeName          string   `json:"size_name"`           // 名称
	DeviceWidthHeight [2]int64 `json:"device_width_height"` // 上机尺寸宽高
}

// 印刷颜色
type Color struct {
	Id        uint64 `json:"id"`         // 记录id
	ColorName string `json:"color_name"` // 名称
}

// 材料
type Material struct {
	Id            uint64  `json:"id"`             // 记录id
	MaterialName  string  `json:"material_name"`  // 材料名称
	MaterialGrams []int64 `json:"material_grams"` // 材料克重
}

// 封面封底
type Cover struct {
	Materials []Material `json:"materials"` // 材料
	Colors    []Color    `json:"colors"`    // 印刷颜色
	Crafts    []Craft    `json:"crafts"`    // 封面封底的工艺要求
}

// 内页
type PageInner struct {
	DefaultPageNum  int                `json:"default_page_num"`  // 默认页数
	Materials       []Material         `json:"materials"`         // 材料
	Colors          []Color            `json:"colors"`            // 印刷颜色
	PageInnerCrafts map[string][]Craft `json:"page_inner_crafts"` // 内页的工艺要求,bind_style.craft_name->具体的Crafts
}

// Tab页
type Tab struct {
	PageNum   int                `json:"page_num"`   // 页数
	Materials []Material         `json:"materials"`  // 材料
	Colors    []Color            `json:"colors"`     // 印刷颜色
	TabCrafts map[string][]Craft `json:"tab_crafts"` // tab页的工艺要求, bind_style.craft_name->具体的Crafts
}

// 产品
type Product struct {
	BindStyles      []Craft   `json:"bind_styles"`       // 装订方式
	Sizes           []Size    `json:"sizes"`             // 成品尺寸
	DefaultPrintNum int       `json:"default_print_num"` // 印刷本数
	PayMethods      []string  `json:"pay_methods"`       // 付款方式
	DeliveryTimes   []string  `json:"delivery_times"`    // 计划货期
	PriceFactors    []float64 `json:"price_factors"`     // 报价系数
	Cover           Cover     `json:"cover"`             // 封面封底
	PageInner       PageInner `json:"page_inner"`        // 内页
	BindCrafts      []Craft   `json:"bind_crafts"`       // 装订要求
	PackageCrafts   []Craft   `json:"package_crafts"`    // 包装要求
	Tab             Tab       `json:"tab"`               // tab页
}

// 自动报价配置
type AutoPriceConfig struct {
	Product    Product `json:"product"`     // 产品
	DollarRate float64 `json:"dollar_rate"` // 美元汇率
}

var (
	DefaultPrintNum = 1000                                                     // 默认印刷本数
	DefaultPageNum  = 4                                                        // 默认印刷本数
	PayMethods      = []string{"Ali Assurance", "T/T", "Paypal", "West Union"} // 付款方式
	DeliveryTimes   = []string{"3-5 working days", "5-10 working days", "10-15 working days",
		"16-20 working days", "21-25 working days", "26-30 working days"} // 计划货期

	// 目前hardcode的动作
	YOBindStyle         = "Y-O 装订"
	HardBindStyle       = "硬壳精装"
	CurrentBindStyles   = []string{YOBindStyle, HardBindStyle}                                          // 装订方式
	CoverMaterials      = []string{"PU面料", "PVC面料", "布面料", "单铜纸", "灰底白-GRAY WHITE BOARD", "双铜纸", "哑粉纸"} // 封面封底的材料
	CoverColorsIds      = []uint64{1, 4}                                                                // 封面封底印刷                                                                  // 封面封底的颜色
	CoverCrafts         = []string{"哑膜", "亮膜", "烫金", "烫银", "局部UV", "击凸", "击凹", "灰板"}                    // 封面封底工艺
	PageInnerMaterials  = []string{"单铜纸", "双胶纸", "双铜纸", "哑粉纸"}                                          // 内页的材料
	PageInnerColorsIds  = []uint64{2, 3, 4, 5, 6}                                                       // 彩色,免印                                                                   // 封面封底的颜色
	BindCrafts          = []string{"YO圈", "护角", "皮筋", "口袋", "丝带", "鸡眼", "装订"}                           // 装订要求
	TabMaterials        = []string{"PVC不干胶", "单铜纸", "普通不干胶", "双胶纸", "双铜纸", "哑粉纸"}                       // 内页的材料
	TabColorsIds        = []uint64{2, 3, 4, 5, 6}                                                       // 彩色,免印                                                                   // 封面封底的颜色
	YOPageInnerCrafts   = []string{"哑膜", "亮膜", "内分阶模切", "Tab首页加膜", "书签", "书封"}                          // YO内页工艺
	HardPageInnerCrafts = []string{"内分阶模切", "Tab首页加膜", "金边", "针孔", "书签", "书封"}                          // 硬壳内页工艺
	YOTabCrafts         = []string{"亮膜", "哑膜", "烫金", "烫银", "内分阶模切", "Tab首页加膜"}                          // YO tab页面工艺
	HardTabCrafts       = []string{"亮膜", "哑膜", "Tab首页加膜", "模切", "内分阶模切"}                                // 硬壳tab页面工艺
)

// @Tags 报价系统
// @Summary 自动报价配置
// @Param Authorization header string true "token"
// @Success 200 {object} AutoPriceConfig
// @Router /v1_0/productPrice/auto/config [get]
func AutoPriceConfigGet(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	// 系统配置
	conf := product_price_m.GetSystemConfig(model)
	// 规格尺寸
	bookSizes := product_price_m.GetSizeConfigByType(model, product_price_e.SizeConfigTypeBook)
	var sizes []Size
	for _, size := range bookSizes {
		sizes = append(sizes, Size{
			Id:                size.ID,
			SizeName:          size.SizeName,
			DeviceWidthHeight: [2]int64{size.DeviceWidth, size.DeviceHeight},
		})
	}
	// 封面封底材料
	var coverMaterials []Material
	_coverMaterials := product_price_m.GetMaterialPriceByName(model, CoverMaterials)
	coverMaterialNameMap := make(map[string][]int64)
	for _, material := range _coverMaterials {
		coverMaterialNameMap[material.MaterialName] = append(coverMaterialNameMap[material.MaterialName], material.MaterialGram)
	}
	for materialName, grams := range coverMaterialNameMap {
		coverMaterials = append(coverMaterials, Material{
			Id:            0,
			MaterialName:  materialName,
			MaterialGrams: grams,
		})
	}
	// 封面封底颜色
	var coverColors []Color
	coverColorPrice := product_price_m.GetColorPriceByIds(model, CoverColorsIds)
	for _, v := range coverColorPrice {
		coverColors = append(coverColors, Color{
			Id:        v.ID,
			ColorName: v.ColorName,
		})
	}
	// 封面封底工艺
	var coverCrafts []Craft
	_coverCrafts := product_price_m.GetCraftByCraftName(model, CoverCrafts)
	for _, v := range _coverCrafts {
		coverCrafts = append(coverCrafts, Craft{
			Id:             v.ID,
			CraftName:      v.CraftName,
			MinSumPrice:    v.MinSumPrice,
			CraftUnitName:  v.CraftUnit,
			CraftUnitPrice: v.CraftPrice,
		})
	}
	// 内页材料
	var pageInnerMaterials []Material
	_pageInnerMaterials := product_price_m.GetMaterialPriceByName(model, PageInnerMaterials)
	pageInnerMaterialNameMap := make(map[string][]int64)
	for _, material := range _pageInnerMaterials {
		pageInnerMaterialNameMap[material.MaterialName] = append(pageInnerMaterialNameMap[material.MaterialName], material.MaterialGram)
	}
	for materialName, grams := range pageInnerMaterialNameMap {
		pageInnerMaterials = append(pageInnerMaterials, Material{
			//Id:            0,
			MaterialName:  materialName,
			MaterialGrams: grams,
		})
	}
	// tab页颜色
	var tabColors []Color
	tabColorPrice := product_price_m.GetColorPriceByIds(model, TabColorsIds)
	for _, v := range tabColorPrice {
		tabColors = append(tabColors, Color{
			Id:        v.ID,
			ColorName: v.ColorName,
		})
	}
	// tab页材料
	var tabMaterials []Material
	_tabMaterials := product_price_m.GetMaterialPriceByName(model, TabMaterials)
	tabMaterialNameMap := make(map[string][]int64)
	for _, material := range _tabMaterials {
		tabMaterialNameMap[material.MaterialName] = append(tabMaterialNameMap[material.MaterialName], material.MaterialGram)
	}
	for materialName, grams := range tabMaterialNameMap {
		tabMaterials = append(tabMaterials, Material{
			//Id:            0,
			MaterialName:  materialName,
			MaterialGrams: grams,
		})
	}
	// 内页颜色
	var pageInnerColors []Color
	pageInnerColorPrice := product_price_m.GetColorPriceByIds(model, PageInnerColorsIds)
	for _, v := range pageInnerColorPrice {
		pageInnerColors = append(pageInnerColors, Color{
			Id:        v.ID,
			ColorName: v.ColorName,
		})
	}
	// 包装要求
	var packageCrafts []Craft
	packages := product_price_m.GetPackageCrafts(model)
	for _, v := range packages {
		packageCrafts = append(packageCrafts, Craft{
			Id:              v.ID,
			CraftName:       v.CraftName,
			CraftBodyName:   v.CraftBodyName,
			MinSumPrice:     v.MinSumPrice,
			CraftUnitName:   v.CraftUnit,
			CraftUnitPrice:  v.CraftPrice,
			CraftUnitType:   0,
			CraftNums:       nil,
			DefaultCraftNum: 0,
		})
	}
	response := AutoPriceConfig{
		DollarRate: conf.DollarRate, // 汇率
		Product: Product{
			Sizes:           sizes,
			DefaultPrintNum: DefaultPrintNum,
			PayMethods:      PayMethods,
			DeliveryTimes:   DeliveryTimes,
			PriceFactors:    []float64{conf.PriceFactor}, // 报价系数 todo
			Cover: Cover{
				Materials: coverMaterials,
				Colors:    coverColors,
				Crafts:    coverCrafts,
			},
			PageInner: PageInner{
				DefaultPageNum: DefaultPageNum,
				Materials:      pageInnerMaterials,
				Colors:         pageInnerColors,
			},
			PackageCrafts: packageCrafts,
			Tab: Tab{
				PageNum:   DefaultPageNum,
				Materials: tabMaterials,
				Colors:    tabColors,
			},
		},
	}
	// 装订要求
	bindCrafts := product_price_m.GetCraftByCraftBodyName(model, BindCrafts)
	for _, bindStyle := range bindCrafts {
		response.Product.BindCrafts = append(response.Product.BindCrafts, Craft{
			Id:             bindStyle.ID,
			CraftName:      bindStyle.CraftName,
			MinSumPrice:    bindStyle.MinSumPrice,
			CraftUnitPrice: bindStyle.CraftPrice,
			CraftUnitName:  bindStyle.CraftUnit,
		})
	}
	// 内页YO装订工艺
	response.Product.PageInner.PageInnerCrafts = make(map[string][]Craft)
	yoPageInnerCrafts := product_price_m.GetCraftByCraftName(model, YOPageInnerCrafts)
	for _, v := range yoPageInnerCrafts {
		response.Product.PageInner.PageInnerCrafts[YOBindStyle] = append(response.Product.PageInner.PageInnerCrafts[YOBindStyle], Craft{
			Id:             v.ID,
			CraftName:      v.CraftName,
			MinSumPrice:    v.MinSumPrice,
			CraftUnitPrice: v.CraftPrice,
			CraftUnitName:  v.CraftUnit,
		})
	}
	// 内页硬壳装订
	hardPageInnerCrafts := product_price_m.GetCraftByCraftName(model, HardPageInnerCrafts)
	for _, v := range hardPageInnerCrafts {
		response.Product.PageInner.PageInnerCrafts[HardBindStyle] = append(response.Product.PageInner.PageInnerCrafts[HardBindStyle], Craft{
			Id:             v.ID,
			CraftName:      v.CraftName,
			MinSumPrice:    v.MinSumPrice,
			CraftUnitPrice: v.CraftPrice,
			CraftUnitName:  v.CraftUnit,
		})
	}
	// Tab YO装订工艺
	response.Product.Tab.TabCrafts = make(map[string][]Craft)
	yoTabCrafts := product_price_m.GetCraftByCraftName(model, YOTabCrafts)
	for _, v := range yoTabCrafts {
		response.Product.Tab.TabCrafts[YOBindStyle] = append(response.Product.Tab.TabCrafts[YOBindStyle], Craft{
			Id:             v.ID,
			CraftName:      v.CraftName,
			MinSumPrice:    v.MinSumPrice,
			CraftUnitPrice: v.CraftPrice,
			CraftUnitName:  v.CraftUnit,
		})
	}
	hardTabCrafts := product_price_m.GetCraftByCraftName(model, HardTabCrafts)
	for _, v := range hardTabCrafts {
		response.Product.Tab.TabCrafts[HardBindStyle] = append(response.Product.Tab.TabCrafts[HardBindStyle], Craft{
			Id:             v.ID,
			CraftName:      v.CraftName,
			MinSumPrice:    v.MinSumPrice,
			CraftUnitPrice: v.CraftPrice,
			CraftUnitName:  v.CraftUnit,
		})
	}
	// 暂时只做的装订方式
	bindStyles := product_price_m.GetCraftByCraftBodyName(model, CurrentBindStyles)
	for _, bindStyle := range bindStyles {
		response.Product.BindStyles = append(response.Product.BindStyles, Craft{
			Id:             bindStyle.ID,
			CraftName:      bindStyle.CraftBodyName,
			MinSumPrice:    bindStyle.MinSumPrice,
			CraftUnitPrice: bindStyle.CraftPrice,
			CraftUnitName:  bindStyle.CraftUnit,
		})
	}
	resp.ResponseOk(c, response)
	return myCtx, nil
}

type AutoPriceReq struct {
	Order   bool `json:"order"` // 生成订单
	Product struct {
		BindStyle    int     `json:"bind_style,omitempty"`
		ClientName   string  `json:"client_name,omitempty"`
		DeliveryTime string  `json:"delivery_time,omitempty"`
		PayExtraDesc string  `json:"pay_extra_desc,omitempty"`
		PayExtraNum  int     `json:"pay_extra_num,omitempty"`
		PayExtraUnit float64 `json:"pay_extra_unit,omitempty"`
		PayMethod    string  `json:"pay_method,omitempty"`
		PriceFactor  float64 `json:"price_factor,omitempty"`
		PrintNum     int     `json:"print_num,omitempty"`
		ProductName  string  `json:"product_name,omitempty"`
		Size         uint64  `json:"size,omitempty"`
		TranDesc     string  `json:"tran_desc,omitempty"`
		TranPrice    float64 `json:"tran_price,omitempty"`
	} `json:"product,omitempty"`
	Cover struct {
		CoverColor        uint64    `json:"cover_color,omitempty"`
		CoverMaterial     string    `json:"cover_material,omitempty"`
		CoverMaterialGram int       `json:"cover_material_gram,omitempty"`
		CoverCraftIds     []uint64  `json:"cover_craft_ids"`   // 封面需要用到的工艺ids
		CoverCraftUnits   []float64 `json:"cover_craft_units"` // 上面对应的单价,要求len(CoverCraftIds)==len(CoverCraftUnits)
		CoverCraftNums    []int     `json:"cover_craft_nums"`  // 数量,同上要求
		CoverOtherCrafts  []struct {
			Name  string  `json:"name,omitempty"`
			Price float64 `json:"price,omitempty"`
		} `json:"cover_other_crafts,omitempty"`
	} `json:"cover,omitempty"`
	Inner struct {
		InnerColor        uint64    `json:"inner_color,omitempty"`
		InnerMaterial     string    `json:"inner_material,omitempty"`
		InnerMaterialGram int       `json:"inner_material_gram,omitempty"`
		InnerPageNum      int       `json:"inner_page_num,omitempty"`
		InnerCraftIds     []uint64  `json:"inner_craft_ids"`   // 内页需要用到工艺ids
		InnerCraftUnits   []float64 `json:"inner_craft_units"` // 上面对应的单价,要求len(InnerCraftIds)==len(InnerCraftUnits)
		InnerCraftNums    []int     `json:"inner_craft_nums"`  // 上面对应的单价,要求len(InnerCraftIds)==len(InnerCraftUnits)
	} `json:"inner,omitempty"`
	HasTab bool `json:"has_tab"` // 是否要Tab页
	Tab    struct {
		TabPageNum      int       `json:"tab_page_num,omitempty"`
		TabMaterial     string    `json:"tab_material,omitempty"`
		TabMaterialGram int       `json:"tab_material_gram,omitempty"`
		TabColor        uint64    `json:"tab_color,omitempty"`
		TabCraftIds     []uint64  `json:"tab_craft_ids"`   // tab面需要用到的工艺ids
		TabCraftUnits   []float64 `json:"tab_craft_units"` // 上面对应的单价,要求len(TabCraftIds)==len(TabCraftUnits)
		TabCraftNums    []int     `json:"tab_craft_nums"`  //
	} `json:"tab,omitempty"`
	Bind struct {
		BindCraftIds   []uint64  `json:"bind_craft_ids"`   //
		BindCraftUnits []float64 `json:"bind_craft_units"` //
		BindCraftNums  []int     `json:"bind_craft_nums"`  //
	} `json:"bind"`
	Package struct {
		PackageCraftIds   []uint64  `json:"package_craft_ids"`
		PackageCraftUnits []float64 `json:"package_craft_units"`
		PackageCraftNums  []int     `json:"package_craft_nums"`
	} `json:"package"`
}

type ProductDetail struct {
	PageNum      string   `json:"page_num"`      // 4P
	MaterialName string   `json:"material_name"` // PU面料等
	MaterialGram string   `json:"material_gram"` // 克重
	ColorCode    string   `json:"color_code"`    // color代号
	CraftNames   []string `json:"craft_names"`   // 用到的工艺名称s
}

type AutoPriceProductDetail struct {
	Size        string        `json:"size"`
	CoverDetail ProductDetail `json:"cover_detail"`
	InnerDetail ProductDetail `json:"inner_detail"`
	HasTab      bool          `json:"has_tab"`
	TabDetail   ProductDetail `json:"tab_detail"`
}

type AutoPriceDetail struct {
	CoverColorPrice    float64 `json:"cover_color_price"`    // 封面印刷价格
	CoverMaterialPrice float64 `json:"cover_material_price"` // 封面材料价格
	InnerColorPrice    float64 `json:"inner_color_price"`    // 内页印刷价格
	InnerMaterialPrice float64 `json:"inner_material_price"` // 内页材料价格
	TabColorPrice      float64 `json:"tab_color_price"`      // tab页印刷价格
	TabMaterialPrice   float64 `json:"tab_material_price"`   // tab页材料价格
	BindingPrice       float64 `json:"binding_price"`        // 装订价格
	PackagingPrice     float64 `json:"packaging_price"`      // 包装价格
	CraftPriceSum      float64 `json:"craft_price_sum"`      // 工艺费用: 封面+内页+tab汇总

	TranDesc      string  `json:"tran_desc"`       // 运输
	TranPrice     float64 `json:"tran_price"`      // 运输成本
	PayExtraDesc  string  `json:"pay_extra_desc"`  // 额外成本
	PayExtraPrice float64 `json:"pay_extra_price"` // 额外成本

	ProducePriceSum float64 `json:"produce_price_sum"` // 生产成本
	AllPriceSum     float64 `json:"all_price_sum"`     // 费用总和
}

func (p *AutoPriceDetail) CalProducePriceSum() {
	// 先改成两位小数
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := fieldValue.Type()

		if fieldType.Kind() == reflect.Float64 {
			newNum, _ := strconv.ParseFloat(fmt.Sprintf("%0.2f", fieldValue.Float()), 64)
			fieldValue.SetFloat(newNum)
		}
	}

	p.ProducePriceSum += p.CoverColorPrice
	p.ProducePriceSum += p.CoverMaterialPrice
	p.ProducePriceSum += p.InnerColorPrice
	p.ProducePriceSum += p.InnerMaterialPrice
	p.ProducePriceSum += p.TabColorPrice
	p.ProducePriceSum += p.TabMaterialPrice
	p.ProducePriceSum += p.BindingPrice
	p.ProducePriceSum += p.PackagingPrice
	p.ProducePriceSum += p.CraftPriceSum
}

func (p *AutoPriceDetail) CalAllProductSum() {
	if p.ProducePriceSum <= 0 { // 简单兼容一下
		p.CalProducePriceSum()
	}
	p.AllPriceSum += p.ProducePriceSum
	p.AllPriceSum += p.TranPrice
	p.AllPriceSum += p.PayExtraPrice
}

type AutoPriceResponse struct {
	ProductDetail   AutoPriceProductDetail `json:"product_detail"`    // 产品明细
	AutoPriceDetail AutoPriceDetail        `json:"auto_price_detail"` // 报价明细
	PayMethod       string                 `json:"pay_method"`        // 支付方式
	DeliverTimes    string                 `json:"deliver_times"`     // 计划货期
}

// @Tags 报价系统
// @Summary 自动报价
// @Param Authorization header string true "token"
// @Param AutoPriceReq body AutoPriceReq true "请求体"
// @Success 200 {object} AutoPriceResponse
// @Router /v1_0/productPrice/auto/generate [post]
func AutoPriceGenerate(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	model := domain.CreateModelContext(myCtx)
	var req AutoPriceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return myCtx, err
	}
	// 规格
	size := product_price_m.GetSizeConfigById(model, req.Product.Size)
	if size.ID <= 0 {
		return myCtx, errors.New("请选择成品尺寸")
	}
	if size.SizeOpenNum <= 0 {
		return myCtx, errors.New("您选择的成品尺寸,开数有问题,请到规格尺寸配置检查")
	}
	// 封面
	coverColor := product_price_m.GetColorPriceById(model, req.Cover.CoverColor)
	coverCrafts := product_price_m.GetCraftByIds(model, req.Cover.CoverCraftIds)
	var coverCraftNames []string
	for _, v := range coverCrafts {
		coverCraftNames = append(coverCraftNames, v.CraftBodyCode)
	}
	for _, v := range req.Cover.CoverOtherCrafts {
		coverCraftNames = append(coverCraftNames, v.Name)
	}
	// 内页
	innerColor := product_price_m.GetColorPriceById(model, req.Inner.InnerColor)
	innerCrafts := product_price_m.GetCraftByIds(model, req.Inner.InnerCraftIds)
	var innerCraftNames []string
	for _, v := range innerCrafts {
		innerCraftNames = append(innerCraftNames, v.CraftBodyCode)
	}
	bindCrafts := product_price_m.GetCraftByIds(model, req.Bind.BindCraftIds)
	var bindCraftNames []string
	for _, v := range bindCrafts {
		bindCraftNames = append(bindCraftNames, v.CraftBodyCode)
	}
	packageCrafts := product_price_m.GetCraftByIds(model, req.Package.PackageCraftIds)
	var packageCraftNames []string
	for _, v := range packageCrafts {
		packageCraftNames = append(packageCraftNames, v.CraftBodyCode)
	}
	// tab页面
	tabColor := product_price_m.GetColorPriceById(model, req.Tab.TabColor)
	tabCrafts := product_price_m.GetCraftByIds(model, req.Tab.TabCraftIds)
	var tabCraftNames []string
	for _, v := range tabCrafts {
		tabCraftNames = append(tabCraftNames, v.CraftBodyCode)
	}
	// 产品材料细节
	productDetail := AutoPriceProductDetail{
		Size: size.SizeCode,
		CoverDetail: ProductDetail{
			PageNum:      "4P",
			MaterialName: req.Cover.CoverMaterial,
			MaterialGram: cast.ToString(req.Cover.CoverMaterialGram),
			ColorCode:    coverColor.ColorCode,
			CraftNames:   coverCraftNames,
		},
		InnerDetail: ProductDetail{
			PageNum:      fmt.Sprintf("%dP", req.Inner.InnerPageNum),
			MaterialName: req.Inner.InnerMaterial,
			MaterialGram: cast.ToString(req.Inner.InnerMaterialGram),
			ColorCode:    innerColor.ColorCode,
			CraftNames:   innerCraftNames,
		},
		HasTab: req.HasTab,
		TabDetail: ProductDetail{
			PageNum:      fmt.Sprintf("%dP", req.Tab.TabPageNum),
			MaterialName: req.Tab.TabMaterial,
			MaterialGram: cast.ToString(req.Tab.TabMaterialGram),
			ColorCode:    tabColor.ColorCode,
			CraftNames:   tabCraftNames,
		},
	}
	// 产品报价

	// 工艺价格
	// 装订、包装、封面|内页|tab工艺要求
	var bindingPrice, packagePrice, coverCraftPrice, innerCraftPrice, tabCraftPrice float64
	bindingPrice = getCraftPrice(model, req.Product.PrintNum, 0, size, req.Bind.BindCraftIds, req.Bind.BindCraftUnits, req.Bind.BindCraftNums)
	packagePrice = getCraftPrice(model, req.Product.PrintNum, 0, size, req.Package.PackageCraftIds, req.Package.PackageCraftUnits, req.Package.PackageCraftNums)
	coverCraftPrice = getCraftPrice(model, req.Product.PrintNum, 4, size, req.Cover.CoverCraftIds, req.Cover.CoverCraftUnits, req.Cover.CoverCraftNums)
	innerCraftPrice = getCraftPrice(model, req.Product.PrintNum, req.Inner.InnerPageNum, size, req.Inner.InnerCraftIds, req.Inner.InnerCraftUnits, req.Inner.InnerCraftNums)
	tabCraftPrice = getCraftPrice(model, req.Product.PrintNum, req.Tab.TabPageNum, size, req.Tab.TabCraftIds, req.Tab.TabCraftUnits, req.Tab.TabCraftNums)
	var craftPrice float64
	craftPrice = coverCraftPrice + innerCraftPrice
	if req.HasTab {
		craftPrice += tabCraftPrice
	}
	for _, v := range req.Cover.CoverOtherCrafts {
		craftPrice += v.Price
	}
	// 价格计算
	// 封面/内页/tab印刷
	// 印刷费要先算出印刷的版数
	// 	单面印刷：张数/开数*2  color_config中的page_cover字段，1就是单面
	// 	双面印刷：P数/开数*2   同上,2就是双面
	//  ps: 张数=P数/2
	// 最后: 版数*印刷单价=印刷费用
	var coverColorPrice, innerColorPrice, tabColorPrice float64
	// 封面
	coverColorPrice = getPrintPrice(model, coverColor, coverColor.PageCover, size.SizeOpenNum, 4, req.Product.PrintNum)
	// 内页
	innerColorPrice = getPrintPrice(model, innerColor, innerColor.PageInner, size.SizeOpenNum, req.Inner.InnerPageNum, req.Product.PrintNum)
	// tab页
	if req.HasTab {
		tabColorPrice = getPrintPrice(model, tabColor, tabColor.PageTag, size.SizeOpenNum, req.Tab.TabPageNum, req.Product.PrintNum)
	}

	// 封面/内页/tab页材料:
	coverMaterial := product_price_m.GetMaterialByNameGram(model, req.Cover.CoverMaterial, req.Cover.CoverMaterialGram)
	innerMaterial := product_price_m.GetMaterialByNameGram(model, req.Inner.InnerMaterial, req.Inner.InnerMaterialGram)
	tabMaterial := product_price_m.GetMaterialByNameGram(model, req.Tab.TabMaterial, req.Tab.TabMaterialGram)
	// 大度纸 : 4 / 开数 * 本数 / 500(令数) * 克重 * 吨价 / 1884
	// 正度纸 : 4 / 开数 * 本数 / 500(令数) * 克重 * 吨价 /2325
	var coverMaterialPrice, innerMaterialPrice, tabMaterialPrice float64
	coverMaterialPrice = getMaterialPrice(model, coverMaterial, size, 4, req.Product.PrintNum)
	innerMaterialPrice = getMaterialPrice(model, innerMaterial, size, req.Inner.InnerPageNum, req.Product.PrintNum)
	if req.HasTab {
		tabMaterialPrice = getMaterialPrice(model, tabMaterial, size, req.Tab.TabPageNum, req.Product.PrintNum)
	}
	priceDetail := AutoPriceDetail{
		CoverColorPrice: coverColorPrice,
		InnerColorPrice: innerColorPrice,
		TabColorPrice:   tabColorPrice,

		CoverMaterialPrice: coverMaterialPrice,
		InnerMaterialPrice: innerMaterialPrice,
		TabMaterialPrice:   tabMaterialPrice,

		BindingPrice:   bindingPrice,
		PackagingPrice: packagePrice,
		CraftPriceSum:  craftPrice,

		PayExtraDesc:  req.Product.PayExtraDesc,
		PayExtraPrice: req.Product.PayExtraUnit * float64(req.Product.PayExtraNum),
		TranDesc:      req.Product.TranDesc,
		TranPrice:     req.Product.TranPrice,
	}
	priceDetail.CalProducePriceSum()
	priceDetail.CalAllProductSum()

	var response = AutoPriceResponse{
		ProductDetail:   productDetail,
		AutoPriceDetail: priceDetail,
		PayMethod:       req.Product.PayMethod,
		DeliverTimes:    req.Product.DeliveryTime,
	}
	if !req.Order {
		reqJson, _ := json.Marshal(req)
		respJson, _ := json.Marshal(response)
		_ = product_price_m.CreateGenerateLog(model, product_price_m.GenerateLog{
			ProductName: req.Product.ProductName,
			ClientName:  req.Product.ClientName,
			Req:         string(reqJson),
			Resp:        string(respJson),
		})
		resp.ResponseOk(c, response)
	} else {
		// 打开Excel文件
		templateFile := "template.xlsx"
		if response.AutoPriceDetail.PayExtraPrice > 0 {
			templateFile = "template_ext.xlsx"
		}
		file, err := excelize.OpenFile(templateFile)
		if err != nil {
			return myCtx, err
		}
		// 额外费用
		if response.AutoPriceDetail.PayExtraPrice > 0 {
			// C7: 额外费用
			C7Value := req.Product.PayExtraDesc
			file.SetCellValue("order", "C7", C7Value)
			// E7: 运输价格
			FE7Value := fmt.Sprintf("US$%.2f", response.AutoPriceDetail.PayExtraPrice)
			file.SetCellValue("order", "E7", FE7Value)
			file.SetCellValue("order", "F7", FE7Value)
			// A8: 总的价格
			A8Value := fmt.Sprintf("TOTAL:US$%.2f", response.AutoPriceDetail.AllPriceSum)
			file.SetCellValue("order", "A8", A8Value)
		} else {
			// A7: 总的价格
			A7Value := fmt.Sprintf("TOTAL:US$:%.2f", response.AutoPriceDetail.AllPriceSum)
			file.SetCellValue("order", "A7", A7Value)
		}
		// B5:产品名称
		B5Value := req.Product.ProductName
		file.SetCellValue("order", "B5", B5Value)
		// D5:产品数量
		D5Value := fmt.Sprintf("%d", req.Product.PrintNum)
		file.SetCellValue("order", "D5", D5Value)
		// D3: 日期
		D3Value := fmt.Sprintf(`
		DATE：%s
		PINO.：ZFA-202401003
		`, time.Now().Format("2006-01-02"))
		file.SetCellValue("order", "D3", D3Value)
		// C6: Door to Door (运输说明)
		C6Value := req.Product.TranDesc
		file.SetCellValue("order", "C6", C6Value)
		// F6: 运输价格
		EF6Value := fmt.Sprintf("US$%.2f", response.AutoPriceDetail.TranPrice)
		file.SetCellValue("order", "E6", EF6Value)
		file.SetCellValue("order", "F6", EF6Value)
		// E5: 生产单价
		if req.Product.PrintNum > 0 {
			E5Value := fmt.Sprintf("US$%0.2f", response.AutoPriceDetail.ProducePriceSum/float64(req.Product.PrintNum))
			file.SetCellValue("order", "E5", E5Value)
		}
		// F5: 生成总价
		F5Value := fmt.Sprintf("US$%0.2f", response.AutoPriceDetail.ProducePriceSum)
		file.SetCellValue("order", "F5", F5Value)

		// C5: 所有工艺
		var C5Value string
		C5Value += fmt.Sprintf("Size:%d*%dmm\n", size.SizeWidth, size.SizeHeight)
		var coverCraftsEnglish string
		for _, v := range coverCraftNames {
			coverCraftsEnglish += getEnglish(v) + " "
		}
		C5Value += fmt.Sprintf("Cover:4P %dg %s %s %s \n", // page gram material color crafts
			coverMaterial.MaterialGram, getEnglish(coverMaterial.MaterialCode), getEnglish(coverColor.ColorCode), coverCraftsEnglish)
		var innerCraftsEnglish string
		for _, v := range innerCraftNames {
			innerCraftsEnglish += getEnglish(v) + ""
		}
		C5Value += fmt.Sprintf("Inside page:%dP %dg %s %s %s \n", req.Inner.InnerPageNum, // page gram material color crafts
			innerMaterial.MaterialGram, getEnglish(innerMaterial.MaterialCode), getEnglish(coverColor.ColorCode), innerCraftsEnglish)
		if req.HasTab {
			var tabCraftsEnglish string
			for _, v := range tabCraftNames {
				tabCraftsEnglish += getEnglish(v)
			}
			C5Value += fmt.Sprintf("Tab:%dP %dg %s %s %s \n", req.Tab.TabPageNum, // page gram material color crafts
				tabMaterial.MaterialGram, getEnglish(tabMaterial.MaterialCode), getEnglish(coverColor.ColorCode), tabCraftsEnglish)
		}
		var bindCraftsEnglish string
		for _, v := range bindCraftNames {
			bindCraftsEnglish += getEnglish(v)
		}
		C5Value += fmt.Sprintf("bind:%v\n", bindCraftsEnglish)
		var packageCraftsEnglish string
		for _, v := range packageCraftNames {
			packageCraftsEnglish += getEnglish(v)
		}
		C5Value += fmt.Sprintf("package:%v\n", packageCraftsEnglish)
		file.SetCellValue("order", "C5", C5Value)

		tempFile := "uploads/file/" + fmt.Sprintf("order_%d.xlsx", time.Now().UnixNano())
		// 保存修改后的Excel文件
		err = file.SaveAs(tempFile)
		if err != nil {
			return myCtx, err
		}
		reqJson, _ := json.Marshal(req)
		respJson, _ := json.Marshal(response)
		if err := product_price_m.CreateOrderGenerate(model, product_price_m.OrderGenerate{
			ProductName: req.Product.ProductName,
			ClientName:  req.Product.ClientName,
			File:        tempFile,
			Status:      1,
			Req:         string(reqJson),
			Resp:        string(respJson),
		}); err != nil {
			return myCtx, err
		}

		//defer os.Remove(tempFile)

		// 设置响应头，告诉浏览器发送的是Excel文件
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=download.xlsx")
		c.Writer.Header().Set("Content-Type", "application/octet-stream")

		// 读取临时文件并发送给客户端
		newFile, err := os.Open(tempFile)
		if err != nil {
			model.Log.Errorf("Failed to open temporary file:%v", err)
			return myCtx, err
		}
		defer newFile.Close()

		_, err = io.Copy(c.Writer, newFile)
		if err != nil {
		} else {

		}
	}
	return myCtx, nil
}

func getEnglish(str string) string {
	arr := strings.Split(str, "_")
	if len(arr) == 2 {
		return arr[1]
	}
	return str
}

// 获取工艺价格
// param printNum:印刷本数
// param pageNum:P数
// param sizeConfig:规格尺寸
// param craftIds:工艺ids
// param units:工艺单价
// param nums:工艺数量
func getCraftPrice(model *domain.Model, printNum, pageNum int, sizeConfig product_price_m.SizeConfig,
	craftIds []mysql.ID, units []float64, nums []int) float64 {
	defer func() {
		if err := recover(); err != nil {
			model.Log.Errorf("getCraftPrice fail:%v", err)
		}
	}()
	var priceSum float64
	for i, v := range craftIds {
		if craft := product_price_m.GetCraftById(model, v); craft.ID > 0 {
			price := craft.MinSumPrice // 底价
			var unitPNum float64       // 计算价
			if craft.CraftUnit == "" || craft.CraftUnit == "件/次" {
				unitPNum = units[i] * float64(nums[i])
			}
			if craft.CraftUnit == "元/本" {
				unitPNum = units[i] * float64(printNum)
			}
			if craft.CraftUnit == "元/m²" {
				// 单价*面积
				// 需要的大纸数: 张数 / 开数 (张数=PageNum/2)
				papers := pageNum / 2 * printNum / sizeConfig.SizeOpenNum // 大纸数
				area := sizeConfig.PerSqmX * sizeConfig.PerSqmY * float64(papers)
				unitPNum = area * units[i]
			}
			if craft.CraftUnit == "元/次" {
				unitPNum = units[i]
			}
			if craft.CraftUnit == "元/本" {
				unitPNum = units[i] * float64(printNum)
			}
			// 其他新的单位

			if unitPNum > price {
				price = unitPNum
			}
			priceSum += price
		}
	}
	return priceSum
}

// 获取印刷价格
// param colorPrice: 印刷配置
// param sizeOpenNum: 开数
// param pageNum: P数
// param printNum: 打印本数
// 公式:
// 印刷费要先算出印刷的版数
//
//		单面印刷：张数/开数*2  color_config中的page_cover字段，1就是单面
//		双面印刷：P数/开数*2   同上,2就是双面
//	 ps: 张数=P数/2
//
// 印刷费用=版数*印刷单价
// 印刷车头数超一千，需要另加50/千车头(1车头=1大纸)
// 大纸 = 张数 / 开数 即 P数/2 / 开数
func getPrintPrice(model *domain.Model, colorPrice product_price_m.ColorPrice, singleDouble int, sizeOpenNum, pageNum, printNum int) float64 {
	price := colorPrice.PrintStartPrice // 开机费
	var banNum int                      // 版数
	if singleDouble == 1 {              // 单面印刷
		banNum = pageNum * 2 / 2 / sizeOpenNum // 版数
	} else if singleDouble == 2 { // 双面印刷
		banNum = pageNum / sizeOpenNum * 2
	}
	actPrice := colorPrice.PrintBasePrice * float64(banNum)
	// 加上印刷车头数
	pages := pageNum / 2 * printNum / sizeOpenNum
	head := pages / 1000
	actPrice += float64(head) * colorPrice.PrintBasePrice2

	if actPrice > price {
		price = actPrice
	}
	return price
}

// 获取材料价格
// 大度纸 : P数 / 开数 * 本数 / 500(令数) * 克重 * 吨价 / 1884
// 正度纸 : P数 / 开数 * 本数 / 500(令数) * 克重 * 吨价 /2325
func getMaterialPrice(model *domain.Model, material product_price_m.MaterialPrice, size product_price_m.SizeConfig, pageNum, printNum int) float64 {
	var sizeDivider float64
	if strings.Contains(size.SizeName, "大度") {
		sizeDivider = 1884
	}
	if strings.Contains(size.SizeName, "正度") {
		sizeDivider = 2325
	}
	price := float64(pageNum) / float64(size.SizeOpenNum) * float64(printNum) / 500 * float64(material.MaterialGram) * material.TonPrice / sizeDivider
	return price
}
