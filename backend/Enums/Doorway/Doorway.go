package Doorway

import (
	"strings"
)

const (
	Null      = 0
	Incoming  = 1
	Outgoing  = 2
	Inventory = 3
	Shipment  = 4
)

func ToDoorway(door interface{}) (doorway byte) {
	door = strings.ToUpper(strings.TrimSpace(door.(string)))
	switch door {
	case "INCOMING":
		doorway = Incoming
	case "OUTGOING":
		doorway = Outgoing
	case "INVENTORY":
		doorway = Inventory
	case "SHIPMENT":
		doorway = Shipment
	default:
		doorway = Null
		// panic(fmt.Errorf("Not a Doorway: %s", door))
	}
	return
}
