package MPLEndpoint

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

func RefreshMasterPartsList(header http.Header) ([][]interface{}, error) {
	parts := IPDatabase.GetView(db, "masterpartslist")
	IPSheets.WriteToSpreadSheet(parts, "'Master Part List'!A1:Z20000", IDMap[header.Get("UserData")], IPSheets.GetSheetsService(header))
	utility.Log(header.Get("UserData"))
	return parts, nil
}

func FindPartForEdit(header http.Header) (interface{}, error) {
	defer utility.TimeTrack(time.Now(), "Find Part")
	var data [][][]interface{}
	var ranges []string
	srv := IPSheets.GetSheetsService(header)

	ch := make(chan bool)
	go IPSheets.ClearRange([]string{"MPLEdit!B9:Z10", "MPLEdit!A11:I16"}, IDMap[header.Get("UserData")], srv, ch)

	SKU := header.Get("RequestData")
	partnumindex := fmt.Sprint(IPDatabase.Convert(db, "parts", SKU, "sku", "index_parts")[0]) //get the parts index

	var partData [][]interface{}
	partData = append(partData, IPDatabase.GetHeaders(db, "mpltext"))
	DatabasePartData, err := IPDatabase.Search(db, "public.mpltext", partnumindex, "index_parts")
	if err != nil {
		return nil, err
	}
	for _, e := range DatabasePartData { //get the parts data
		partData = append(partData, e)
	}

	var venderData [][]interface{}
	venderData = append(venderData, IPDatabase.GetHeaders(db, "partsvendertext"))
	DatabaseVendorData, err := IPDatabase.Search(db, "public.partsvendertext", SKU, "\"IP SKU\"")
	if err != nil {
		return nil, err
	}
	for _, e := range DatabaseVendorData { //get the parts vendor data
		venderData = append(venderData, e)
	}
	var noArray [][]interface{}
	for i := 0; i < 5; i++ {
		noArray = append(noArray, []interface{}{"no"})
	}
	//combine data and their ranges for spreadsheet write
	data = append(data, partData)
	ranges = append(ranges, "MPLEdit!B9:Z10")
	data = append(data, venderData)
	ranges = append(ranges, "MPLEdit!A11:H16")
	data = append(data, noArray)
	ranges = append(ranges, "MPLEdit!I12:I16")

	<-ch
	IPSheets.BatchWriteToSheetNoClear(data, ranges, IDMap[header.Get("UserData")], srv)

	return data[0][0][1], nil
}

func SaveMPLEdit(header http.Header) (interface{}, error) {
	defer utility.TimeTrack(time.Now(), "Save")
	headerRanges := []string{"MPLEdit!B9:Z9", "MPLEdit!A11:H11"}  //Sheet ranges of headers
	valueRanges := []string{"MPLEdit!B10:Z10", "MPLEdit!A12:H16"} //Sheet ranges of data
	deleteFlagsRange := "MPLEdit!I12:I16"
	ranges := append(headerRanges, valueRanges...)
	ranges = append(ranges, deleteFlagsRange)
	read := IPSheets.BatchGet(ranges, MPLID, IPSheets.GetSheetsService(header)) //Get Values from sheets

	headerValues := read[:len(headerRanges)] //split header and value reads back into separate variables
	values := read[len(headerRanges) : len(headerRanges)+len(valueRanges)]
	deleteFlags := read[len(headerRanges)+len(valueRanges):]

	sku := values[0][0][1] //save the SKU from sheets
	//TODO: should do a autogeneration of SKU here
	utility.OverWriteColumn(values[1], sku, utility.GetHeaderLocation(headerValues[1][0], "IP SKU")) //write SKU into the vendor locations for linking

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
	_, err = FindPartForEdit(header)                 //write data back to sheet
	if err != nil {
		return nil, err
	}

	return values[0][0][1], nil
}

func KeywordSearch(header http.Header) (interface{}, error) {
	defer utility.TimeTrack(time.Now(), "Keyword Search")
	keyRanges := []string{"KeywordSearch!D1:D2"}
	ranges := []string{"KeywordSearch!B10:M10000"}
	srv := IPSheets.GetSheetsService(header)

	ch := make(chan bool)
	go IPSheets.ClearRange(ranges, IDMap[header.Get("UserData")], srv, ch)

	keysInterface := IPSheets.BatchGetCol(keyRanges, MPLID, srv)
	// fmt.Print(keysInterface)
	keys := utility.IntfToString(keysInterface[0][0])
	for len(keys) < 2 { //append the empty keywords
		keys = append(keys, "")
	}
	keys = []string{keys[0], keys[0], keys[1], keys[1]} //duplicate entries for multicolumn search
	keys = utility.AddWildCards(keys)

	keycolumns := []string{"technical_desc", "customer_desc", "name", "part_number"}
	combiners := []string{"", " OR ", ") AND (", " OR "}
	SearchResult, err := IPDatabase.MultiLIKE(db, "keywordsearch", keys, keycolumns, combiners)
	if err != nil {
		return nil, err
	}
	partData := [][][]interface{}{SearchResult}

	<-ch
	IPSheets.BatchWriteToSheetNoClear(partData, ranges, IDMap[header.Get("UserData")], srv)
	return partData, nil
}
