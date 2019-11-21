package utility

import (
	"fmt"
	"time"

	"github.com/vishalkuo/bimap"
)

func IntfToString(data []interface{}) []string {
	//convert an interface slice to string slice
	out := make([]string, len(data))
	for i, e := range data {
		out[i] = fmt.Sprintf("%v", e)
	}
	return out
}

func RearrangeHeaders(headerMap *bimap.BiMap, sheetsHeaders []interface{}) []interface{} {
	defer TimeTrack(time.Now(), "rearrange: ")
	// mapper := BuildMap(headerMap, []int{0, 1})
	var headers []interface{}
	for _, e := range sheetsHeaders {
		e, _ := headerMap.GetInverse(e)
		headers = append(headers, e)
	}

	return headers
}

func BuildMap(data [][]interface{}, colIndex []int) *bimap.BiMap {
	biMap := bimap.NewBiMap()
	if len(colIndex) == 1 {
		colIndex = append(colIndex, len(data[0])-1)
	}
	for i, _ := range data {
		biMap.Insert(data[i][colIndex[0]], data[i][colIndex[1]])
	}
	return biMap
}

func ParseNulls(data [][]interface{}) [][]interface{} { //puts '' around data, replaces empty spots with NULL
	for i, e := range data {
		for ri, re := range e {
			if re == "" {
				data[i][ri] = "NULL"
			} else {
				data[i][ri] = "'" + fmt.Sprint(re) + "'"
			}
		}
	}
	Log(data)
	return data
}
