package Parsing

import (
	"backend/Validation"
)

var UnitConvTable = &unitConv{}

func GetUnitConv() *[]UnitConvTableData { // Easy access func
	return UnitConvTable.Get(UnitConvTable).(*[]UnitConvTableData)
}

type unitConv struct {
	Sheet
	// EmptyCollection []UnitConvTableData // Shadow this to the correct type when using
}

func (s *unitConv) Init() {
	s.Range = "'Data Validation'!B:D"
	s.SpreadsheetID = INVMANGMENT
	// s.EmptyData = &UnitConvTableStruct{}
}

func (s *unitConv) Parse() {
	Sheetdata := s.getSheet()

	collection := new([]UnitConvTableData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(UnitConvTableStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type UnitConvTableStruct struct {
	SheetParseBase // Base struct for Line method inheritance
	UnitConvTableData
}

type UnitConvTableData struct { // Must be Exported
	Unit       string
	Multiplier float64
	Type       string
}

func (data *UnitConvTableStruct) convData(line []interface{}) {
	const (
		unitCol = 0
		multCol = 1
		typeCol = 2
	)
	data.Unit = Validation.ConvString(line[unitCol])
	data.Multiplier = Validation.ConvNumPos(line[multCol])
	data.Type = Validation.ConvString(line[typeCol])
}

// Adds new Data
func (new *UnitConvTableStruct) appendNew(data interface{}) {
	*data.(*[]UnitConvTableData) = append(*data.(*[]UnitConvTableData), new.UnitConvTableData)
}
