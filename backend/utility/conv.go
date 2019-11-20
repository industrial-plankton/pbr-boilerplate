package utility

import (
	"fmt"

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

func RearrangeHeaders(headerMap [][]interface{}, sheetsHeaders []interface{}) []interface{} {
	mapper := BuildMap(headerMap)
	var headers []interface{}
	for _, e := range sheetsHeaders {
		e, _ := mapper.GetInverse(e)
		headers = append(headers, e)
	}

	return headers
}

func BuildMap(data [][]interface{}) *bimap.BiMap {
	biMap := bimap.NewBiMap()
	indexCol := len(data[0]) - 1
	for i, _ := range data {
		biMap.Insert(data[i][0], data[i][indexCol])
	}
	return biMap
}
