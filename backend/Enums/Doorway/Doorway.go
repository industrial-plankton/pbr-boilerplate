package Doorway

import (
	"strings"
)

const (
	Null      byte = 0
	Incoming  byte = 1
	Outgoing  byte = 2
	Inventory byte = 3
	Shipment  byte = 4
)

func ToDoorway(door interface{}) (doorway byte) {
	// door = strings.ToUpper(strings.TrimSpace(door.(string)))
	// switch door {
	// case "Purchase Mfct":
	// 	doorway = Incoming
	// case "OUTGOING":
	// 	doorway = Outgoing
	// case "INVENTORY":
	// 	doorway = Inventory
	// case "SHIPMENT":
	// 	doorway = Shipment
	// default:
	// 	doorway = Null
	// 	// panic(fmt.Errorf("Not a Doorway: %s", door))
	// }
	door = strings.TrimSpace(door.(string))
	switch door {
	case "Purchase Mfct",
		"SS",
		"Purchase R&D A-Cons",
		"Purchase R&D A-Tform",
		"Purchase R&D R-Cons",
		"Purchase R&D R-Tform",
		"Purchase R&D Non-SRED",
		"Sale Drop Shipment",
		"Warranty Drop Shipment",
		"PC Parts Drop Shipment",
		"Work Order",
		"Shop Equipment/Tools",
		"Office Supplies":
		doorway = Incoming
	case "Sale",
		"Warranty",
		"ASA/CWS/PC",
		"Marketing Parts",
		"Shrink",
		"Obsolete Part Exit",
		"R&D Consumed for Algae",
		"R&D Transformed Algae",
		"R&D Consumed Rotifer",
		"R&D Transformed Rotifer",
		"R&D Non-SRED":
		doorway = Outgoing
	default:
		doorway = Null
		// panic(fmt.Errorf("Not a Doorway: %s", door))
	}
	return
}
