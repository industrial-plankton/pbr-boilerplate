package IPSheets

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(name, " took ", elapsed)
}

func WriteToSpreadSheet(SQLData [][]interface{}, rangeData string, spreadsheetId string, srv *sheets.Service) {
	defer timeTrack(time.Now(), "Write to "+rangeData)
	ctx := context.Background()
	//spreadsheetId := "1Hi0PrHe53q4JhNetcJ_y3WrDIJ9qocVEd4irMunxVyE"

	// How the input data should be interpreted.
	valueInputOption := "RAW" //"USER_ENTERED"
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

	cl := &sheets.ClearValuesRequest{
		// TODO: Add desired fields of the request body.
	}

	_, err := srv.Spreadsheets.Values.Clear(spreadsheetId, rangeData, cl).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("%#v\n", resp)
}
