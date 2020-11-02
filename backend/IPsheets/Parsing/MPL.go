package Parsing

import (
	"backend/Validation"

	"github.com/jinzhu/copier"
)

var Mpl = &mpl{} // Create variable for reference
func GetMpl() map[string]MplData { // Easy access func
	return Mpl.Get(Mpl).(map[string]MplData)
}

type mpl struct { // Create new sheet type
	Sheet
	EmptyCollection map[string]MplData // Shadow this to the correct type
}

func (s *mpl) Init() { // Initialize sheet specific data
	s.Range = "'Master Part List'!A:AK"
	s.SpreadsheetID = INVMANGMENT
	s.EmptyData = &mplStruct{} // Empty struct of Data going into the collection
	s.EmptyCollection = make(map[string]MplData)
}

func (s *mpl) Parse() { // Change type to correct sheet struct
	Sheetdata := s.getSheet()

	data := s.EmptyCollection //EmptyCollection
	copier.Copy(&data, &s.EmptyCollection)
	for i, e := range Sheetdata {
		var newData Line
		copier.Copy(&newData, &s.EmptyData)
		newData.processNew(i, e, newData)
		newData.appendNew(&data)
	}
	s.AllData = data
}

type mplStruct struct { // Struct that inherits base and data structs
	SheetParseBase // Base struct for Line method inheritance
	MplData
	SKU string
}

// Data from the mpl that is important Must Export
type MplData struct {
	// SKU                     string
	Desc string
	// Supplier                string
	// SupplierPN              string
	// OrderType               string
	// Units string
	// CostPerUnit             float32
	// Currency                string
	// CostPerEaCAD            float32
	// ShippingCostPerUnitCAD  float32
	// Price2                  float32
	// SellMarkup              float32
	// SellPricePerEaUSD       float32
	// HSCode                  float32
	// CountryOfOrigin         string
	// ReorderQtySupplierUnits float32
	// PartLocation            string
	CountType bool
	// ShelfNum                int8
	// SpecSheet               string
	// UnitForSale             string
	// SellPricePerEaUSDStatic float32
	// CostPerUnitCAD          float32
	// LeadTime                uint8
	// DefaultDoorway          string
	// MinShelfQtySupplierUnit uint16
	// OrderByAdjustment       int8
	// DynamicPartClass        uint16
}

func (data *mplStruct) convData(line []interface{}) { // Converts interfaces{} to struct values
	const skuCol = 0
	const descCol = 1
	const countCol = 21
	data.SKU = Validation.ConvStringUpper(line[skuCol])
	data.Desc = Validation.ConvString(line[descCol])
	data.CountType = Validation.ConvBool(line[countCol])
	// data[SKU] = Data{desc, count}
}

func (new *mplStruct) appendNew(data interface{}) { // Adds new Data
	temp := *data.(*map[string]MplData)
	temp[new.SKU] = MplData{new.Desc, new.CountType}
	// *data.(*[]MplData) = append(*data.(*[]MplData), new.MplData)
}
