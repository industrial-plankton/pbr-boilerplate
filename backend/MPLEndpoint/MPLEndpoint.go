package MPLEndpoint

import (
	// "errors"
	"net/http"

	"backend/IPDatabase"
	"backend/IPsheets"
)

type sampleRequest struct {
	ID int `json:"ID"`
}

func RefreshMasterPartsList(Auth http.Header) ([][]interface{}, error) {
	headers := []interface{}{"IP SKU", "Technical Description", "Customer Description", "Supplier (Main)", "Main Supplier PN", "Supplier (Secondary)", "Secondary supplier PN", "Extra Info", "Order Type", "unit", "Cost/unit", "currency", "Cost/ea or ft (CAD)", "Shipping Cost/ea (CAD)", "Sell Markup", "Sell price per ea or ft (USD, calculated)", "HS Code", "Re-order q-ty", "Part Location", "Book Type", "Sell price per ea or ft (USD, Static)"}
	parts := IPDatabase.GetView(db, "public.new_view", headers)
	IPSheets.WriteToSpreadSheet(parts, "'Master Part List'!A1:Y1400", MPLID, IPSheets.GetSheetsService(Auth))

	return parts, nil
}
