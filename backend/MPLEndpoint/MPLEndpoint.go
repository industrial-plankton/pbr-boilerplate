package MPLEndpoint

import (
	// "errors"
	"backend/IPDatabase"
	"backend/IPsheets"
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
