package Tracking

import (
	"backend/Enums/Doorway"
	"backend/IPSheets"
	"backend/Validation"
	"fmt"
	"strings"
	"time"
)

const (
	SpreadsheetID = "1pdhA4p4n4LbOQCrJgmSDZOzHBtV6mIfF2JUUrtxvGuc"
	Range         = "2020 Tracking!A:U"
	// Columns
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

func Get() []Data {
	if allData == nil {
		Refresh()
	}
	return allData
}

func Refresh() {
	allData = parse()
}

// Data from the MPL that is important
type Data struct {
	Sku       string
	Qty       float64
	RecQty    float64
	Doorway   byte
	Completed bool
	OrderDate time.Time
	ExpDate   time.Time
	RecDate   time.Time
}

// parse maps the interface to the mpl Data struct
func parse() []Data {
	Sheetdata := IPSheets.BatchGet([]string{Range}, SpreadsheetID, nil)[0]

	data := []Data{}
	for i, e := range Sheetdata {
		processNew(i, e, &data)
	}

	return data
}

func processNew(i int, e []interface{}, data *[]Data) {
	defer handleError(i)
	newData := newData(e)
	newData.rejectData()
	newData.checkWarnings()
	newData.assumeData()
	appendNew(newData, data)
}

// Handles errors thrown by newData, continues if only minor
func handleError(line int) {
	err := recover()
	if err != nil {
		if !strings.Contains(err.(error).Error(), "&minor&") && !strings.Contains(err.(error).Error(), "index") && line != 0 { //Print off errors that dont contain the &minor& flag, and index errors
			fmt.Println(err.(error).Error(), ", line:", line+1)
		}
	}
}

// Adds new Data
func appendNew(new Data, data *[]Data) {
	*data = append(*data, new)
}

// Formats and Checks new Data struct
func newData(line []interface{}) (data Data) {
	defer func() {
		err := recover().(error)
		if strings.Contains(err.Error(), "index") {
			return
		}
		if strings.Contains(err.Error(), "&continue&") { //Print off continue flagged errors
			fmt.Println(err)
		} else {
			// Rethrow
			panic(err)
		}
	}()
	data.Doorway = Doorway.ToDoorway(line[doorwayCol])
	data.Sku = Validation.Sku(line[skuCol])
	data.Qty = Validation.ConvNum(line[qtyCol])
	data.OrderDate = Validation.ConvDate(line[orderCol])
	data.ExpDate = Validation.ConvDate(line[exCol])
	data.RecDate = Validation.ConvDate(line[recCol])
	data.RecQty = Validation.ConvNum(line[recQtyCol])
	data.Completed = Validation.ConvBool(line[fullyReceivedCol])

	return
}

// Reject data that doesn't make sense
func (data Data) rejectData() {
	if data.OrderDate.IsZero() && data.ExpDate.IsZero() && data.RecDate.IsZero() {
		panic(fmt.Errorf("%s", "&minor& No Dates filled"))
	}

}

// Check for Warnign as handle them
func (data *Data) checkWarnings() {
	defer func() {
		err := recover().(error)
		fmt.Println(err)
	}()
	data.warningData()
}

// Warns of strange data, assumes data if it can
func (data *Data) warningData() {
	var sb strings.Builder

	if data.Completed && (time.Now().Before(data.RecDate) || data.RecDate.IsZero()) {
		sb.WriteString(fmt.Sprintln("Received true but Received Date invalid:", data.RecDate))
		data.RecDate = time.Now()
	}

	if !data.Completed && !data.ExpDate.IsZero() && data.ExpDate.Before(time.Now().AddDate(0, 0, -2)) {
		sb.WriteString(fmt.Sprintln("Warning: Unfulfilled order with out of date Expected Date in the past:", data.ExpDate))
		data.ExpDate = time.Now().AddDate(0, 0, 2)
	}

	Estring := strings.TrimSpace(sb.String())
	if !(Estring == "") {
		panic(fmt.Errorf("%s", Estring))
	}
}

// Reject data that doesn't make sense
func (data *Data) assumeData() {
	// Assume weirdly missing dates
	if data.OrderDate.IsZero() {
		if !data.ExpDate.IsZero() {
			data.OrderDate = data.ExpDate
		} else {
			data.ExpDate = data.RecDate
			data.OrderDate = data.RecDate
		}
	}

	// Assume missing dates for arrived parts
	if data.Completed && data.Doorway == Doorway.Incoming {
		if data.RecDate.IsZero() {
			if !data.ExpDate.IsZero() {
				data.RecDate = data.ExpDate
			} else {
				data.ExpDate = data.OrderDate
				data.RecDate = data.OrderDate
			}
		}
	}

	if !data.Completed && data.Doorway == Doorway.Incoming && data.ExpDate.IsZero() {
		// TODO Assume Lead Time
	}
}
