package IPObjects

import (
	"backend/IPSheets/Parsing"
	"fmt"
	"time"
)

type Sku struct {
	Name       string
	mplInfo    Parsing.MplData
	leadTime   float64
	costCAD    float64
	invDate    time.Time
	countedQty float64
}

const (
	S = 1
	B = 2
	M = 3
	P = 4
	R = 5
	W = 6
)

func (s Sku) MplData() Parsing.MplData {
	if s.mplInfo.Desc == "" {
		s.mplInfo = Parsing.GetPart(s.Name)
	}
	return s.mplInfo
}

func (s Sku) Category() (cat byte) {
	name := s.Name[0]
	switch name {
	case 'S':
		cat = S
	case 'B':
		cat = B
	case 'M':
		cat = M
	case 'P':
		cat = P
	case 'R':
		cat = R
	case 'W':
		cat = W
	default:
		cat = 0
		// panic(fmt.Errorf("Not a Doorway: %s", door))
	}
	return
}

func (s Sku) CostPerUnitCAD() float64 {
	if s.costCAD == 0.0 {
		if s.Category() == S {
			cost := 0.0
			for K, V := range Parsing.CreateFlaggedBOM(s.Name, false, 0) {
				var sku Sku
				sku.Name = K
				cost += V * sku.CostPerUnitCAD()
			}
			s.costCAD = cost
		} else {
			s.costCAD = s.MplData().CostPerUnit * ConvertCurrency(s.Currency())
		}
	}
	return s.costCAD
}

func (s Sku) Currency() string {
	return Parsing.GetVendor(s.MplData().Supplier).Currency
}

func (s Sku) ConvertUnit() float64 {
	mult := ConvertUnit(s.MplData().Units)
	if mult == 0 {
		fmt.Println(s.Name + " " + s.MplData().Units + " bad units")
	}
	return mult
}

func (s Sku) LeadTime() (days float64) {
	if s.leadTime == 0 {
		if s.MplData().LeadTime != 0 {
			s.leadTime = s.MplData().LeadTime
		} else {
			s.leadTime = Parsing.GetVendor(s.MplData().Supplier).Leadtime
		}
	}
	return s.leadTime
}

func ConvertUnit(unit string) float64 {
	for _, e := range *Parsing.GetUnitConv() {
		if e.Unit == unit {
			return e.Multiplier
		}
	}
	return 0
}

func ConvertCurrency(currency string) float64 {
	for _, e := range *Parsing.GetCurrencyConv() {
		if e.Currency == currency {
			return e.Multiplier
		}
	}
	return 1
}
