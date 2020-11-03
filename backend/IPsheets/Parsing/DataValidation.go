package Parsing

import (
	"backend/Validation"
)

var DataValidation = &dataVal{}

func GetValidation() []DataValidationData { // Easy access func
	return DataValidation.Get(DataValidation).([]DataValidationData)
}

type dataVal struct {
	Sheet
	// EmptyCollection []DataValidationData // Shadow this to the correct type when using
}

func (s *dataVal) Init() {
	s.Range = "'Data Validation'!B:D"
	s.SpreadsheetID = INVMANGMENT
	// s.EmptyData = &dataValidationStruct{}
}

func (s *dataVal) Parse() {
	Sheetdata := s.getSheet()

	collection := new([]DataValidationData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(dataValidationStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type dataValidationStruct struct {
	SheetParseBase // Base struct for Line method inheritance
	DataValidationData
}

type DataValidationData struct { // Must be Exported
	Unit       string
	Multiplier float64
	Type       string
}

func (data *dataValidationStruct) convData(line []interface{}) {
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
func (new *dataValidationStruct) appendNew(data interface{}) {
	*data.(*[]DataValidationData) = append(*data.(*[]DataValidationData), new.DataValidationData)
}
