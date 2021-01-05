package Parsing

import (
	"backend/Validation"
	"fmt"
	"time"
)

var Inv = &inv{} // Create variable for reference
func GetInv() map[string]InvData { // Easy access func
	return Inv.Get(Inv).(map[string]InvData)
}

func GetPartInv(sku string) InvData { // Easy access func
	return GetInv()[sku]
}

func CheckInv() (errs []error) {
	_ = GetInv()
	return Inv.Errors
}

type inv struct { // Create new sheet type
	Sheet
	// EmptyCollection map[string]InvData // Shadow this to the correct type
}

func (s *inv) Init() { // Initialize sheet specific data
	s.Range = "'Physical Inventory Entry'!A:F"
	s.SpreadsheetID = MATTRACK
	// s.EmptyData = &invStruct{} // Empty struct of Data going into the collection
	// s.EmptyCollection = make(map[string]InvData)
}

func (s *inv) Parse() { // Change type to correct sheet struct
	Sheetdata := s.getSheet()

	collection := make(map[string]InvData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(invStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type invStruct struct { // Struct that inherits base and data structs
	SheetParseBase // Base struct for Line method inheritance
	InvData
	SKU string
}

// Data from the inv that is important Must Export
type InvData struct {
	Qty  float64
	Date time.Time
}

func (data *invStruct) convData(line []interface{}) { // Converts interfaces{} to struct values
	const skuCol = 0
	const qtyCol = 3
	const dateCol = 4
	data.SKU = Validation.ConvStringUpper(line[skuCol])
	data.Qty = Validation.ConvNumPos(line[qtyCol])
	data.Date = Validation.ConvDate(line[dateCol])
}

// Reject data that doesn't make sense
func (data *invStruct) rejectData() {
	// panic on bad values
	if _, ok := GetSubs()[data.SKU]; ok && !GetPart(data.SKU).CountType {
		panic(fmt.Errorf("Non-count subassembly %s", data.SKU))
	}
}

func (new *invStruct) appendNew(data interface{}) { // Adds new Data
	// for key, Qty := range CreateBOM(new.SKU, true, 0) { // Split notcount
	// 	if data.(map[string]InvData)[key].Date.Before(new.Date) {
	// 		data.(map[string]InvData)[key] = InvData{
	// 			new.Qty * Qty,
	// 			new.Date}
	// 	}
	// }
	// temp := data.(map[string]InvData)
	if data.(map[string]InvData)[new.SKU].Date.Before(new.Date) {
		data.(map[string]InvData)[new.SKU] = InvData{
			new.Qty,
			new.Date}
	}
	// *data.(*[]InvData) = append(*data.(*[]InvData), new.InvData)
}
