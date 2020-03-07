package utility

import (
	"fmt"
	"strings"
	"time"

	"github.com/vishalkuo/bimap"
)

//IntfToString convert an interface slice to string slice
func IntfToString(data []interface{}) []string {
	out := make([]string, len(data))
	for i, e := range data {
		out[i] = fmt.Sprintf("%v", e)
	}
	return out
}

//AddWildCards adds regular expresstion wildcars to begining and end of each word
func AddWildCards(array []string) []string {
	for i, e := range array {
		term := strings.ReplaceAll(e, " ", ".*)(.*")
		array[i] = "'(.*" + strings.TrimSuffix(term, "\n") + ".*)'"
	}
	return array
}

//SanatizeSQLVariables adds regular expresstion wildcars to begining and end of each word
func SanitizeSQLVariables(array []string) []string {
	for i, e := range array {
		array[i] = SanitizeSQLVariable(e)
	}
	return array
}

//SanatizeSQLVariables removes all " and ' from Variable
func SanitizeSQLVariable(Variable string) string {
	Variable = strings.ReplaceAll(Variable, "'", "''")
	Variable = strings.ReplaceAll(Variable, `"`, `""`)
	return Variable
}

//todo make function to turn array into (0)(1)(2) for and conditions or (0)|(1)|(2) for or
func RegxAndString(variables []string) string {
	var buffer strings.Builder
	for _, e := range variables {
		buffer.WriteString("(%")
		buffer.WriteString(e)
		buffer.WriteString("%)")
	}
	return buffer.String()
}

func RegxOrString(variables []string) string {
	var buffer strings.Builder
	for _, e := range variables {
		buffer.WriteString("(%")
		buffer.WriteString(e)
		buffer.WriteString("%)|")
	}
	return strings.Trim(buffer.String(), "|")
}

//RearrangeHeaders rearanges and converts Header slice into database header sclice
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

//BuildMap Creates a map of Text and Index, assumes index is farthest right if not specified
func BuildMap(data [][]interface{}, colIndex []int) *bimap.BiMap {
	biMap := bimap.NewBiMap()
	if len(colIndex) == 1 {
		colIndex = append(colIndex, len(data[0])-1)
	}
	for i := range data {
		biMap.Insert(data[i][colIndex[0]], data[i][colIndex[1]])
	}
	return biMap
}

//ParseNulls puts '' around data, replaces empty spots with NULL
func ParseNulls(data [][]interface{}) [][]interface{} {
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

//FindPrimIndexLocation locates the position of the primary Key (has prefix "index_")
func FindPrimIndexLocation(columns []interface{}) int {
	for i, e := range columns {
		if strings.HasPrefix(fmt.Sprint(e), "index_") {
			return i
		}
	}
	return -1
}

//FindUnIndexedLocation locates the position of the Key (table+"_index")
func FindUnIndexedLocation(table string, columns []interface{}) int {
	for i, e := range columns {
		if e == table+"_index" {
			return i
		}
	}
	return -1
}

//FindTranslationTables determines which translation tables need to be pulled from the database
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

//GetHeaderLocation returns the location of a string inside the slice
func GetHeaderLocation(columns []interface{}, header string) int {
	for i, e := range columns {
		if e == header {
			return i
		}
	}
	return -1
}

//OverWriteColumn fill a colum with a value
func OverWriteColumn(data [][]interface{}, value interface{}, column int) [][]interface{} {
	for i := range data {
		data[i][column] = value
	}
	return data
}

//FillifEmpty fills nil and "" points with "value"
func FillifEmpty(data [][]interface{}, value interface{}, column int) [][]interface{} {
	for i := range data {
		if data[i][column] != nil && data[i][column] != "" {
			continue
		}
		data[i][column] = value
	}
	return data
}

//MatchSizes ensures rectangle interface by adding nulls
func MatchSizes(data [][]interface{}, size []interface{}) [][]interface{} {
	row := make([]interface{}, len(size))
	for i := range data {
		if len(data[i]) < len(size) {
			data[i] = append(data[i], row[len(data[i]):]...)
		} else {
			data[i] = data[i][:len(size)]
		}
	}
	return data
}

//SetSize enforces data[][] to be rectangular, with a width of size
func SetSize(data [][]interface{}, size int) [][]interface{} {
	row := make([]interface{}, size)
	for i := range data {
		if len(data[i]) < size {
			data[i] = append(data[i], row[len(data[i]):]...)
		} else {
			data[i] = data[i][:size]
		}

	}
	return data
}

//ConcatSplitData , if data for one table is split accross multiple ranges this will combine them for database entry, only use rectangular matrixs
func ConcatSplitData(data [][][]interface{}) [][]interface{} {
	CombinedData := data[0] //initialize CombinedData with the first range

	for i := 1; i < len(data); i++ { //loop through remaining ranges
		for r := range data[i] { //loop through each row
			if len(data[i]) > len(CombinedData) { //append a row of nulls if CombinedData doesn't have another row
				nullrow := make([]interface{}, len(CombinedData[0]))
				CombinedData = append(CombinedData, nullrow)
			}
			CombinedData[r] = append(CombinedData[r], data[i][r]...) //add the contents of the new row to CombinedData
		}
	}
	return CombinedData
}

//PopColumn Removes the specified column from the table data
func PopColumn(table [][]interface{}, column int) [][]interface{} {
	//Nil checks
	if len(table) > 0 {
		if len(table[0]) > 0 {
			var newTable [][]interface{}
			//Check column is in range
			if column >= len(table[0]) {
				fmt.Println("column out of bounds")
				return table
			}
			for _, e := range table {
				var newColumn []interface{}
				if column == 0 {
					newColumn = e[1:]
				} else if column == len(table[0])-1 {
					newColumn = e[:column-1]
				} else {
					newColumn = append(e[:column], e[column+1:])
				}
				newTable = append(newTable, newColumn)
			}
			return newTable
		}
	}
	return table
}
