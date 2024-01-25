package product_price_r

import (
	"fmt"
	"frozen-go-cms/_const/enum/product_price_e"
	"frozen-go-cms/common/domain"
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/domain/model/product_price_m"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/tealeg/xlsx"
	"io"
	"os"
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
	HardTabCrafts       = []string{"亮膜", "哑膜", "Tab首页加膜", "啤"}                                          // 硬壳tab页面工艺
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
		CoverCraftXs      []float64 `json:"cover_craft_xs"`    // 上面对应的面积x
		CoverCraftYs      []float64 `json:"cover_craft_ys"`    // 上面对应的面积y
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
		InnerCraftIds     []uint64  `json:"inner_craft_ids"`     // 内页需要用到工艺ids
		InnerCraftUnits   []float64 `json:"inner_craft_units"`   // 上面对应的单价,要求len(InnerCraftIds)==len(InnerCraftUnits)
		InnerCraftXs      []float64 `json:"inner_craft_xs"`      // 上面对应的面积x
		InnerCraftYs      []float64 `json:"inner_craft_ys"`      // 上面对应的面积y
		InnerBindIds      []uint64  `json:"inner_bind_ids"`      // 内页用到装订工艺ids
		InnerBindUnits    []float64 `json:"inner_bind_units"`    // 上面对应的单价,要求len(InnerBindIds)==len(InnerBindUnits)
		InnerBindNums     []int     `json:"inner_bind_nums"`     // 上面对应的数量
		InnerPackageIds   []uint64  `json:"inner_package_ids"`   // 内页用到的包装工艺ids
		InnerPackageUnits []float64 `json:"inner_package_units"` // 上面对应的单价,要求len(InnerPackageIds)==len(InnerPackageUnits)
	} `json:"inner,omitempty"`
	HasTab bool `json:"has_tab"` // 是否要Tab页
	Tab    struct {
		TabPageNum      int       `json:"tab_page_num,omitempty"`
		TabMaterial     string    `json:"tab_material,omitempty"`
		TabMaterialGram int       `json:"tab_material_gram,omitempty"`
		TabColor        uint64    `json:"tab_color,omitempty"`
		TabCraftIds     []uint64  `json:"tab_craft_ids"`   // tab面需要用到的工艺ids
		TabCraftUnits   []float64 `json:"tab_craft_units"` // 上面对应的单价,要求len(TabCraftIds)==len(TabCraftUnits)
		TabCraftXs      []float64 `json:"tab_craft_xs"`    // 上面对应的面积x
		TabCraftYs      []float64 `json:"tab_craft_ys"`    // 上面对应的面积y
	} `json:"tab,omitempty"`
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
	innerBinds := product_price_m.GetCraftByIds(model, req.Inner.InnerBindIds)
	innerPackages := product_price_m.GetCraftByIds(model, req.Inner.InnerPackageIds)
	var innerCraftNames []string
	for _, v := range innerCrafts {
		innerCraftNames = append(innerCraftNames, v.CraftBodyCode)
	}
	for _, v := range innerBinds {
		innerCraftNames = append(innerCraftNames, v.CraftBodyCode)
	}
	for _, v := range innerPackages {
		innerCraftNames = append(innerCraftNames, v.CraftBodyCode)
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
	// 产品报价 final todo
	coverMaterial := product_price_m.GetMaterialByNameGram(model, req.Cover.CoverMaterial, req.Cover.CoverMaterialGram)
	innerMaterial := product_price_m.GetMaterialByNameGram(model, req.Inner.InnerMaterial, req.Inner.InnerMaterialGram)
	var bindingPrice, packagePrice float64
	for _, v := range innerBinds {
		bindingPrice += v.MinSumPrice // todo 这是工艺的价格计算，这个很复杂，需要用面积/个数/单价的，不过有个最低单价
	}
	for _, v := range innerPackages {
		packagePrice += v.MinSumPrice // todo 这是工艺的价格计算，这个很复杂，需要用面积/个数/单价的，不过有个最低单价
	}
	// 工艺价格
	var craftPrice float64
	for _, v := range coverCrafts {
		craftPrice += v.MinSumPrice
	}
	for _, v := range req.Cover.CoverOtherCrafts {
		craftPrice += v.Price
	}
	for _, v := range innerCrafts {
		craftPrice += v.MinSumPrice
	}
	var tabColorPrice, tabMaterialPrice float64
	if req.HasTab {
		tabColorPrice = tabColor.PrintBasePrice
		tabMaterial := product_price_m.GetMaterialByNameGram(model, req.Tab.TabMaterial, req.Tab.TabMaterialGram)
		tabMaterialPrice = tabMaterial.LowPrice * float64(req.Tab.TabPageNum)
		for _, v := range innerCrafts {
			craftPrice += v.MinSumPrice
		}
	}
	priceDetail := AutoPriceDetail{
		CoverColorPrice:    coverColor.PrintBasePrice,  // BasePrice乘以BaseNum?
		CoverMaterialPrice: coverMaterial.LowPrice * 4, // 印刷跟本书有关的。
		InnerColorPrice:    innerColor.PrintBasePrice,
		InnerMaterialPrice: innerMaterial.LowPrice * float64(req.Inner.InnerPageNum),
		TabColorPrice:      tabColorPrice,
		TabMaterialPrice:   tabMaterialPrice,
		BindingPrice:       bindingPrice,
		PackagingPrice:     packagePrice,
		CraftPriceSum:      craftPrice,

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
		resp.ResponseOk(c, response)
	} else {
		templateFile, err := xlsx.OpenFile("template.xlsx")
		if err != nil {
			model.Log.Errorf("Failed to open template file:%v", err)
			return myCtx, err
		}

		// 在C5格子写入数据
		sheet := templateFile.Sheets[0]
		cell := sheet.Cell(4, 2) // C5的索引是(4, 2)
		cell.Value = "Hello, World!"

		// 保存为临时文件
		tempFile := "temp.xlsx"
		err = templateFile.Save(tempFile)
		if err != nil {
			model.Log.Errorf("Failed to save temporary file:%v", err)
			return myCtx, err
		}
		defer os.Remove(tempFile)

		// 设置响应头，告诉浏览器发送的是Excel文件
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=download.xlsx")
		c.Writer.Header().Set("Content-Type", "application/octet-stream")

		// 读取临时文件并发送给客户端
		file, err := os.Open(tempFile)
		if err != nil {
			model.Log.Errorf("Failed to open temporary file:%v", err)
			return myCtx, err
		}
		defer file.Close()

		_, err = io.Copy(c.Writer, file)
		if err != nil {
		}
	}
	return myCtx, nil
}
