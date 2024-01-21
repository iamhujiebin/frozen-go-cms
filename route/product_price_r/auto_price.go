package product_price_r

import (
	"frozen-go-cms/_const/enum/product_price_e"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/domain/model/product_price_m"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
)

// 工艺
type Craft struct {
	Id              uint64                        `json:"id"`                // 记录id
	CraftName       string                        `json:"craft_name"`        // 名称
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
}

// 内页
type PageInner struct {
	DefaultPageNum int        `json:"default_page_num"` // 默认页数
	Materials      []Material `json:"materials"`        // 材料
	Colors         []Color    `json:"colors"`           // 印刷颜色
}

// Tab页
type Tab struct {
	PageNum   int        `json:"page_num"`  // 页数
	Materials []Material `json:"materials"` // 材料
	Colors    []Color    `json:"colors"`    // 印刷颜色
}

// 产品
type Product struct {
	BindStyles          []Craft            `json:"bind_styles"`            // 装订方式
	Sizes               []Size             `json:"sizes"`                  // 成品尺寸
	DefaultPrintNum     int                `json:"default_print_num"`      // 印刷本数
	PayMethods          []string           `json:"pay_methods"`            // 付款方式
	DeliveryTimes       []string           `json:"delivery_times"`         // 计划货期
	PriceFactors        []float64          `json:"price_factors"`          // 报价系数
	Cover               Cover              `json:"cover"`                  // 封面封底
	CoverCrafts         []Craft            `json:"cover_crafts"`           // 封面封底的工艺要求, bind_style.craft_name->具体的Crafts
	PageInner           PageInner          `json:"page_inner"`             // 内页
	PageInnerCrafts     map[string][]Craft `json:"page_inner_crafts"`      // 内页的工艺要求,bind_style.craft_name->具体的Crafts
	PageInnerBindStyles []Craft            `json:"page_inner_bind_styles"` // 内页的装订要求,bind_style.craft_name->具体的Crafts
	PageInnerPackage    []Craft            `json:"page_inner_package"`     // 内页的包装要求,bind_style.craft_name->具体的Crafts
	Tab                 Tab                `json:"tab"`                    // tab页
	TabCrafts           map[string][]Craft `json:"tab_crafts"`             // tab页的工艺要求, bind_style.craft_name->具体的Crafts
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
	PageInnerBindStyles = []string{"YO圈", "护角", "皮筋", "口袋", "丝带", "鸡眼", "装订"}                           // 内页的装订要求-工艺
	TabMaterials        = []string{"PVC不干胶", "单铜纸", "普通不干胶", "双胶纸", "双铜纸", "哑粉纸"}                       // 内页的材料
	TabColorsIds        = []uint64{2, 3, 4, 5, 6}                                                       // 彩色,免印                                                                   // 封面封底的颜色
	YOPageInnerCrafts   = []string{"哑膜", "亮膜", "内分阶模切", "Tab首页加膜", "书签", "书封"}                          // YO内页工艺
	HardPageInnerCrafts = []string{"内分阶模切", "Tab首页加膜", "金边", "针孔", "书签", "书封"}                          // 硬壳内页工艺
	YOTabCrafts         = []string{"亮膜", "哑膜", "烫金", "烫银", "内分阶模切", "Tab首页加膜"}                          // YO tab页面工艺
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
	var pageInnerPackages []Craft
	packages := product_price_m.GetPackageCrafts(model)
	for _, v := range packages {
		pageInnerPackages = append(pageInnerPackages, Craft{
			Id:              v.ID,
			CraftName:       v.CraftName,
			MinSumPrice:     v.MinSumPrice,
			CraftUnitName:   v.CraftUnit,
			CraftUnitType:   0,
			CraftUnitPrice:  0,
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
			},
			CoverCrafts: coverCrafts,
			PageInner: PageInner{
				DefaultPageNum: DefaultPageNum,
				Materials:      pageInnerMaterials,
				Colors:         pageInnerColors,
			},
			PageInnerPackage: pageInnerPackages,
			Tab: Tab{
				PageNum:   DefaultPageNum,
				Materials: tabMaterials,
				Colors:    tabColors,
			},
		},
	}
	// 内页装订要求
	pageInnerBindStyles := product_price_m.GetCraftByCraftBodyName(model, PageInnerBindStyles)
	for _, bindStyle := range pageInnerBindStyles {
		response.Product.PageInnerBindStyles = append(response.Product.PageInnerBindStyles, Craft{
			Id:             bindStyle.ID,
			CraftName:      bindStyle.CraftName,
			MinSumPrice:    bindStyle.MinSumPrice,
			CraftUnitPrice: bindStyle.CraftPrice,
			CraftUnitName:  bindStyle.CraftUnit,
		})
	}
	// 内页YO装订工艺
	response.Product.PageInnerCrafts = make(map[string][]Craft)
	yoPageInnerCrafts := product_price_m.GetCraftByCraftName(model, YOPageInnerCrafts)
	for _, v := range yoPageInnerCrafts {
		response.Product.PageInnerCrafts[YOBindStyle] = append(response.Product.PageInnerCrafts[YOBindStyle], Craft{
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
		response.Product.PageInnerCrafts[HardBindStyle] = append(response.Product.PageInnerCrafts[HardBindStyle], Craft{
			Id:             v.ID,
			CraftName:      v.CraftName,
			MinSumPrice:    v.MinSumPrice,
			CraftUnitPrice: v.CraftPrice,
			CraftUnitName:  v.CraftUnit,
		})
	}
	// Tab YO装订工艺
	response.Product.TabCrafts = make(map[string][]Craft)
	yoTabCrafts := product_price_m.GetCraftByCraftName(model, YOTabCrafts)
	for _, v := range yoTabCrafts {
		response.Product.TabCrafts[YOBindStyle] = append(response.Product.TabCrafts[YOBindStyle], Craft{
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
