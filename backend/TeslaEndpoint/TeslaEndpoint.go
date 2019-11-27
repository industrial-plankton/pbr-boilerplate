package TeslaEndpoint

import (
	// "errors"
	"backend/IPDatabase"
	"backend/IPsheets"
	"backend/utility"
	"fmt"
	"net/http"
	// "time"
)

type sampleRequest struct {
	ID int `json:"ID"`
}

func FindShipForEdit(header http.Header) (interface{}, error) {
	var data [][][]interface{}
	var ranges []string

	PONum := header.Get("RequestData")
	shipIndex := fmt.Sprint(IPDatabase.Convert(db, "shipments", PONum, "out_num", "index_shipments")[0]) //get the index

	var ShipData [][]interface{}
	ShipData = append(ShipData, IPDatabase.GetHeaders(db, "shiptext"))
	DatabaseShipData, err := IPDatabase.Search(db, "public.shiptext", shipIndex, "index_shipments")
	if err != nil {
		return nil, err
	}
	for _, e := range DatabaseShipData { //insert the database data
		ShipData = append(ShipData, e)
	}

	var TrackData [][]interface{}
	TrackData = append(TrackData, IPDatabase.GetHeaders(db, "partstrackingtext"))
	DatabaseTrackData, err := IPDatabase.Search(db, "public.partsvendertext", shipIndex, "\"Our PO\"")
	if err != nil {
		return nil, err
	}
	for _, e := range DatabaseTrackData { //get the parts vendor data
		TrackData = append(TrackData, e)
	}
	var noArray [][]interface{}
	for i := 0; i < 5; i++ {
		noArray = append(noArray, []interface{}{"no"})
	}
	//combine data and their ranges for spreadsheet write
	data = append(data, ShipData)
	ranges = append(ranges, "MPLEdit!B9:Z10")
	data = append(data, TrackData)
	ranges = append(ranges, "MPLEdit!A11:H16")
	data = append(data, noArray)
	ranges = append(ranges, "MPLEdit!I12:I16")

	IPSheets.BatchWriteToSheet(data, ranges, IDMap[header.Get("UserData")], IPSheets.GetSheetsService(header))

	return data[0][0][1], nil
}

func SaveTeslaEdit(header http.Header) (interface{}, error) {
	// defer utility.TimeTrack(time.Now(), "Save")
	headerRanges := []string{"MPLEdit!B9:Z9", "MPLEdit!A11:H11"}  //Sheet ranges of headers
	valueRanges := []string{"MPLEdit!B10:Z10", "MPLEdit!A12:H16"} //Sheet ranges of data
	deleteFlagsRange := "MPLEdit!I12:I16"
	ranges := append(headerRanges, valueRanges...)
	ranges = append(ranges, deleteFlagsRange)
	read := IPSheets.BatchGet(ranges, IDMap[header.Get("UserData")], IPSheets.GetSheetsService(header)) //Get Values from sheets

	headerValues := read[:len(headerRanges)] //split header and value reads back into separate variables
	values := read[len(headerRanges) : len(headerRanges)+len(valueRanges)]
	deleteFlags := read[len(headerRanges)+len(valueRanges):]

	sku := values[0][0][1] //save the SKU from sheets
	//TODO: should do a autogeneration of SKU here
	utility.OverWriteColumn(values[1], sku, utility.GetHeaderLocation(headerValues[1][0], "IP SKU")) //write SKU into the vendor locations

	headerMapBase := IPDatabase.GetView(db, "column_map")     //collect header collection from db
	headerMap := utility.BuildMap(headerMapBase, []int{0, 1}) //make a map with the header data
	tablesToUpdate := []string{"parts", "parts_vendor"}

	for i := range values {
		err := IPDatabase.UpdateOrAdd(db, tablesToUpdate[i], headerMap, values[i], headerValues[i][0])
		if err != nil {
			return nil, err
		}
	}

	for i := range values[1] {
		if len(deleteFlags) > 0 {
			if len(deleteFlags[0])-1 > i {
				if deleteFlags[0][i][0] == "yes" { //check for yes's in the column
					err := IPDatabase.Delete(db, tablesToUpdate[1], fmt.Sprint(values[1][i][0])) // send the primaryIndex to be deleted
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	SKU, err := IPDatabase.Search(db, tablesToUpdate[0], fmt.Sprint(sku), "sku") // get SKU
	if err != nil {
		return nil, err
	}
	header.Set("RequestData", fmt.Sprint(SKU[0][0])) //set header value to SKU for FindPartForEdit
	_, err = FindShipForEdit(header)                 //write data back to sheet
	if err != nil {
		return nil, err
	}

	return values[0][0][1], nil
}
