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

func RefreshMasterPartsList(Auth http.Header) (parts [][]interface{}, err error) {

	parts = IPDatabase.GetMPL(db)
	IPSheets.WriteToSpreadSheet(parts, "'Master Part List'!A2:Y1400", MPLID, IPSheets.GetSheetsService(Auth))

	return parts, nil
}
