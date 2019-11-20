package MPLEndpoint

import (
	// "errors"
	"backend/IPDatabase"
	"backend/IPsheets"
	"backend/utility"
	"fmt"
	"net/http"
)

type sampleRequest struct {
	ID int `json:"ID"`
}

func RefreshMasterPartsList(header http.Header) ([][]interface{}, error) {
	// headers := []interface{}{"IP SKU", "Technical Description", "Customer Description", "Supplier (Main)", "Main Supplier PN", "Supplier (Secondary)", "Secondary supplier PN", "Extra Info", "Order Type", "unit", "Cost/unit", "currency", "Cost/ea or ft (CAD)", "Shipping Cost/ea (CAD)", "Sell Markup", "Sell price per ea or ft (USD, calculated)", "HS Code", "Re-order q-ty", "Part Location", "Book Type", "Sell price per ea or ft (USD, Static)"}
	parts := IPDatabase.GetView(db, "new_view" /*, headers*/)
	IPSheets.WriteToSpreadSheet(parts, "'Master Part List'!A1:Y1400", MPLID, IPSheets.GetSheetsService(header))
	return parts, nil
}

func FindPartForEdit(header http.Header) (interface{}, error) {
	var data [][][]interface{}
	var ranges []string

	SKU := header.Get("RequestData")
	partnumindex := fmt.Sprintf("%v", IPDatabase.Convert(db, "parts", SKU, "sku", "index_parts")[0]) //get the parts index

	var partData [][]interface{}
	partData = append(partData, IPDatabase.GetHeaders(db, "mpltext"))
	for _, e := range IPDatabase.Search(db, "public.mpltext", partnumindex, "part_index") { //get the parts data
		partData = append(partData, e)
	}

	var venderData [][]interface{}
	venderData = append(venderData, IPDatabase.GetHeaders(db, "partsvendertext"))
	for _, e := range IPDatabase.Search(db, "public.partsvendertext", partnumindex, "part_index") { //get the parts vendor data
		venderData = append(venderData, e)
	}

	//combine data and their ranges for spreadsheet write
	data = append(data, partData)
	ranges = append(ranges, "MPLEdit!B9:Z10")
	data = append(data, venderData)
	ranges = append(ranges, "MPLEdit!A11:Z16")

	IPSheets.BatchWriteToSheet(data, ranges, MPLID, IPSheets.GetSheetsService(header))

	return data[0][0][1], nil
}

func SaveMPLEdit(header http.Header) (interface{}, error) {
	// ranges := []string{"MPLEdit!B10:Z10", "MPLEdit!A12:Z16"}
	ranges := []string{"MPLEdit!B9:Z9", "MPLEdit!B10:Z10", "MPLEdit!A12:Z16"}
	values := IPSheets.BatchGet(ranges, MPLID, IPSheets.GetSheetsService(header))
	// sheetHeaders := IPSheets.BatchGet([]string{"MPLEdit!B10:Z10"}, MPLID, IPSheets.GetSheetsService(header))
	// fmt.Println(values)
	utility.Log(values)
	// columns := IPDatabase.GetHeaders(db, "parts")
	headerMap := IPDatabase.GetView(db, "column_map")
	partIndex := fmt.Sprint(values[0][0][0])

	columns := utility.RearrangeHeaders(headerMap, values[0][0])

	if IPDatabase.Exists(db, "parts", partIndex, "index_parts") {
		for _, section := range values[1:] {
			for _, row := range section {
				IPDatabase.Update(db, "parts", partIndex, columns, row)
			}
		}
	} else {
		IPDatabase.Insert(db, "parts", columns, values[0])
	}

	return values[0][0][1], nil
}

func KeywordSearch(header http.Header) (interface{}, error) {
	keyRanges := []string{"KeywordSearch!D1:D2"}
	srv := IPSheets.GetSheetsService(header)
	keysInterface := IPSheets.BatchGetCol(keyRanges, MPLID, srv)
	fmt.Print(keysInterface)
	keys := utility.IntfToString(keysInterface[0][0])
	for len(keys) < 2 { //append the empty keywords
		keys = append(keys, "")
	}
	keys = []string{keys[0], keys[0], keys[1], keys[1]} //duplicate entries for multicolumn search
	keys = IPDatabase.AddWildCards(keys)
	// keycolumns := []string{"(t.technical_desc::text", "OR t.customer_desc::text)", " And (t.name::text", "OR t.part_number::text"}
	keycolumns := []string{"technical_desc", "customer_desc", "name", "part_number"}
	combiners := []string{"", "OR", ")AND(", "OR"}
	partData := [][][]interface{}{IPDatabase.MultiLIKE(db, "public.keywordsearch", keys, keycolumns, combiners)}
	fmt.Print(partData)
	ranges := []string{"KeywordSearch!B10:M10000"}
	IPSheets.BatchWriteToSheet(partData, ranges, MPLID, srv)
	return partData, nil
}
