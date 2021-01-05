package Validation

import (
	"fmt"
	"strings"
)

func Sku(sku interface{}) (SKU string) {
	SKU = strings.ToUpper(ConvString(sku))
	if SKU == "" {
		panic(fmt.Errorf("%s", "&minor& No SKU"))
	}
	if len(SKU) != 5 {
		panic(fmt.Errorf("Invalid SKU: %s", sku.(string)))
	}
	return
}
