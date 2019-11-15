package IPDatabase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(name, " took ", elapsed)
}

func GetUnits(mysqlDB *sqlx.DB) {
	// defer timeTrack(time.Now(), "Unit")
	// fetch all places from the db
	var values [][]string
	rows, _ := mysqlDB.Query("SELECT `parts`.`index_Parts`, `unit`.`unit`, `parts`.`IP SKU` FROM `demodb`.`parts` `parts`, `demodb`.`unit` `unit` WHERE `parts`.`units` = `unit`.`index_unit` ORDER BY `parts`.`index_Parts`")

	// iterate over each row
	for rows.Next() {
		var units string
		var index_Parts int
		var SKU string

		// note that city can be NULL, so we use the NullString type
		_ = rows.Scan(&index_Parts, &units, &SKU)
		rowdata := []string{units, SKU, strconv.Itoa(index_Parts)}
		values = append(values, rowdata)
	}
	sliceofslices := values[:]
	fmt.Println(sliceofslices, "\n")
}
func GetHeaders(mysqlDB *sqlx.DB, table string) []interface{} { //DOESNT WORK
	defer timeTrack(time.Now(), "Get Headers of "+table)
	results := make(map[string]interface{})
	var interfaceCol []interface{}
	var headers []interface{}

	err := mysqlDB.QueryRowx("SELECT * FROM " + table).MapScan(results)
	if err != nil {
		fmt.Println(err)
		return headers
	}

	for _, i := range results {
		interfaceCol = []interface{}{i}
		headers = append(headers, interfaceCol)
	}
	return headers
}

func GetView(mysqlDB *sqlx.DB, view string, headers []interface{}) [][]interface{} {
	defer timeTrack(time.Now(), "View "+view)
	// fetch all places from the db
	var values [][]interface{}
	values = append(values, headers)
	rows, _ := mysqlDB.Queryx("SELECT * FROM " + view)

	// iterate over each row
	for rows.Next() {
		rowdata, _ := rows.SliceScan()
		values = append(values, rowdata)
	}

	return values
}

func Search(mysqlDB *sqlx.DB, table string, key string, keycolumn string) [][]interface{} {
	defer timeTrack(time.Now(), "Search for "+key)
	var values [][]interface{}

	rows, err := mysqlDB.Queryx("SELECT * FROM " + table + " t WHERE t." + keycolumn + "::text='" + key + "'::text")
	if err != nil {
		fmt.Println(err)
		return values
	}

	for rows.Next() {
		rowdata, _ := rows.SliceScan()
		values = append(values, rowdata)
	}

	return values
}
