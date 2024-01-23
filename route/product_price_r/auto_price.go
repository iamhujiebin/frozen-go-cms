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
			CraftBodyName:   v.CraftBodyName,
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

type AutoPriceReq struct {
	Product struct {
		BindStyle    int     `json:"bind_style,omitempty"`
		ClientName   string  `json:"client_name,omitempty"`
		DeliveryTime string  `json:"delivery_time,omitempty"`
		PayExtraDesc string  `json:"pay_extra_desc,omitempty"`
		PayExtraNum  int     `json:"pay_extra_num,omitempty"`
		PayExtraUnit int     `json:"pay_extra_unit,omitempty"`
		PayMethod    string  `json:"pay_method,omitempty"`
		PriceFactor  float64 `json:"price_factor,omitempty"`
		PrintNum     int     `json:"print_num,omitempty"`
		ProductName  string  `json:"product_name,omitempty"`
		Size         int     `json:"size,omitempty"`
		TranDesc     string  `json:"tran_desc,omitempty"`
		TranPrice    int     `json:"tran_price,omitempty"`
	} `json:"product,omitempty"`
	Cover struct {
		CoverColor        int     `json:"cover_color,omitempty"`
		CoverCraftCheck42 bool    `json:"cover_craft_check_42,omitempty"`
		CoverCraftCheck46 bool    `json:"cover_craft_check_46,omitempty"`
		CoverCraftCheck50 bool    `json:"cover_craft_check_50,omitempty"`
		CoverCraftCheck51 bool    `json:"cover_craft_check_51,omitempty"`
		CoverCraftCheck52 bool    `json:"cover_craft_check_52,omitempty"`
		CoverCraftCheck53 bool    `json:"cover_craft_check_53,omitempty"`
		CoverCraftCheck54 bool    `json:"cover_craft_check_54,omitempty"`
		CoverCraftCheck55 bool    `json:"cover_craft_check_55,omitempty"`
		CoverCraftUnit42  int     `json:"cover_craft_unit_42,omitempty"`
		CoverCraftUnit46  float64 `json:"cover_craft_unit_46,omitempty"`
		CoverCraftUnit50  float64 `json:"cover_craft_unit_50,omitempty"`
		CoverCraftUnit51  float64 `json:"cover_craft_unit_51,omitempty"`
		CoverCraftUnit52  float64 `json:"cover_craft_unit_52,omitempty"`
		CoverCraftUnit53  float64 `json:"cover_craft_unit_53,omitempty"`
		CoverCraftUnit54  float64 `json:"cover_craft_unit_54,omitempty"`
		CoverCraftUnit55  float64 `json:"cover_craft_unit_55,omitempty"`
		CoverCraftX54     int     `json:"cover_craft_x_54,omitempty"`
		CoverCraftX55     int     `json:"cover_craft_x_55,omitempty"`
		CoverCraftY54     int     `json:"cover_craft_y_54,omitempty"`
		CoverCraftY55     int     `json:"cover_craft_y_55,omitempty"`
		CoverMaterial     string  `json:"cover_material,omitempty"`
		CoverMaterialsNum int     `json:"cover_materials_num,omitempty"`
		CoverOtherCrafts  []struct {
			Name  string `json:"name,omitempty"`
			Price string `json:"price,omitempty"`
		} `json:"cover_other_crafts,omitempty"`
	} `json:"cover,omitempty"`
	Inner struct {
		InnerBindCheck19  bool    `json:"inner_bind_check_19,omitempty"`
		InnerBindCheck20  bool    `json:"inner_bind_check_20,omitempty"`
		InnerBindCheck21  bool    `json:"inner_bind_check_21,omitempty"`
		InnerBindCheck22  bool    `json:"inner_bind_check_22,omitempty"`
		InnerBindCheck23  bool    `json:"inner_bind_check_23,omitempty"`
		InnerBindCheck24  bool    `json:"inner_bind_check_24,omitempty"`
		InnerBindCheck25  bool    `json:"inner_bind_check_25,omitempty"`
		InnerBindCheck28  bool    `json:"inner_bind_check_28,omitempty"`
		InnerBindCheck29  bool    `json:"inner_bind_check_29,omitempty"`
		InnerBindUnit19   int     `json:"inner_bind_unit_19,omitempty"`
		InnerBindUnit20   int     `json:"inner_bind_unit_20,omitempty"`
		InnerBindUnit21   int     `json:"inner_bind_unit_21,omitempty"`
		InnerBindUnit22   int     `json:"inner_bind_unit_22,omitempty"`
		InnerBindUnit23   int     `json:"inner_bind_unit_23,omitempty"`
		InnerBindUnit24   int     `json:"inner_bind_unit_24,omitempty"`
		InnerBindUnit25   int     `json:"inner_bind_unit_25,omitempty"`
		InnerBindUnit28   float64 `json:"inner_bind_unit_28,omitempty"`
		InnerBindUnit29   float64 `json:"inner_bind_unit_29,omitempty"`
		InnerColor        int     `json:"inner_color,omitempty"`
		InnerCraftCheck26 bool    `json:"inner_craft_check_26,omitempty"`
		InnerCraftCheck27 bool    `json:"inner_craft_check_27,omitempty"`
		InnerCraftCheck35 bool    `json:"inner_craft_check_35,omitempty"`
		InnerCraftCheck36 bool    `json:"inner_craft_check_36,omitempty"`
		InnerCraftCheck54 bool    `json:"inner_craft_check_54,omitempty"`
		InnerCraftCheck55 any     `json:"inner_craft_check_55,omitempty"`
		InnerCraftUnit26  float64 `json:"inner_craft_unit_26,omitempty"`
		InnerCraftUnit27  float64 `json:"inner_craft_unit_27,omitempty"`
		InnerCraftUnit35  float64 `json:"inner_craft_unit_35,omitempty"`
		InnerCraftUnit36  float64 `json:"inner_craft_unit_36,omitempty"`
		InnerCraftUnit54  float64 `json:"inner_craft_unit_54,omitempty"`
		InnerCraftUnit55  float64 `json:"inner_craft_unit_55,omitempty"`
		InnerCraftX54     any     `json:"inner_craft_x_54,omitempty"`
		InnerCraftX55     any     `json:"inner_craft_x_55,omitempty"`
		InnerCraftY54     any     `json:"inner_craft_y_54,omitempty"`
		InnerCraftY55     any     `json:"inner_craft_y_55,omitempty"`
		InnerMaterial     string  `json:"inner_material,omitempty"`
		InnerMaterialsNum int     `json:"inner_materials_num,omitempty"`
		InnerPageNum      int     `json:"inner_page_num,omitempty"`
	} `json:"inner,omitempty"`
	Tab struct {
		TabPageNum      int     `json:"tab_page_num,omitempty"`
		TabColor        int     `json:"tab_color,omitempty"`
		TabCraftCheck35 bool    `json:"tab_craft_check_35,omitempty"`
		TabCraftCheck36 bool    `json:"tab_craft_check_36,omitempty"`
		TabCraftCheck52 bool    `json:"tab_craft_check_52,omitempty"`
		TabCraftCheck53 bool    `json:"tab_craft_check_53,omitempty"`
		TabCraftCheck54 bool    `json:"tab_craft_check_54,omitempty"`
		TabCraftCheck55 bool    `json:"tab_craft_check_55,omitempty"`
		TabCraftUnit35  float64 `json:"tab_craft_unit_35,omitempty"`
		TabCraftUnit36  float64 `json:"tab_craft_unit_36,omitempty"`
		TabCraftUnit52  float64 `json:"tab_craft_unit_52,omitempty"`
		TabCraftUnit53  int     `json:"tab_craft_unit_53,omitempty"`
		TabCraftUnit54  float64 `json:"tab_craft_unit_54,omitempty"`
		TabCraftUnit55  float64 `json:"tab_craft_unit_55,omitempty"`
		TabCraftX54     int     `json:"tab_craft_x_54,omitempty"`
		TabCraftX55     int     `json:"tab_craft_x_55,omitempty"`
		TabCraftY54     int     `json:"tab_craft_y_54,omitempty"`
		TabCraftY55     int     `json:"tab_craft_y_55,omitempty"`
		TabMaterial     string  `json:"tab_material,omitempty"`
		TabMaterialsNum int     `json:"tab_materials_num,omitempty"`
	} `json:"tab,omitempty"`
}

