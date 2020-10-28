package mpl

import (
	"backend/IPSheets"
	"backend/Validation"
)

const (
	SpreadsheetID = "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo"
	Range         = "'Master Part List'!A:AK"
	// Zero indexed Columns
	skuCol   = 0
	descCol  = 1
	countCol = 21
)

var allData map[string]Data

func Get() map[string]Data {
	if allData == nil {
		Refresh()
	}
	return allData
}

func Refresh() {
	allData = parse()
}

// Data from the MPL that is important
type Data struct {
	// SKU                     string
	Desc string
	// Supplier                string
	// SupplierPN              string
	// OrderType               string
	// Units                   string
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

/// parse maps the interface to the mpl Data struct
func parse() map[string]Data {
	Sheetdata := IPSheets.BatchGet([]string{Range}, SpreadsheetID, nil)[0]

	data := make(map[string]Data)
	for _, e := range Sheetdata {
		// err := goodSub(e)
		// if err != nil {
		// 	fmt.Println(err, ", line:", i+1)
		// 	continue
		// }
		SKU := Validation.ConvStringUpper(e[skuCol])
		desc := Validation.ConvString(e[descCol])
		count := Validation.ConvBool(e[countCol])
		data[SKU] = Data{desc, count}
	}

	return data
}

// func countToBool(count string) bool {
// 	count = strings.ToUpper(strings.TrimSpace(count))
// 	if count == "COUNT" {
// 		return true
// 	} else {
// 		return false
// 	}
// }
