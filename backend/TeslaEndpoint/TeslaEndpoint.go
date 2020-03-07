package TeslaEndpoint

import (
	// "errors"
	"backend/IPDatabase"
	"backend/IPSheets"
	"backend/utility"
	"fmt"
	"net/http"
	"time"
)

type sampleRequest struct {
	ID int `json:"ID"`
}

func FindShipForEdit(header http.Header) (interface{}, error) {
	defer utility.TimeTrack(time.Now(), "Load Ship")
	var data [][][]interface{}
	var ranges []string

	//Start Clearing the Sheet in Parrallel
	ch := make(chan bool)
	go IPSheets.ClearRange([]string{"'Tesla Generator'!D1:T4", "'Tesla Generator'!B9:T50"}, IDMap[header.Get("UserData")], IPSheets.GetSheetsService(header), ch)

	PONum := header.Get("RequestData")
	shipIndex := fmt.Sprint(IPDatabase.Convert(db, "shipments", PONum, "our_num", "index_shipments")[0]) //get the index

	var ShipData [][]interface{}
	ShipData = append(ShipData, IPDatabase.GetHeaders(db, "shiptext"))
	DatabaseShipData, err := IPDatabase.Search(db, "public.shiptext", shipIndex, "index_shipments")
	if err != nil {
		return nil, err
	}
	for _, e := range DatabaseShipData { //append the database data
		ShipData = append(ShipData, e)
	}

	var TrackData [][]interface{}
	var DatabaseTrackData [][]interface{}
	if ShipData[1][2] == "Incomming" { //check incomming or outgoing
		TrackData = append(TrackData, IPDatabase.GetHeaders(db, "trackingtextincomming"))
		DatabaseTrackData, err = IPDatabase.Search(db, "trackingtextincomming", PONum, "\"Our Number\"")
	} else {
		TrackData = append(TrackData, IPDatabase.GetHeaders(db, "trackingtextcustomer"))
		DatabaseTrackData, err = IPDatabase.Search(db, "trackingtextcustomer", PONum, "\"Our Number\"")
	}
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
	var ShipData2 [][]interface{}
	//split Shipdata
	for i, e := range ShipData {
		ShipData[i] = e[:10]
		ShipData2 = append(ShipData2, e[10:])
	}

	//combine data and their ranges for spreadsheet write
	data = append(data, ShipData)
	ranges = append(ranges, "'Tesla Generator'!D1:Z2")
	data = append(data, ShipData2)
	ranges = append(ranges, "'Tesla Generator'!E3:Z4")
	//Split TrackData so it doesn't fill in the AutoFill Row
	data = append(data, TrackData[:1])
	ranges = append(ranges, "'Tesla Generator'!B8:Z9") //include autofill range so it gets cleared
	data = append(data, TrackData[1:])
	ranges = append(ranges, "'Tesla Generator'!B10:Z50")
	data = append(data, noArray)
	ranges = append(ranges, "'Tesla Generator'!T10:T50")

	<-ch //Wait until Sheet Cleared
	IPSheets.BatchWriteToSheetNoClear(data, ranges, IDMap[header.Get("UserData")], IPSheets.GetSheetsService(header))

	return data[0][0][1], nil
}

