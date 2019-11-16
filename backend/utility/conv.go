package utility

import (
	"fmt"
)

func IntfToString(data []interface{}) []string {
	//convert an interface slice to string slice
	out := make([]string, len(data))
	for i, e := range data {
		out[i] = fmt.Sprintf("%v", e)
	}
	return out
}
