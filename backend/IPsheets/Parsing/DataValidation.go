package Parsing

import (
	"backend/Validation"

	"github.com/jinzhu/copier"
)

var DataValidation = &dataVal{}

func GetValidation() []DataValidationData { // Easy access func
	return DataValidation.Get(DataValidation).([]DataValidationData)
}

type dataVal struct {
	Sheet
	EmptyCollection []DataValidationData // Shadow this to the correct type when using
}

func (s *dataVal) Init() {
	s.Range = "'Data Validation'!B:D"
	s.SpreadsheetID = INVMANGMENT
	s.EmptyData = &dataValidationStruct{}
}

func (s *dataVal) Parse() {
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

type dataValidationStruct struct {
	SheetParseBase // Base struct for Line method inheritance
	DataValidationData
}

type DataValidationData struct { // Must be Exported
	One   string
	Two   string
	Three string
}

func (data *dataValidationStruct) convData(line []interface{}) {
	data.One = Validation.ConvString(line[0])
	data.Two = Validation.ConvString(line[1])
	data.Three = Validation.ConvString(line[2])
}

// Adds new Data
func (new *dataValidationStruct) appendNew(data interface{}) {
	*data.(*[]DataValidationData) = append(*data.(*[]DataValidationData), new.DataValidationData)
}
