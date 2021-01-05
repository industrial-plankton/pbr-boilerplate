package Parsing

import (
	"backend/Validation"
)

var CurrencyConvTable = &currencyConv{}

func GetCurrencyConv() *[]CurrencyConvTableData { // Easy access func
	return CurrencyConvTable.Get(CurrencyConvTable).(*[]CurrencyConvTableData)
}

type currencyConv struct {
	Sheet
	// EmptyCollection []CurrencyConvTableData // Shadow this to the correct type when using
}

func (s *currencyConv) Init() {
	s.Range = "'Data Validation'!I:J"
	s.SpreadsheetID = INVMANGMENT
	// s.EmptyData = &currencyConvTableStruct{}
}

func (s *currencyConv) Parse() {
	Sheetdata := s.getSheet()

	collection := new([]CurrencyConvTableData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(currencyConvTableStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type currencyConvTableStruct struct {
	SheetParseBase // Base struct for Line method inheritance
	CurrencyConvTableData
}

type CurrencyConvTableData struct { // Must be Exported
	Currency   string
	Multiplier float64
}

func (data *currencyConvTableStruct) convData(line []interface{}) {
	const (
		currencyCol = 0
		multCol     = 1
	)
	data.Currency = Validation.ConvString(line[currencyCol])
	data.Multiplier = Validation.ConvNumPos(line[multCol])
}

// Adds new Data
func (new *currencyConvTableStruct) appendNew(data interface{}) {
	*data.(*[]CurrencyConvTableData) = append(*data.(*[]CurrencyConvTableData), new.CurrencyConvTableData)
}
