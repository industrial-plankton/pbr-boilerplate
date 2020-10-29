package IPSheets

import (
	"fmt"
	"strings"
)

const (
	SpreadsheetID = "1pdhA4p4n4LbOQCrJgmSDZOzHBtV6mIfF2JUUrtxvGuc"
	Range         = "2020 Tracking!A:U"
	// Column Consts
	skuCol           = 5
	qtyCol           = 6
	recQtyCol        = 17
	doorwayCol       = 3
	fullyReceivedCol = 19
	orderCol         = 12
	exCol            = 14
	recCol           = 16
)

var allData []Data
var emptyCollection = []Data{}
var emptyData = Data{}

type Sheet struct {
	Range           string // Assign this to the correct Value when using
	SpreadsheetID   string // Assign this to the correct Value when using
	allData         interface{}
	EmptyCollection interface{} // Shadow this to the correct type when using
	EmptyData       Line        // Shadow this to the correct type when using
}

type Line interface {
	ProcessNew(i int, e []interface{}, obj Line)
	AppendNew(interface{})
	handleError(i int)
	newData(e []interface{}, obj Line)
	rejectData()
	checkWarnings()
	warningData()
	assumeData()
	ConvData(line []interface{})
}

func (s Sheet) Get() interface{} { // Shadow this to the correct type when using
	if s.allData == nil {
		s.Refresh()
	}
	return s.allData //.(emptyCollectiontype)
}

func (s Sheet) Refresh() {
	s.allData = s.parse()
}

// Data from the MPL that is important
type Data struct {
	// Sku       string
	// Qty       float64
	// RecQty    float64
	// Doorway byte
	// Completed bool
	// OrderDate time.Time
	// ExpDate   time.Time
	// RecDate   time.Time
}

// parse maps the interface to the mpl Data struct
func (s Sheet) parse() interface{} {
	Sheetdata := BatchGet([]string{s.Range}, s.SpreadsheetID, nil)[0]

	data := s.EmptyCollection
	for i, e := range Sheetdata {
		newData := s.EmptyData
		newData.ProcessNew(i, e, newData)
		newData.AppendNew(&data)
	}

	return data
}

func (data *Data) ProcessNew(i int, e []interface{}, obj Line) {
	// var line Line = data

	defer obj.handleError(i)
	obj.newData(e, obj)
	obj.rejectData()
	obj.checkWarnings()
	obj.assumeData()
}

// Handles errors thrown by newData, continues if only minor
func (data *Data) handleError(line int) {
	err := recover()
	if err != nil {
		if !strings.Contains(err.(error).Error(), "&minor&") && !strings.Contains(err.(error).Error(), "index") && line != 0 { //Print off errors that dont contain the &minor& flag, and index errors
			fmt.Println(err.(error).Error(), ", line:", line+1)
		}
	}
}

// Adds new Data
func (new *Data) AppendNew(data interface{}) {
	// var line line = new
	data = append(*data.(*[]Data), *new)

}

// Formats and Checks new Data struct
func (data *Data) newData(newline []interface{}, obj Line) {
	// var ref Line = data
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
	obj.ConvData(newline)
	return
}

func (data *Data) ConvData(line []interface{}) {
	// data.Doorway = Doorway.ToDoorway(line[doorwayCol])
	// data.Sku = Validation.Sku(line[skuCol])
	// data.Qty = Validation.ConvNum(line[qtyCol])
}

// Reject data that doesn't make sense
func (data *Data) rejectData() {
	// panic on bad values
}

// Check for Warnign as handle them
func (data *Data) checkWarnings() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err.(error))
		}
	}()
	data.warningData()
}

// Warns of strange data, assumes data if it can
func (data *Data) warningData() {
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

// Reject data that doesn't make sense
func (data *Data) assumeData() {
	// Assume any missing data that you can
}
