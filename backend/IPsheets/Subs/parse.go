package Subs

import (
	"backend/IPSheets"
	"backend/IPSheets/mpl"
	"backend/Validation"
	"fmt"
	"strings"
)

const (
	SpreadsheetID = "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo"
	Range         = "'SubAssem Builder'!A:I"
	// Columns
	parent = 0
	child  = 1
	qty    = 2
	loc    = 8
)

var allData map[string][]Data

func Get() map[string][]Data {
	if allData == nil {
		Refresh()
	}
	return allData
}

func Refresh() {
	allData = parse()
}

type Data struct {
	Child    string
	Qty      float64
	Location string
}

// parse maps the interface to the Data struct
func parse() map[string][]Data {
	Sheetdata := IPSheets.BatchGet([]string{Range}, SpreadsheetID, nil)[0]

	data := make(map[string][]Data)
	for i, e := range Sheetdata {
		processNew(i, e, data)
	}

	return data
}

func processNew(i int, e []interface{}, data map[string][]Data) {
	defer handleError(i)
	newData, Parent := newData(e)
	newData.rejectData()
	// newData.checkWarnings()
	// newData.assumeData()
	appendNew(newData, Parent, data)
}

// Handles errors thrown by newData, continues if only minor
func handleError(line int) {
	err := recover()
	if err != nil {
		if !strings.Contains(err.(error).Error(), "&minor&") && !strings.Contains(err.(error).Error(), "index") && line != 0 { //Print off errors that dont contain the &minor& flag, and index errors
			fmt.Println(err.(error).Error(), ", line:", line+1)
		}
	}
}

// Adds new Data
func appendNew(new Data, Parent string, data map[string][]Data) {
	if data[Parent] == nil {
		data[Parent] = []Data{}
	}
	// Empty Map Keys default to types nil val
	data[Parent] = append(data[Parent], new)
}

// Formats and Checks new Data struct
func newData(line []interface{}) (data Data, Parent string) {
	defer func() {
		err := recover().(error)
		if strings.Contains(err.Error(), "index") {
			return
		}
		if strings.Contains(err.Error(), "&continue&") { //Print off continue flagged errors
			fmt.Println(err)
		} else {
			// Rethrow
			panic(err)
		}
	}()
	Parent = Validation.ConvString(line[parent])
	data.Child = Validation.ConvString(line[child])
	data.Qty = Validation.ConvNum(line[qty])
	data.Location = Validation.ConvString(line[loc])

	return
}

func (data Data) rejectData() {
	if (data.Child == "") || (data.Qty == 0) {
		panic(fmt.Errorf("Missing Data, , Child: %s, Qty: %g", data.Child, data.Qty))
	}
}

func FindOffspring(parent string, OnlyImportant bool) map[string]float64 {
	sheetData := Get()
	mpl := mpl.Get()

	lostChildren := make(map[string]float64)

	children, hasChildren := sheetData[parent]
	if hasChildren {
		for _, line := range children {
			_, alsoHasChildren := sheetData[line.Child]
			if alsoHasChildren && (!OnlyImportant || !mpl[line.Child].CountType) {
				grandkids := FindOffspring(line.Child, OnlyImportant)
				for grandchild, amount := range grandkids {
					lostChildren[grandchild] += amount * line.Qty
				}
			} else {
				lostChildren[line.Child] += line.Qty
			}
		}
	} else {
		lostChildren[parent] = 1
	}

	return lostChildren
}
