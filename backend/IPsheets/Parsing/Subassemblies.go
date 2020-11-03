package Parsing

import (
	"backend/Validation"
	"fmt"
	"regexp"
	"strings"
)

var Subs = &subs{} // Create variable for reference
func GetSubs() map[string][]SubsData { // Easy access func
	return Subs.Get(Subs).(map[string][]SubsData)
}

type subs struct { // Create new sheet type
	Sheet
	// EmptyCollection map[string][]SubsData // Shadow this to the correct type
}

func (s *subs) Init() { // Initialize sheet specific data
	s.Range = "'SubAssem Builder'!A:I"
	s.SpreadsheetID = INVMANGMENT
	// s.EmptyData = &subsStruct{} // Empty struct of Data going into the collection
	// s.EmptyCollection = make(map[string][]SubsData)
}

func (s *subs) Parse() { // Change type to correct sheet struct
	Sheetdata := s.getSheet()

	collection := make(map[string][]SubsData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(subsStruct) //EmptyData
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
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
	// Must be parsed in column order else could return before evaluation
	data.Parent = Validation.ConvStringUpper(line[parent])
	data.Child = Validation.ConvStringUpper(line[child])
	data.Qty = Validation.ConvNum(line[qty])
	data.Location = Validation.ConvString(line[loc])
}

func (new *subsStruct) appendNew(data interface{}) { // Adds new Data
	temp := data.(map[string][]SubsData)
	// temp[new.Parent] = SubsData{new.Desc, new.CountType}
	// *data.(*[]SubsData) = append(*data.(*[]SubsData), new.SubsData)

	if temp[new.Parent] == nil {
		temp[new.Parent] = make([]SubsData, 0) //{}
	}
	// Empty Map Keys default to types nil val
	// fmt.Println(new)
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

// Wrapper for recursive BOM, returns a map of parents parts and their Qtys
func CreateBOM(parent string, OnlyImportant bool, depth int) (BOM map[string]float64) {
	BOM = make(map[string]float64)
	recursiveBOM(parent, 1, OnlyImportant, depth, 0, BOM)
	return
}

//  Holy instantaneous
func recursiveBOM(parent string, multiple float64, OnlyImportant bool, depth int, currentDepth int, BOM map[string]float64) {
	sheetData := GetSubs()
	mpl := GetMpl()

	children, hasChildren := sheetData[parent]
	if !hasChildren || (depth != 0 && currentDepth >= depth) || (OnlyImportant && mpl[parent].CountType) {
		BOM[parent] += multiple
	} else {
		for _, line := range children {
			recursiveBOM(line.Child, line.Qty*multiple, OnlyImportant, depth, currentDepth+1, BOM)
		}
	}
}

// Wrapper for recursive BOM, returns a map of parents parts and their Qtys
func CreateFlaggedBOM(parent string, OnlyImportant bool, depth int) (BOM map[string]float64) {
	BOM = make(map[string]float64)
	r := regexp.MustCompile(`[A-Z]+[0-9]+`)
	matches := r.FindAllString(Validation.ConvStringUpper(parent), -1)
	fmt.Println(matches)
	recursiveFlaggedBOM(matches[0], matches[1:], 1, OnlyImportant, depth, 0, BOM)
	return
}

//  Holy instantaneous
func recursiveFlaggedBOM(parent string, flags []string, multiple float64, OnlyImportant bool, depth int, currentDepth int, BOM map[string]float64) {
	sheetData := GetFalgsubs()
	mpl := GetMpl()

	children, hasChildren := sheetData[parent]
	if !hasChildren || (depth != 0 && currentDepth >= depth) || (OnlyImportant && mpl[parent].CountType) {
		BOM[parent] += multiple
	} else {
		for _, line := range children {
			if flagCompare(flags, line.Flags) {
				recursiveFlaggedBOM(line.Child, flags, line.Qty*multiple, OnlyImportant, depth, currentDepth+1, BOM)
			}
		}
	}
}

func flagCompare(flags []string, toMatch string) bool {
	if toMatch == "" {
		return true
	}
	for _, flag := range flags {
		if strings.Contains(toMatch, flag) {
			return true
		}
	}
	return false
}