type AutoPriceResponse struct {
	ChMessage string `json:"chMessage"` // 中文信息
	EnMessage string `json:"enMessage"` // 英文信息
}

// @Tags 报价系统
// @Summary 自动报价
// @Param Authorization header string true "token"
// @Success 200 {object} AutoPriceResponse
// @Router /v1_0/productPrice/auto/generate [post]
func AutoPriceGenerate(c *gin.Context) (*mycontext.MyContext, error) {
	myCtx := mycontext.CreateMyContext(c.Keys)
	//model := domain.CreateModelContext(myCtx)
	var req AutoPriceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		//return myCtx, err
	}
	var response AutoPriceResponse
	response.ChMessage = `
封面印刷：300.00元
封面材料：2903.15元
内页印刷：300.00元
内页材料：108.67元
装订价格：2500.00元
工艺费用 ：534.80元
包装要求 ：400.00元
 – – – – – – – – – – – – – – – – –  
生产成本 ：7046.62元

额外成本：0元
运输成本：0.00元
 – – – – – – – – – – – – – – – – – – – – – – – – – – – – – – – –
费用合计 ：7046.62元
`
	response.EnMessage = `
Size:140*210mm 
Cover: 4P 270g PU   Leather  4C+0C  Matt Lamination   +Grey Board 
Inside page: 4P 210g Gloss Art  Paper  C1S  4C+4C     Poly Bag   Cartoning 
Wire bound
 By Ali Assurance    Price    Total: USD1371.76
`
	resp.ResponseOk(c, response)
	return myCtx, nil
}
