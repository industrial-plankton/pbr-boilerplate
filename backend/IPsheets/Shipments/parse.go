package Shipments

// import (
// 	"backend/IPSheets"
// 	"backend/Validation"
// 	"fmt"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// const (
// 	spreadsheetID = "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo"
// 	Range         = "'Shipments'!A:N"
// 	timeLayout    = "2006/01/02"
// 	// Zero indexed Columns
// 	skusCol = 0
// 	staged  = 11
// 	shipped = 12
// 	alloted = 10
// )

// var allData []Data

// func Get() []Data {
// 	if allData == nil {
// 		Refresh()
// 	}
// 	return allData
// }

// func Refresh() {
// 	allData = parse()
// }

// // Data from the MPL that is important
// type Data struct {
// 	Sku     string
// 	Qty     float64
// 	Staged  time.Time
// 	Shipped time.Time
// 	Alloted time.Time
// }

// /// parse maps the interface to the mpl Data struct
// func parse() []Data {
// 	Sheetdata := IPSheets.BatchGet([]string{Range}, spreadsheetID, nil)[0]

// 	data := []Data{}
// 	for i, e := range Sheetdata {
// 		if len(e) <= 12 {
// 			fmt.Println("Missing Data", e, "line:", i+1)
// 			continue
// 		}
// 		CombSKUs := Validation.ConvStringUpper(e[skusCol])
// 		// allotedDate := time.Parse(strings.Trim(e[skusCol].(string), " \r\n"))
// 		allotedDate, err := time.Parse(timeLayout, strings.TrimSpace(e[alloted].(string)))
// 		if err != nil {
// 			fmt.Println(err, ", alloted Date line:", i+1)
// 			continue
// 		}
// 		stagedDate, err := time.Parse(timeLayout, strings.TrimSpace(e[staged].(string)))
// 		if err != nil && strings.TrimSpace(e[staged].(string)) != "" {
// 			fmt.Println(err, ", staged Date line:", i+1)
// 			continue
// 		}
// 		shippedDate, err := time.Parse(timeLayout, strings.TrimSpace(e[shipped].(string)))
// 		if err != nil && strings.TrimSpace(e[shipped].(string)) != "" {
// 			fmt.Println(err, ", shipped Date line:", i+1)
// 		}

// 		SKUs := strings.Split(CombSKUs, ",")
// 		for _, sku := range SKUs {
// 			qtys := strings.Split(strings.TrimSpace(sku), "*")
// 			var qty float64 = 1
// 			if len(qtys) > 1 {
// 				qty, _ = strconv.ParseFloat(strings.TrimSpace(qtys[1]), 32)
// 			}
// 			newData := Data{strings.TrimSpace(qtys[0]), qty, stagedDate, shippedDate, allotedDate}
// 			err := validate(newData)
// 			if err != nil {
// 				fmt.Println(err, ", line:", i+1)
// 				continue
// 			}
// 			data = append(data, newData)
// 		}

// 	}

// 	return data
// }

// func validate(data Data) error {
// 	if len(data.Sku) != 5 {
// 		return fmt.Errorf("Invalid SKU, %s", data.Sku)
// 	}
// 	if data.Alloted.IsZero() && data.Shipped.IsZero() && data.Staged.IsZero() {
// 		return fmt.Errorf("Missing/Unparsable Dates")
// 	}
// 	if data.Qty == 0 {
// 		return fmt.Errorf("Quantity Error for sku: %s , Qty: %g", data.Sku, data.Qty)
// 	}
// 	return nil
// }
