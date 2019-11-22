package IPSheets

import (
	"log"
	"time"

	"backend/utility"

	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

func WriteToSpreadSheet(SQLData [][]interface{}, rangeData string, spreadsheetId string, srv *sheets.Service) {
	defer utility.TimeTrack(time.Now(), "Write to "+rangeData)

	ctx := context.Background()

	// How the input data should be interpreted.
	valueInputOption := "RAW" //"USER_ENTERED"
	data := SQLData

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: valueInputOption,
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
		// log.Fatal(err)
		utility.Log(err)
	}

	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		// log.Fatal(err)
		utility.Log(err)
	}

	// TODO: Change code below to process the `resp` object:
	// fmt.Printf("%#v\n", resp)
}

func BatchWriteToSheet(SQLData [][][]interface{}, rangeData []string, spreadsheetId string, srv *sheets.Service) {
	defer utility.TimeTrack(time.Now(), "Batch Write to sheet")
	ctx := context.Background()

	// How the input data should be interpreted.
	valueInputOption := "RAW" //"USER_ENTERED"

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: valueInputOption,
	}
	cl := &sheets.ClearValuesRequest{
		// TODO: Add desired fields of the request body.
	}

	// The combine ValueRanges.
	for i, e := range rangeData {
		rb.Data = append(rb.Data, &sheets.ValueRange{
			Range:  e,
			Values: SQLData[i],
		})
		_, err := srv.Spreadsheets.Values.Clear(spreadsheetId, e, cl).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err := srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	// fmt.Printf("%#v\n", resp)
}

func BatchGet(rangeData []string, spreadsheetId string, srv *sheets.Service) [][][]interface{} {
	defer utility.TimeTrack(time.Now(), "Batch Read from sheet")
	ctx := context.Background()

	resp, err := srv.Spreadsheets.Values.BatchGet(spreadsheetId).Ranges(rangeData...).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	var data [][][]interface{}
	for _, e := range resp.ValueRanges {
		v := *e
		data = append(data, v.Values)
	}

	return data
}

func BatchGetCol(rangeData []string, spreadsheetId string, srv *sheets.Service) [][][]interface{} {
	defer utility.TimeTrack(time.Now(), "Batch Read from sheet")
	ctx := context.Background()

	resp, err := srv.Spreadsheets.Values.BatchGet(spreadsheetId).MajorDimension("COLUMNS").Ranges(rangeData...).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	var data [][][]interface{}
	for _, e := range resp.ValueRanges {
		v := *e
		data = append(data, v.Values)
	}

	return data
}