func SaveTeslaEdit(header http.Header) (interface{}, error) {
	defer utility.TimeTrack(time.Now(), "Save")
	DivPoint := 3 //point at which ship data becomes tracking data
	//Sheet ranges of headers
	headerRanges := []string{"'Tesla Generator'!D1:Z1",
		"'Tesla Generator'!E3:G3",
		"'Tesla Generator'!I3:Z3",
		"'Tesla Generator'!B8:D8",
		"'Tesla Generator'!H8:N8",
		"'Tesla Generator'!P8:Q8"}
	//Sheet ranges of data
	valueRanges := []string{"'Tesla Generator'!D2:Z2",
		"'Tesla Generator'!E4:G4",
		"'Tesla Generator'!I4:Z4",
		"'Tesla Generator'!B9:D50",
		"'Tesla Generator'!H9:N50",
		"'Tesla Generator'!P9:Q50"}
	deleteFlagsRange := "'Tesla Generator'!T2:T50"
	ranges := append(headerRanges, valueRanges...)
	ranges = append(ranges, deleteFlagsRange)
	read := IPSheets.BatchGet(ranges, IDMap[header.Get("UserData")], IPSheets.GetSheetsService(header)) //Get Values from sheets

	//split the batch read back into separate variables
	headerValues := read[:len(headerRanges)]
	values := read[len(headerRanges) : len(headerRanges)+len(valueRanges)]
	deleteFlags := read[len(headerRanges)+len(valueRanges):]

	//Process autofill row
	TrackValuesGroup := values[DivPoint:]
	for mi, matrix := range TrackValuesGroup {
		if len(matrix) != 0 { //nil check
			for vali, val := range matrix[0] {
				utility.FillifEmpty(TrackValuesGroup[mi], val, vali) // fill in the data
			}
			TrackValuesGroup[mi] = TrackValuesGroup[mi][1:] // remove the autofill row
		}
		//enforce size
		utility.MatchSizes(TrackValuesGroup[mi], headerValues[mi+DivPoint][0])
	}
	//enforce size
	ShipValuesGroup := values[:DivPoint]
	for i := range ShipValuesGroup {
		utility.SetSize(ShipValuesGroup[i], len(headerValues[i][0]))
	}

	//Combine the data thats split accross ranges
	ShipHeaders := utility.ConcatSplitData(headerValues[:DivPoint])
	ShipValues := utility.ConcatSplitData(ShipValuesGroup)
	TrackHeaders := utility.ConcatSplitData(headerValues[DivPoint:])
	TrackValues := utility.ConcatSplitData(TrackValuesGroup)

	PO := ShipValues[0][4]                                                                                           //save the PO from sheets
	TrackValues = utility.OverWriteColumn(TrackValues, PO, utility.GetHeaderLocation(TrackHeaders[0], "Our Number")) //write PO into the tracking section for linking

	headerMapBase := IPDatabase.GetView(db, "column_map")     //collect header collection from db
	headerMap := utility.BuildMap(headerMapBase, []int{0, 1}) //make a map with the header data
	tablesToUpdate := []string{"shipments", "partstracking"}

	err := IPDatabase.UpdateOrAdd(db, tablesToUpdate[0], headerMap, ShipValues, ShipHeaders[0])
	if err != nil {
		return nil, err
	}

	err = IPDatabase.UpdateOrAdd(db, tablesToUpdate[1], headerMap, TrackValues, TrackHeaders[0])
	if err != nil {
		return nil, err
	}

	for i := range values[1] {
		if len(deleteFlags) > 0 {
			if len(deleteFlags[0])-1 > i {
				if deleteFlags[0][i][0] == "yes" { //check for yes's in the column
					err := IPDatabase.Delete(db, tablesToUpdate[1], fmt.Sprint(TrackValues[i][0])) // send the primaryIndex to be deleted
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	po, err := IPDatabase.Search(db, tablesToUpdate[0], fmt.Sprint(PO), "our_num") // get fresh PO from DB
	if err != nil {
		return nil, err
	}

	header.Set("RequestData", fmt.Sprint(po[0][4])) //set header value to PO for FindShipForEdit
	_, err = FindShipForEdit(header)                //write data back to sheet
	if err != nil {
		return nil, err
	}

	return ShipValues[0][1], nil
}

func SearchPOs(header http.Header) (interface{}, error) {
	defer utility.TimeTrack(time.Now(), "PO Search")
	PONum := header.Get("RequestData")
	srv := IPSheets.GetSheetsService(header)
	ranges := []string{"'Search Tree'!A1:A100"}

	// ch := make(chan bool)
	// IPSheets.ClearRange(ranges, IDMap[header.Get("UserData")], srv, ch)

	var ShipData [][]interface{}
	ShipData, err := IPDatabase.Filter(db, "public.shipments", PONum, "our_num")
	if err != nil {
		return nil, err
	}
	utility.Log(ShipData)

	//combine data and their ranges for spreadsheet write
	data := make([][][]interface{}, 1)
	if len(ShipData) > 100 {
		ShipData = ShipData[:100]
	}
	data[0] = ShipData
	utility.Log(data)
	utility.Log(ranges)
	// <-ch
	IPSheets.BatchWriteToSheet(data, ranges, IDMap[header.Get("UserData")], srv)

	return nil, nil
}
