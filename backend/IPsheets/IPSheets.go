package IPSheets

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

func WriteToSpreadSheet(SQLData [][]interface{}, rangeData string, spreadsheetId string, srv *sheets.Service) {
	// defer timeTrack(time.Now(), "Write")
	ctx := context.Background()
	//spreadsheetId := "1Hi0PrHe53q4JhNetcJ_y3WrDIJ9qocVEd4irMunxVyE"

	// How the input data should be interpreted.
	valueInputOption := "USER_ENTERED"
	// rangeData := "'Master Part List'!A2:Y1400"
	data := SQLData

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: valueInputOption,

		// TODO: Add desired fields of the request body.
	}
	// The new values to apply to the spreadsheet.
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: data,
	})

	resp, err := srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("%#v\n", resp)
}
