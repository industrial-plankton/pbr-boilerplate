package Parsing

import (
	"fmt"
	"strings"

	"backend/IPSheets"
)

const (
	MATTRACK    = "1pdhA4p4n4LbOQCrJgmSDZOzHBtV6mIfF2JUUrtxvGuc"
	INVMANGMENT = "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo"
)

type Sheet struct {
	Range         string // Assign this to the correct Value when using
	SpreadsheetID string // Assign this to the correct Value when using
	AllData       interface{}
	// EmptyCollection []interface{} // Shadow this to the correct type when using
	// EmptyData       Line          // Assign this the correct type when using
	Errors []error
}

type Line interface {
	processNew(i int, e []interface{}, obj Line, collection interface{}, errors *[]error)
	appendNew(interface{}) //* Must Override
	handleError(i int, errors *[]error)
	newData(e []interface{}, obj Line)
	rejectData()
	checkWarnings(line int, obj Line, errors *[]error)
	warningData()
	assumeData(errors *[]error)
	convData(line []interface{}) //* Must Override
}

type SheetParse interface {
	Get(ref SheetParse) interface{}
	Parse() //* Must Override
	Init()  //* Must Override
	getSheet() [][]interface{}
	GetErrors() []error
}

//* Must Override
func (s *Sheet) Init() { // Shadow this to the correct type when using
	s.Range = "2020 Tracking!A:U"
	s.SpreadsheetID = "1pdhA4p4n4LbOQCrJgmSDZOzHBtV6mIfF2JUUrtxvGuc"
	// s.EmptyData = &SheetParseBase{}
}

func (s *Sheet) Get(ref SheetParse) interface{} { // Shadow this to the correct type when using // OR Assert correct type when receiving
	if s.SpreadsheetID == "" {
		ref.Init()
	}
	if s.AllData == nil {
		ref.Parse()
	}
	return s.AllData //.(emptyCollectiontype)
}

//* Must Inherit or implement all Line interface Methods
type SheetParseBase struct {
}

type Part struct {
	Sku string
	Qty float64
}

//* Must Override, should usually just be a copy with the attached type corrected
func (s *Sheet) Parse() {
	// Sheetdata := s.getSheet()
	collection := new([]interface{}) //EmptyCollection
	// for i, e := range Sheetdata {
	// newData := new(Line) //EmptyData
	// newData.processNew(i, e, newData, collection, &s.Errors)
	// newData.appendNew(&data)
	// }

	s.AllData = collection
}

func (s *Sheet) getSheet() [][]interface{} {
	return IPSheets.BatchGet([]string{s.Range}, s.SpreadsheetID, nil)[0]
}

func (s *Sheet) GetErrors() []error {
	return s.Errors
}

//* Must Override
// Adds new Data to the collection
func (new *SheetParseBase) appendNew(data interface{}) {
	// func (new *Data) appendNew(data interface{}) {
	// 	*data.(*[]Data) = append(*data.(*[]Data), *new)
	// }
	*data.(*[]interface{}) = append(*data.(*[]interface{}), *new)
}

//* Must Override
// Converts interfaces into Stuct data
func (data *SheetParseBase) convData(line []interface{}) {
	// Must be parsed in column order else could return before evaluation
	// data.Doorway = Doorway.ToDoorway(line[doorwayCol])
	// data.Sku = Validation.Sku(line[skuCol])
	// data.Qty = Validation.ConvNum(line[qtyCol])
}

func (data *SheetParseBase) processNew(i int, e []interface{}, obj Line, collection interface{}, errors *[]error) {
	defer obj.handleError(i, errors)
	obj.newData(e, obj)
	obj.rejectData()
	obj.checkWarnings(i, obj, errors)
	obj.assumeData(errors)
	// fmt.Println(i+1, obj)
	obj.appendNew(collection)
}

// Formats and Checks new Data struct
func (data *SheetParseBase) newData(newline []interface{}, obj Line) {
	defer func() {
		err := recover()
		if err != nil {
			if strings.Contains(err.(error).Error(), "index") {
				return
			}
			if strings.Contains(err.(error).Error(), "&continue&") { //Print off continue flagged errors
				fmt.Println(err)
			} else {
				// Rethrow
				panic(err)
			}
		}
	}()
	//  Convert Data using Validation conversion functions
	obj.convData(newline)
}

// Handles errors thrown by newData, continues if only minor
func (data *SheetParseBase) handleError(line int, errors *[]error) {
	err := recover()
	if err != nil {
		if !strings.Contains(err.(error).Error(), "&minor&") && !strings.Contains(err.(error).Error(), "index") && line != 0 { //Print off errors that dont contain the &minor& flag, and index errors
			msg := fmt.Errorf("%v %v %v", err.(error).Error(), ", line:", line+1)
			fmt.Println(msg)
			*errors = append(*errors, msg)
		}
	}
}

// Reject data that doesn't make sense
func (data *SheetParseBase) rejectData() {
	// panic on bad values
}

// Check for Warning and handle them
func (data *SheetParseBase) checkWarnings(line int, obj Line, errors *[]error) {
	defer func() {
		err := recover()
		if err != nil {
			msg := fmt.Errorf("%v %v %v", err.(error).Error(), ", line:", line+1)
			fmt.Println(msg)
			*errors = append(*errors, msg)
			// fmt.Println(err.(error))
		}
	}()
	obj.warningData()
}

// Warns of strange data, assumes data if it can
func (data *SheetParseBase) warningData() {
	var sb strings.Builder

	// Write to sb on weird data, adjust values if you can, panic at end
	// if data.Completed && (time.Now().Before(data.RecDate) || data.RecDate.IsZero()) {
	// 	sb.WriteString(fmt.Sprintln("Received true but Received Date invalid:", data.RecDate))
	// 	data.RecDate = time.Now()
	// }

	Estring := strings.TrimSpace(sb.String())
	if !(Estring == "") {
		panic(fmt.Errorf("%s", Estring))
	}
}

// Assume any missing data that you can
func (data *SheetParseBase) assumeData(errors *[]error) {
}
