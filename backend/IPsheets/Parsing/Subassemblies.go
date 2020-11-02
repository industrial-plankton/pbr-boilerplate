package Parsing

import (
	"backend/Validation"

	"github.com/jinzhu/copier"
)

var Subs = &subs{} // Create variable for reference
func GetSubs() map[string][]SubsData { // Easy access func
	return Subs.Get(Subs).(map[string][]SubsData)
}

type subs struct { // Create new sheet type
	Sheet
	EmptyCollection map[string][]SubsData // Shadow this to the correct type
}

func (s *subs) Init() { // Initialize sheet specific data
	s.Range = "'SubAssem Builder'!A:I"
	s.SpreadsheetID = INVMANGMENT
	s.EmptyData = &subsStruct{} // Empty struct of Data going into the collection
	s.EmptyCollection = make(map[string][]SubsData)
}

func (s *subs) Parse() { // Change type to correct sheet struct
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

type subsStruct struct { // Struct that inherits base and data structs
	SheetParseBase // Base struct for Line method inheritance
	SubsData
	Parent string
}

// Data from the subs that is important Must Export
type SubsData struct {
	Child    string
	Qty      float64
	Location string
}

func (data *subsStruct) convData(line []interface{}) { // Converts interfaces{} to struct values
	const (
		parent = 0
		child  = 1
		qty    = 2
		loc    = 8
	)
	data.Parent = Validation.ConvStringUpper(line[parent])
	data.Child = Validation.ConvStringUpper(line[child])
	data.Location = Validation.ConvString(line[loc])
	data.Qty = Validation.ConvNum(line[qty])
	// data[SKU] = Data{desc, count}
}

func (new *subsStruct) appendNew(data interface{}) { // Adds new Data
	temp := *data.(*map[string][]SubsData)
	// temp[new.Parent] = SubsData{new.Desc, new.CountType}
	// *data.(*[]SubsData) = append(*data.(*[]SubsData), new.SubsData)

	if temp[new.Parent] == nil {
		temp[new.Parent] = []SubsData{}
	}
	// Empty Map Keys default to types nil val
	temp[new.Parent] = append(temp[new.Parent], new.SubsData)
}

//
//

func FindOffspring(parent string, OnlyImportant bool) map[string]float64 {
	sheetData := GetSubs()
	mpl := GetMpl()

	lostChildren := make(map[string]float64)

	children, hasChildren := sheetData[parent]
	if hasChildren {
		for _, line := range children {
			_, alsoHasChildren := sheetData[line.Child]
			if alsoHasChildren && (!OnlyImportant || !mpl[line.Child].CountType) {
				grandkids := FindOffspring(line.Child, OnlyImportant)
				for grandchild, amount := range grandkids {
					lostChildren[grandchild] += amount * line.Qty
				}
			} else {
				lostChildren[line.Child] += line.Qty
			}
		}
	} else {
		lostChildren[parent] = 1
	}

	return lostChildren
}
