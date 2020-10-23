package Tracking

import (
	"backend/IPSheets"
	"fmt"
	"time"
)

const (
	SpreadsheetID = "1pdhA4p4n4LbOQCrJgmSDZOzHBtV6mIfF2JUUrtxvGuc"
	Range         = "2020 Tracking!A:U"
	skuCol        = 5
	qtyCol        = 6
	// unitCol       = 2
	orderCol = 12
	exCol    = 14
	shipCol  = 16
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
	Sku string
	Qty float64
	// Unit      string
	OrderDate time.Time
	ExpDate   time.Time
	ShipDate  time.Time
}

// parse maps the interface to the mpl Data struct
func parse() []Data {
	Sheetdata := IPSheets.BatchGet([]string{Range}, SpreadsheetID, nil)[0]

	data := []Data{}
	for i, e := range Sheetdata {
		for len(e) <= 19 {
			e = append(e, "")
		}
		SKU := IPSheets.ConvString(e[skuCol])
		if SKU == "" {
			continue
		}
		// Unit := IPSheets.ConvString(e[unitCol])
		Qty, err := IPSheets.ConvNum(e[qtyCol])
		Order, err := IPSheets.ConvDate(e[orderCol])
		Ex, err := IPSheets.ConvDate(e[exCol])
		Ship, err := IPSheets.ConvDate(e[shipCol])
		if err != nil {

		}
		newData := Data{SKU, Qty, Order, Ex, Ship}
		err = validate(newData)
		if err != nil {
			fmt.Println(err, ", line:", i+1)
			continue
		}
		data = append(data, newData)
	}

	return data
}

func validate(data Data) error {
	if len(data.Sku) != 5 {
		return fmt.Errorf("Invalid SKU, %s", data.Sku)
	}

	return nil
}
