package DataValidationEndpoint

import (
	// "errors"
	"backend/IPDatabase"
	"backend/IPsheets"
	"backend/utility"

	// "fmt"
	"net/http"
	"time"
)

type sampleRequest struct {
	ID int `json:"ID"`
}

func Refresh(header http.Header) (interface{}, error) {
	defer utility.TimeTrack(time.Now(), "Load Ship")
	// var data [][][]interface{}
	ranges := []string{"'Database Data Validation'!A1:D60",
		"'Database Data Validation'!E1:F60",
		"'Database Data Validation'!G1:H60",
		"'Database Data Validation'!I1:J60",
		"'Database Data Validation'!K1:L60",
		"'Database Data Validation'!M1:N60",
		"'Database Data Validation'!O1:Q60"}
	tables := []string{"unit",
		"order_type",
		"book_type",
		"terms",
		"how_we_order",
		"how_we_pay",
		"currency"}

	data := make([][][]interface{}, len(tables))
	for i := range tables {
		// DatabaseData := IPDatabase.GetView(db, tables[i])
		// data = append(data, DatabaseData)
		data[i] = IPDatabase.GetView(db, tables[i])
	}

	IPSheets.BatchWriteToSheetNoClear(data, ranges, IDMap[header.Get("UserData")], IPSheets.GetSheetsService(header))

	return nil, nil
}
