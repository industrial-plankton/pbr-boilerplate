package Parsing

import (
	"backend/Validation"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var Shipments = &shipments{}

func GetShip() []ShipmentsData { // Easy access func
	return Shipments.Get(Shipments).([]ShipmentsData)
}

type shipments struct {
	Sheet
	// EmptyCollection []ShipmentsData // Shadow this to the correct type when using
}

func (s *shipments) Init() {
	s.Range = "'Shipments'!A:N"
	s.SpreadsheetID = INVMANGMENT
	// s.EmptyData = &shipmentsStruct{}
}

func (s *shipments) Get(ref SheetParse) interface{} { // Shadow this to the correct type when using // OR Assert correct type when receiving
	if s.SpreadsheetID == "" {
		ref.Init()
	}
	if s.AllData == nil {
		ref.Parse()
	}
	return s.AllData //.(emptyCollectiontype)
}

func (s *shipments) Parse() {
	Sheetdata := s.getSheet()

	collection := new([]ShipmentsData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(shipmentsStruct)
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type shipmentsStruct struct {
	SheetParseBase // Base struct for Line method inheritance
	ShipmentsData
}

type ShipmentsData struct { // Must be Exported
	Parts   []Part
	Staged  time.Time
	Shipped time.Time
	Alloted time.Time
}

func (data *shipmentsStruct) convData(line []interface{}) {
	const (
		skusCol = 0
		alloted = 10
		staged  = 11
		shipped = 12
	)
	data.Alloted = Validation.ConvDate(line[alloted])
	data.Staged = Validation.ConvDate(line[staged])
	data.Shipped = Validation.ConvDate(line[shipped])

	parts := strings.Split(Validation.ConvStringUpper(line[skusCol]), ",")
	data.Parts = make([]Part, len(parts))
	for i, ship := range parts {
		for split, info := range strings.Split(ship, "*") {
			if split > 1 {
				panic(fmt.Errorf("%s", "Incorrect SKU formating"))
			}
			_, err := strconv.ParseFloat(info, 64)
			if err != nil {
				data.Parts[i].Sku = Validation.ConvStringUpper(info)
			} else {
				data.Parts[i].Qty = Validation.ConvNumPos(info)
			}
		}
	}

}

// Adds new Data
func (new *shipmentsStruct) appendNew(data interface{}) {
	*data.(*[]ShipmentsData) = append(*data.(*[]ShipmentsData), new.ShipmentsData)
}

// Reject data that doesn't make sense
func (data *shipmentsStruct) rejectData() {
	// panic on bad values
	if len(data.Parts) == 0 {
		panic(fmt.Errorf("%s", "no parts"))
	}

	if data.Staged.Equal(data.Alloted) && data.Staged.Equal(data.Shipped) {
		panic(fmt.Errorf("%s", "Bad Dates"))
	}

	for _, part := range data.Parts {
		if part.Sku == "" {
			panic(fmt.Errorf("%s", "nil parts"))
		}
	}
}

// Assume any missing data that you can
func (data *shipmentsStruct) assumeData(errors *[]error) {
	for i, part := range data.Parts {
		if part.Qty == 0 {
			data.Parts[i].Qty = 1
		}
	}

	// Assume ship 2 years from now if no date
	if data.Shipped.IsZero() {
		data.Shipped = time.Now().AddDate(2, 0, 0)
	}
}
