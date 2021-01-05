package Parsing

import (
	"backend/Validation"
)

var Vendors = &vendors{}

func GetVendors() map[string]vendorsData { // Easy access func
	return Vendors.Get(Vendors).(map[string]vendorsData)
}

func GetVendor(vendor string) vendorsData { // Easy access func
	return GetVendors()[vendor]
}

func CheckVendors() (errs []error) {
	_ = GetVendors()
	return Vendors.Errors
}

type vendors struct {
	Sheet
	// EmptyCollection map[string]vendorsData // Shadow this to the correct type when using
}

func (s *vendors) Init() {
	s.Range = "'Vendor Master Sheet'!A:S"
	s.SpreadsheetID = INVMANGMENT
	// s.EmptyData = &vendorsStruct{}
}

func (s *vendors) Parse() {
	Sheetdata := s.getSheet()

	collection := make(map[string]vendorsData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(vendorsStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type vendorsStruct struct {
	SheetParseBase // Base struct for Line method inheritance
	vendorsData
	Name string
}

type vendorsData struct { // Must be Exported
	Currency     string
	FreeShipping bool
	Leadtime     float64
}

func (data *vendorsStruct) convData(line []interface{}) {
	const (
		nameCol         = 0
		CurrencyCol     = 13
		FreeShippingCol = 15
		LeadtimeCol     = 17
	)
	data.Name = Validation.ConvString(line[nameCol])
	data.Currency = Validation.ConvString(line[CurrencyCol])
	data.FreeShipping = Validation.ConvBool(line[FreeShippingCol])
	data.Leadtime = Validation.ConvNum(line[LeadtimeCol])
}

// Adds new Data
func (new *vendorsStruct) appendNew(data interface{}) {
	data.(map[string]vendorsData)[new.Name] = vendorsData{
		new.Currency,
		new.FreeShipping,
		new.Leadtime}
}
