package Parsing

import (
	"backend/Validation"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
)

var Shipments = &shipments{}

func GetShip() []ShipmentsData { // Easy access func
	return Shipments.Get(Shipments).([]ShipmentsData)
}

type shipments struct {
	Sheet
	EmptyCollection []ShipmentsData // Shadow this to the correct type when using
}

func (s *shipments) Init() {
	s.Range = "'Shipments'!A:N"
	s.SpreadsheetID = INVMANGMENT
	s.EmptyData = &shipmentsStruct{}
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

	data := s.EmptyCollection //EmptyCollection
	copier.Copy(&data, &s.EmptyCollection)
	for i, e := range Sheetdata {
		if i == 1 {
			continue
		}
		var newData Line
		copier.Copy(&newData, &s.EmptyData)
		newData.processNew(i, e, newData)
		newData.appendNew(&data)
	}
	s.AllData = data
}

type shipmentsStruct struct {
	SheetParseBase // Base struct for Line method inheritance
	ShipmentsData
}

type ShipmentsData struct { // Must be Exported
	Sku     string
	Qty     float64
	Staged  time.Time
	Shipped time.Time
	Alloted time.Time
}

func (data *shipmentsStruct) convData(line []interface{}) {
	const (
		skusCol = 0
		staged  = 11
		shipped = 12
		alloted = 10
	)
	data.Sku = Validation.ConvStringUpper(line[skusCol])
	data.Qty = 1 // default Qty to 1
	data.Alloted = Validation.ConvDate(line[alloted])
	data.Staged = Validation.ConvDate(line[staged])
	data.Shipped = Validation.ConvDate(line[shipped])

}

// Adds new Data
func (new *shipmentsStruct) appendNew(data interface{}) {
	for _, ship := range strings.Split(new.ShipmentsData.Sku, ",") {
		newLine := new.ShipmentsData
		copier.Copy(&newLine, new.ShipmentsData)
		for split, info := range strings.Split(ship, "*") {
			if split > 1 {
				// panic(fmt.Errorf("%s", "Incorrect SKU formating"))
				fmt.Println("%s", "Incorrect SKU formating")
				continue
			}
			_, err := strconv.ParseFloat(info, 64)
			if err != nil {
				newLine.Sku = Validation.ConvStringUpper(info)
			} else {
				newLine.Qty = Validation.ConvNumPos(info)
			}
		}
		if newLine.Sku != "" && newLine.Qty != 0 {
			*data.(*[]ShipmentsData) = append(*data.(*[]ShipmentsData), newLine)
			fmt.Println(newLine)
		}
	}
	// *data.(*[]ShipmentsData) = append(*data.(*[]ShipmentsData), new.ShipmentsData)

}

// Reject data that doesn't make sense
func (data *shipmentsStruct) rejectData() {
	if data.Sku == "" {
		panic(fmt.Errorf("%s", "nil sku"))
	}
	if data.Staged == data.Alloted && data.Staged == data.Shipped {
		panic(fmt.Errorf("%s", "Dates Invarient"))
	}
	// panic on bad values
}
