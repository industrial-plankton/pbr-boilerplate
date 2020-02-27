package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang/gddo/httputil/header"
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
			if len(data[i]) > len(CombinedData) { //append a row of nulls if CombinedData doesnt have another row
				nullrow := make([]interface{}, len(CombinedData[0]))
				CombinedData = append(CombinedData, nullrow)
			}
			CombinedData[r] = append(CombinedData[r], data[i][r]...) //add the contents of the new row to CombinedData
		}
	}
	return CombinedData
}

//ParseJSON , a heavily error checking JSON Parser
func ParseJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we are using the gddo/httputil/header
	// package to parse and extract the value here, so the check works
	// even if the client includes additional charset or boundary
	// information in the header.
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Check that the request body only contained a single JSON object.
	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "JSON Data: %+v", data)
}
