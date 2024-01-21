package product_price_e

type CraftUnitType int
type SizeConfigType int

const (
	// 工艺单位 0:无 1:元/张 2:元/本 3:元/次 4:元/m² 5:元/cm²
	CraftUnitTypeEmpty    CraftUnitType = 0
	CraftUnitTypePerPage  CraftUnitType = 1
	CraftUnitTypePerBook  CraftUnitType = 2
	CraftUnitTypePerPiece CraftUnitType = 3
	CraftUnitTypeMSqr     CraftUnitType = 4
	CraftUnitTypePerCmSqr CraftUnitType = 5

	// 规格尺寸类型 1:书本 2:天地盖盒子 3:卡片
	SizeConfigTypeBook SizeConfigType = 1
	SizeConfigTypeBox  SizeConfigType = 2
	SizeConfigTypeCard SizeConfigType = 3
)
