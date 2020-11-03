package Parsing

// import (
// 	"backend/Validation"

// 	"github.com/jinzhu/copier"
// )

// var DataValidation = &dataVal{} // Create variable for reference

// type dataVal struct { // Create new sheet type
// 	Sheet
// 	EmptyCollection []DataValidationData // Shadow this to the correct type
// }

// func (s *dataVal) Init() { // Initialize sheet specific data
// 	s.Range = "'Data Validation'!B:D"
// 	s.SpreadsheetID = "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo"
// 	s.EmptyData = &dataValidationStruct{} // Empty struct of Data going into the collection
// }

// func (s *dataVal) Parse() { // Change type to correct sheet struct
// 	Sheetdata := s.getSheet()

// 	collection := s.EmptyCollection //EmptyCollection
// 	copier.Copy(&collection, &s.EmptyCollection)
// 	for i, e := range Sheetdata {
// 		var newData Line
// 		copier.Copy(&newData, &s.EmptyData)
// 		newData.processNew(i, e, newData, &collection, &s.Errors)
// 	}
// 	s.AllData = collection
// }

// type dataValidationStruct struct { // Struct that inherits base and data structs
// 	SheetParseBase // Base struct for Line method inheritance
// 	DataValidationData
// }

// type DataValidationData struct { // Data going into the collection, Must be Exported
// 	One   string
// 	Two   string
// 	Three string
// }

// func (data *dataValidationStruct) convData(line []interface{}) { // Converts interfaces{} to struct values
// Must be parsed in column order else could return before evaluation
// 	data.One = Validation.ConvString(line[0])
// 	data.Two = Validation.ConvString(line[1])
// 	data.Three = Validation.ConvString(line[2])
// }

// func (new *dataValidationStruct) appendNew(data interface{}) { // Adds new Data
// 	*data.(*[]DataValidationData) = append(*data.(*[]DataValidationData), new.DataValidationData)
// }
