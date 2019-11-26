package utility

import (
	"fmt"
	"strings"
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

func AddWildCards(array []string) []string { //adds regular expresstion wildcars to begining and end of each word
	for i, e := range array {
		term := strings.ReplaceAll(e, " ", ".*)(.*")
		array[i] = "'(.*" + strings.TrimSuffix(term, "\n") + ".*)'"
	}
	return array
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
			if re == "" || re == nil {
				data[i][ri] = "NULL"
			} else {
				data[i][ri] = "'" + fmt.Sprint(re) + "'"
			}
		}
	}
	return data
}

func FindPrimIndexLocation(columns []interface{}) int {
	for i, e := range columns {
		if strings.HasPrefix(fmt.Sprint(e), "index_") {
			return i
		}
	}
	return -1
}

func FindUnIndexedLocation(table string, columns []interface{}) int {
	for i, e := range columns {
		if e == table+"_index" {
			return i
		}
	}
	return -1
}

func FindTranslationTables(table string, columns []interface{}) []string {
	var translationTables []string
	for i := range columns {
		columnString := fmt.Sprint(columns[i])
		if strings.HasSuffix(columnString, "_index") {
			if !strings.HasPrefix(columnString, table) { //dont add the one referencing the table itself
				translationTables = append(translationTables, strings.TrimSuffix(columnString, "_index"))
			}
		}
	}
	return translationTables
}

func GetHeaderLocation(columns []interface{}, header string) int {
	for i, e := range columns {
		if e == header {
			return i
		}
	}
	return -1
}

func OverWriteColumn(data [][]interface{}, value interface{}, column int) [][]interface{} { //fill a colum with a value
	for i := range data {
		data[i][column] = value
	}
	return data
}

func MatchSizes(data [][]interface{}, size []interface{}) [][]interface{} {
	row := make([]interface{}, len(size))
	for i := range data {
		data[i] = append(data[i], row[len(data[i]):]...)
	}
	return data
}
