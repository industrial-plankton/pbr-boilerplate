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
	headers := []interface{}{"IP SKU", "Technical Description", "Customer Description", "Supplier (Main)", "Main Supplier PN", "Supplier (Secondary)", "Secondary supplier PN", "Extra Info", "Order Type", "unit", "Cost/unit", "currency", "Cost/ea or ft (CAD)", "Shipping Cost/ea (CAD)", "Sell Markup", "Sell price per ea or ft (USD, calculated)", "HS Code", "Re-order q-ty", "Part Location", "Book Type", "Sell price per ea or ft (USD, Static)"}
	parts := IPDatabase.GetView(db, "public.new_view", headers)
	IPSheets.WriteToSpreadSheet(parts, "'Master Part List'!A1:Y1400", MPLID, IPSheets.GetSheetsService(header))

	return parts, nil
}

func FindPartForEdit(header http.Header) (interface{}, error) {
	var data [][][]interface{}
	var ranges []string
	partnum := header.Get("RequestData")
	// partHeaders := IPDatabase.GetHeaders(db, "parts")
	partData := IPDatabase.Search(db, "public.mpltext", partnum, "sku")
	partnumindex := fmt.Sprintf("%v", partData[0][0])
	// venderHeaders := IPDatabase.GetHeaders(db, "partsvendertext")
	venderData := IPDatabase.Search(db, "public.partsvendertext", partnumindex, "part_index")

	// partHeaders = append(partHeaders, venderHeaders)

	data = append(data, partData)
	ranges = append(ranges, "MPLEdit!B10:Z10")
	data = append(data, venderData)
	ranges = append(ranges, "MPLEdit!A12:Z16")

	IPSheets.BatchWriteToSheet(data, ranges, MPLID, IPSheets.GetSheetsService(header))

	return data[0][0][1], nil
}

func SaveMPLEdit(header http.Header) (interface{}, error) {
	ranges := []string{"MPLEdit!B10:Z10", "MPLEdit!A12:Z16"}
	values := IPSheets.BatchGet(ranges, MPLID, IPSheets.GetSheetsService(header))
	fmt.Print(values)
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
