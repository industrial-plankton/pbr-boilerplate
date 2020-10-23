package Subs

import (
	"backend/IPSheets"
	"backend/IPSheets/mpl"
	"fmt"
	"strconv"
	"strings"
)

const (
	SpreadsheetID = "1-dbCTFEWV4oIv7fDaNldJ-mecOVzxupWrKUcugFcdNo"
	Range         = "'SubAssem Builder'!A:C"
	parent        = 0
	child         = 1
	qty           = 2
)

var allData map[string]map[string]float32

func Get() map[string]map[string]float32 {
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
func parse() map[string]map[string]float32 {
	Sheetdata := IPSheets.BatchGet([]string{Range}, SpreadsheetID, nil)[0]

	Subs := make(map[string]map[string]float32)
	var Parent, Child string
	for i, e := range Sheetdata {
		err := validate(e)
		if err != nil {
			fmt.Println(err, ", line:", i+1)
			continue
		}
		Parent = strings.ToUpper(strings.TrimSpace(e[parent].(string)))
		Child = strings.ToUpper(strings.TrimSpace(e[child].(string)))
		q, numErr := strconv.ParseFloat(e[qty].(string), 32)
		Qty := float32(q)
		if numErr == nil {
			if Subs[Parent] == nil {
				Subs[Parent] = make(map[string]float32)
			}
			// Empty Map Keys default to types nil val
			Subs[Parent][Child] += Qty
		}
	}

	return Subs
}

func validate(data []interface{}) error {
	if len(data) < 3 {
		return fmt.Errorf("Missing Data, %s", data)
	}
	if (data[parent].(string) == "") || (data[child].(string) == "") || (data[qty].(string) == "") {
		return fmt.Errorf("Missing Data, Parent: %s , Child: %s, Qty: %s", data[parent], data[child], data[qty])
	}
	return nil
}

func FindOffspring(parent string, OnlyImportant bool) map[string]float32 {
	sheetData := Get()
	mpl := mpl.Get()

	lostChildren := make(map[string]float32)

	children, hasChildren := sheetData[parent]
	if hasChildren {
		for child, qty := range children {
			_, alsoHasChildren := sheetData[child]
			if alsoHasChildren && (!OnlyImportant || !mpl[child].CountType) {
				grandkids := FindOffspring(child, OnlyImportant)
				for grandchild, amount := range grandkids {
					lostChildren[grandchild] += amount * qty
				}
			} else {
				lostChildren[child] += qty
			}
		}
	} else {
		lostChildren[parent] = 1
	}

	return lostChildren
}
