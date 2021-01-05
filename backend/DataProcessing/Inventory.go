package DataProcessing

import (
	"backend/Enums/Doorway"
	"backend/IPObjects"
	"backend/IPSheets"
	"backend/IPSheets/Parsing"
	"fmt"
	"os"
	"time"
)

type Inventory struct {
	Sku               IPObjects.Sku
	PhysicallyCounted float64
	Arrived           float64
	AssumedArrived    float64
	Sent              float64
	Alloted           float64
	OnOrder           float64
}

func (i Inventory) Theoretical() float64 {
	return i.PhysicallyCounted + i.Arrived + i.AssumedArrived - i.Sent
}

func (i Inventory) Excess() float64 {
	return i.PhysicallyCounted + i.Arrived + i.AssumedArrived - i.Sent - i.Alloted
}

func CurrentInv() {
	var Inv Inventory
	output := [][]interface{}{}
	time := time.Now()
	for K := range Parsing.GetInv() {
		Inv = NewInventory(K, time)
		line := []interface{}{Inv.Sku.Name, Inv.PhysicallyCounted, Inv.Arrived, Inv.AssumedArrived, Inv.Sent, Inv.Alloted, Inv.OnOrder, Inv.Theoretical(), Inv.Excess()}
		output = append(output, line)
	}
	IPSheets.WriteToSpreadSheet(output, "'Golang CurrInv'!A2:I", Parsing.MATTRACK, nil)
}

func NewInventory(sku string, date time.Time) (inv Inventory) {
	var SKU IPObjects.Sku
	SKU.Name = sku
	inv.Sku = SKU
	inv.PhysicallyCounted = Parsing.GetPartInv(sku).Qty
	inv.ParseTracking(Parsing.GetPartInv(sku).Date, date)
	inv.ParseShipments(Parsing.GetPartInv(sku).Date, date)
	return
}

func (inv *Inventory) ParseTracking(invDate time.Time, forcastDate time.Time) {
	for _, e := range *Parsing.GetTrack() {
		for bomParts, Qty := range Parsing.CreateBOM(e.Sku, true, 0) {
			if Qty <= 0 {
				continue
			}
			if bomParts == inv.Sku.Name {
				if !e.RecDate.IsZero() && e.RecDate.Before(invDate) && e.Completed {
					continue // Skip completed lines before the inventory date
				}
				if e.Doorway == Doorway.Incoming {

					if invDate.Unix() <= e.RecDate.Unix() {
						inv.Arrived += (e.RecQty * inv.Sku.ConvertUnit() * Qty)
					}

					if !e.Completed {
						// if Expected date is in the past adjust it to a couple days from now (so we dont assume orders have arrived that literally haven't)
						if e.ExpDate.Before(time.Now()) {
							e.ExpDate = time.Now().AddDate(0, 0, 2)
						}
						// if expected date is after the date we are checking its still on order
						if forcastDate.Before(e.ExpDate) {
							inv.OnOrder += (e.Qty * inv.Sku.ConvertUnit() * Qty)
							if !e.RecDate.IsZero() {
								inv.OnOrder -= (e.RecQty * inv.Sku.ConvertUnit() * Qty)
							}
						} else if forcastDate.Unix() > e.ExpDate.Unix() { // if the date we are checking has past
							inv.AssumedArrived += (e.Qty * inv.Sku.ConvertUnit() * Qty)
							if !e.RecDate.IsZero() {
								inv.AssumedArrived -= (e.RecQty * inv.Sku.ConvertUnit() * Qty)
							}
						}
					}

				} else if e.Doorway == Doorway.Outgoing {
					if e.Completed {
						if invDate.Before(e.RecDate) {
							inv.Sent += (e.RecQty * Qty)
						}
					} else {
						if forcastDate.Unix() > e.OrderDate.Unix() {
							inv.Alloted += (e.Qty * Qty)
						}
					}
				}
			}
		}
	}
}

func (inv *Inventory) ParseShipments(invDate time.Time, forcastDate time.Time) {
	for _, e := range *Parsing.GetShip() {
		bom := make(map[string]float64)
		for _, p := range e.Parts {
			Parsing.AmendBOM(p.Sku, p.Qty, true, 0, bom)
		}
		for part, qty := range bom {
			if qty <= 0 {
				continue
			}
			if part == inv.Sku.Name {
				if os.Getenv("DEBUG") == "1" {
					fmt.Println(fmt.Sprint(qty) + " " + fmt.Sprint(e.Alloted) + " " + fmt.Sprint(e.Shipped))
				}

				if invDate.Before(e.Shipped) && forcastDate.Unix() > e.Shipped.Unix() {
					inv.Sent += qty
					if os.Getenv("DEBUG") == "1" {
						fmt.Println("Shipped")
					}
				} else if invDate.Before(e.Shipped) && forcastDate.Unix() > e.Alloted.Unix() {
					inv.Alloted += qty
					if os.Getenv("DEBUG") == "1" {
						fmt.Println("Alloted")
					}

				}

			}
		}
	}
}
