package Parsing

import (
	"backend/Validation"
)

var Falgsubs = &flagsubs{} // Create variable for reference
func GetFalgsubs() map[string][]flagsubsData { // Easy access func
	return Falgsubs.Get(Falgsubs).(map[string][]flagsubsData)
}

type flagsubs struct { // Create new sheet type
	Sheet
	// EmptyCollection map[string][]flagsubsData // Shadow this to the correct type
}

func (s *flagsubs) Init() { // Initialize sheet specific data
	s.Range = "'Built Option 1'!A1:E9"
	s.SpreadsheetID = "1E39oBmjMDgToLJvQ034OrRF6c_OaWDrkVv-mLLvR-MM"
	// s.EmptyData = &flagsubsStruct{} // Empty struct of Data going into the collection
	// s.EmptyCollection = make(map[string][]flagsubsData)
}

func (s *flagsubs) Parse() { // Change type to correct sheet struct
	Sheetdata := s.getSheet()

	collection := make(map[string][]flagsubsData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(flagsubsStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type flagsubsStruct struct { // Struct that inherits base and data structs
	SheetParseBase // Base struct for Line method inheritance
	flagsubsData
	Parent string
}

// Data from the flagsubs that is important Must Export
type flagsubsData struct {
	Child    string
	Flags    string
	Qty      float64
	Location string
}

func (data *flagsubsStruct) convData(line []interface{}) { // Converts interfaces{} to struct values
	const (
		parent = 0
		flags  = 2
		child  = 3
		qty    = 4
		// loc    = 8
	)
	// Must be parsed in column order else could return before evaluation
	data.Parent = Validation.ConvStringUpper(line[parent])
	data.Child = Validation.ConvStringUpper(line[child])
	data.Qty = Validation.ConvNum(line[qty])
	data.Flags = Validation.ConvStringUpperNilable(line[flags])
}

func (new *flagsubsStruct) appendNew(data interface{}) { // Adds new Data
	temp := data.(map[string][]flagsubsData)
	// temp[new.Parent] = flagsubsData{new.Desc, new.CountType}
	// *data.(*[]flagsubsData) = append(*data.(*[]flagsubsData), new.flagsubsData)

	if temp[new.Parent] == nil {
		temp[new.Parent] = make([]flagsubsData, 0) //{}
	}
	// Empty Map Keys default to types nil val
	// fmt.Println(new)
	temp[new.Parent] = append(temp[new.Parent], new.flagsubsData)
}
