package Parsing

import (
	"backend/IPSheets"
	"backend/Validation"
	"fmt"
)

var DataValidation = IPSheets.Sheet{}

// type DataVal struct {
// 	IPSheets.Sheet
// 	emptyCollection []Data // Shadow this to the correct type when using
// 	emptyData       Data   // Shadow this to the correct type when using
// }

type Data struct {
	IPSheets.Data
	One   string
	Two   string
	Three string
}

func GetDataVal() []Data {
	if DataValidation.Range == "" {
		DataValidation.Range = "'Data Validation'!B:D"
		DataValidation.SpreadsheetID = "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo"
		DataValidation.EmptyCollection = []Data{}
		DataValidation.EmptyData = &Data{}
	}
	return DataValidation.Get().([]Data) //.(emptyCollectiontype)
}

var allData []Data

func Get() []Data {
	if allData == nil {
		Refresh()
	}
	return allData
}

func Refresh() {
	allData = parse()
}

// parse maps the interface to the mpl Data struct
func parse() []Data {
	Sheetdata := IPSheets.BatchGet([]string{"'Data Validation'!B:D"}, "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo", nil)[0]
	fmt.Print("wor")
	data := []Data{}
	for i, e := range Sheetdata {
		newData := &Data{}
		var inter IPSheets.Line = newData
		inter.ProcessNew(i, e, inter)
		inter.AppendNew(&data)
	}

	return data
}

func (data *Data) ConvData(line []interface{}) {
	fmt.Println("woooos")
	data.One = Validation.ConvString(line[0])
	data.Two = Validation.ConvString(line[1])
	data.Three = Validation.ConvString(line[2])
}

// Adds new Data
func (new *Data) AppendNew(data interface{}) {

	data = append(*data.(*[]Data), *new)

}
