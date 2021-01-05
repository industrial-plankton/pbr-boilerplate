package Parsing

import (
	"backend/Validation"
)

var Mpl = &mpl{} // Create variable for reference
func GetMpl() map[string]MplData { // Easy access func
	return Mpl.Get(Mpl).(map[string]MplData)
}

func GetPart(sku string) MplData { // Easy access func
	return GetMpl()[sku]
}

func CheckMpl() (errs []error) {
	_ = GetMpl()
	return Mpl.Errors // Currently there aren't any possible errors
}

type mpl struct { // Create new sheet type
	Sheet
	// EmptyCollection map[string]MplData // Shadow this to the correct type
}

func (s *mpl) Init() { // Initialize sheet specific data
	s.Range = "'Master Part List'!A:AK"
	s.SpreadsheetID = INVMANGMENT
	// s.EmptyData = &mplStruct{} // Empty struct of Data going into the collection
	// s.EmptyCollection = make(map[string]MplData)
}

func (s *mpl) Parse() { // Change type to correct sheet struct
	Sheetdata := s.getSheet()

	collection := make(map[string]MplData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(mplStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type mplStruct struct { // Struct that inherits base and data structs
	SheetParseBase // Base struct for Line method inheritance
	MplData
	SKU string
}

// Data from the mpl that is important Must Export
type MplData struct {
	// SKU                     string
	Desc     string
	Supplier string
	// SupplierPN              string
	// OrderType               string
	Units       string
	CostPerUnit float64
	// Currency                string
	// CostPerEaCAD            float32
	ShippingCostPerUnitCAD float64
	// Price2                  float32
	// SellMarkup              float32
	// SellPricePerEaUSD       float32
	// HSCode                  float32
	// CountryOfOrigin         string
	ReorderQtySupplierUnits float64
	// PartLocation            string
	CountType bool
	// ShelfNum                int8
	// SpecSheet               string
	// UnitForSale             string
	// SellPricePerEaUSDStatic float32
	// CostPerUnitCAD          float32
	LeadTime float64
	// DefaultDoorway          string
	MinShelfQtySupplierUnit float64
	// OrderByAdjustment       int8
	// DynamicPartClass        uint16
}

func (data *mplStruct) convData(line []interface{}) { // Converts interfaces{} to struct values
	const skuCol = 0
	const descCol = 1
	const supplierCol = 3
	const unitsCol = 9
	const costPerCol = 10 //Make Price calculator for subassemblies
	const ShipCostCol = 13
	const reOrderCol = 19
	const countCol = 21
	const leadCol = 32
	const minQtyCol = 34
	data.SKU = Validation.ConvStringUpper(line[skuCol])
	data.Desc = Validation.ConvString(line[descCol])
	data.Supplier = Validation.ConvStringUpper(line[supplierCol])
	data.Units = Validation.ConvString(line[unitsCol])
	data.CostPerUnit = Validation.ConvNum(line[costPerCol])
	data.ShippingCostPerUnitCAD = Validation.ConvNum(line[ShipCostCol])
	data.ReorderQtySupplierUnits = Validation.ConvNum(line[reOrderCol])
	data.CountType = Validation.ConvBool(line[countCol])
	data.LeadTime = Validation.ConvNum(line[leadCol])
	data.MinShelfQtySupplierUnit = Validation.ConvNumPos(line[minQtyCol])
	// data[SKU] = Data{desc, count}
}

func (new *mplStruct) appendNew(data interface{}) { // Adds new Data
	// temp := data.(map[string]MplData)
	data.(map[string]MplData)[new.SKU] = MplData{
		new.Desc,
		new.Supplier,
		new.Units,
		new.CostPerUnit,
		new.ShippingCostPerUnitCAD,
		new.ReorderQtySupplierUnits,
		new.CountType,
		new.LeadTime,
		new.MinShelfQtySupplierUnit}
	// *data.(*[]MplData) = append(*data.(*[]MplData), new.MplData)
}
