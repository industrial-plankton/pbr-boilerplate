package Parsing

import (
	"backend/Enums/Doorway"
	"backend/Validation"
	"fmt"
	"strings"
	"time"
)

var Track = &track{} // Create variable for reference
func GetTrack() *[]TrackData { // Easy access func
	return Track.Get(Track).(*[]TrackData)
}

func CheckTrack() (errs []error) {
	_ = GetTrack()
	return Track.Errors
}

type track struct { // Create new sheet type
	Sheet
	EmptyCollection []TrackData // Shadow this to the correct type
}

func (s *track) Init() { // Initialize sheet specific TrackData
	s.Range = "2020 Tracking!A:U"
	s.SpreadsheetID = MATTRACK
}

func (s *track) Parse() { // Change type to correct sheet struct
	Sheetdata := s.getSheet()

	collection := new([]TrackData) //EmptyCollection
	for i, e := range Sheetdata {
		newData := new(trackStruct)
		newData.processNew(i, e, newData, collection, &s.Errors)
	}
	s.AllData = collection
}

type trackStruct struct { // Struct that inherits base and TrackData structs
	SheetParseBase // Base struct for Line method inheritance
	TrackData
}

// TrackData from the track that is important Must Export
type TrackData struct {
	Sku       string
	Qty       float64
	RecQty    float64
	Doorway   byte
	Completed bool
	OrderDate time.Time
	ExpDate   time.Time
	RecDate   time.Time
}

func (data *trackStruct) convData(line []interface{}) { // Converts interfaces{} to struct values
	const (
		doorwayCol       = 3
		skuCol           = 5
		qtyCol           = 6
		orderCol         = 12
		exCol            = 14
		recCol           = 16
		recQtyCol        = 17
		fullyReceivedCol = 19
	)
	data.Doorway = Doorway.ToDoorway(line[doorwayCol])
	data.Sku = Validation.Sku(line[skuCol])
	data.Qty = Validation.ConvNum(line[qtyCol])
	data.OrderDate = Validation.ConvDate(line[orderCol])
	data.ExpDate = Validation.ConvDate(line[exCol])
	data.RecDate = Validation.ConvDate(line[recCol])
	data.RecQty = Validation.ConvNum(line[recQtyCol])
	data.Completed = Validation.ConvBool(line[fullyReceivedCol])
}

func (new *trackStruct) appendNew(data interface{}) { // Adds new TrackData
	*data.(*[]TrackData) = append(*data.(*[]TrackData), new.TrackData)
}

// Reject TrackData that doesn't make sense
func (data *trackStruct) rejectData() {
	if data.OrderDate.IsZero() && data.ExpDate.IsZero() && data.RecDate.IsZero() {
		panic(fmt.Errorf("%s", "&minor& No Dates filled"))
	}

}

// Warns of strange TrackData, assumes TrackData if it can
func (data *trackStruct) warningData() {
	var sb strings.Builder

	if data.Completed && (time.Now().Before(data.RecDate) || data.RecDate.IsZero()) {
		sb.WriteString(fmt.Sprintln("Received true but Received Date invalid:", data.RecDate, "Assuming using expected or ordered date"))
	}

	if !data.Completed && !data.ExpDate.IsZero() && data.ExpDate.Before(time.Now().AddDate(0, 0, -2)) {
		sb.WriteString(fmt.Sprintln("Warning: Unfulfilled order with out of date Expected Date in the past:", data.ExpDate))
		data.ExpDate = time.Now().AddDate(0, 0, 2)
	}

	Estring := strings.TrimSpace(sb.String())
	if !(Estring == "") {
		panic(fmt.Errorf("%s", Estring))
	}
}

// Reject TrackData that doesn't make sense
func (data *trackStruct) assumeData(errors *[]error) {
	// Assume weirdly missing dates
	if data.OrderDate.IsZero() {
		if !data.ExpDate.IsZero() {
			data.OrderDate = data.ExpDate
		} else {
			data.ExpDate = data.RecDate
			data.OrderDate = data.RecDate
		}
	}

	// Assume missing dates for arrived parts
	if data.Completed && data.Doorway == Doorway.Incoming {
		if data.RecDate.IsZero() {
			if !data.ExpDate.IsZero() {
				data.RecDate = data.ExpDate
			} else {
				data.ExpDate = data.OrderDate
				data.RecDate = data.OrderDate
			}
		}
	}

	if !data.Completed && data.Doorway == Doorway.Incoming && data.ExpDate.IsZero() {
		// TODO Assume Lead Time
		lead := GetPart(data.Sku).LeadTime
		if lead == 0 {
			lead = GetVendor(GetPart(data.Sku).Supplier).Leadtime
		}
		data.ExpDate = data.OrderDate.AddDate(0, 0, int(lead))
	}
}
