package IPDatabase

import (
	"backend/utility"
	"fmt"
	"strconv"
	"strings"
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

func GetHeaders(mysqlDB *sqlx.DB, table string) []interface{} {
	defer timeTrack(time.Now(), "Get Headers of "+table)
	// fetch all places from the db
	var headers []interface{}
	SQL := "SELECT * FROM information_schema.columns WHERE table_schema = 'public' AND table_name='" + table + "';"
	rows, err := mysqlDB.Queryx(SQL)
	if err != nil {
		utility.Log(err)
		// fmt.Println(err)
		return headers
	}
	// iterate over each row
	for rows.Next() {
		rowdata, _ := rows.SliceScan()
		headers = append(headers, rowdata[3]) //column names are stored in column 4 of information_schema.columns
	}
	utility.Log("Headers found: ")
	utility.Log(headers)
	// fmt.Print("Headers found: ")
	// fmt.Println(headers)
	return headers
}

func GetView(mysqlDB *sqlx.DB, view string /*, headers []interface{}*/) [][]interface{} {
	defer timeTrack(time.Now(), "View "+view)
	SQL := "SELECT * FROM " + view
	values, _ := headerQuery(mysqlDB, SQL, view)
	// fetch all places from the db
	// var values [][]interface{}
	// values = append(values, headers)
	// rows, err := mysqlDB.Queryx("SELECT * FROM " + view)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return values
	// }
	// // iterate over each row
	// for rows.Next() {
	// 	rowdata, _ := rows.SliceScan()
	// 	values = append(values, rowdata)
	// }

	return values
}

func Search(mysqlDB *sqlx.DB, table string, key string, keycolumn string) [][]interface{} {
	defer timeTrack(time.Now(), "Search for "+key)
	// var values [][]interface{}
	SQL := "SELECT * FROM " + table + " t WHERE t." + keycolumn + "::text='" + key + "'::text"
	values, _ := standQuery(mysqlDB, SQL)
	// rows, err := mysqlDB.Queryx("SELECT * FROM " + table + " t WHERE t." + keycolumn + "::text='" + key + "'::text")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return values
	// }

	// for rows.Next() {
	// 	rowdata, _ := rows.SliceScan()
	// 	values = append(values, rowdata)
	// }

	return values
}

func MultiLIKE(mysqlDB *sqlx.DB, table string, keys []string, keycolumns []string, combiners []string) [][]interface{} {
	//KeyColumns should be appended with ::text when relevent
	//Keys should be in the form 'key', '%key%' or '_key_', _ is single wildcard, % is multi wildcard
	defer timeTrack(time.Now(), "MultiSearch")
	// var values [][]interface{}
	conditions := "("
	for i, e := range keycolumns {
		conditions = conditions + combiners[i] + " LOWER(" + e + ") ~~ LOWER(" + keys[i] + ")"
	}
	conditions = conditions + ")"
	fmt.Print(conditions)
	SQL := "SELECT * FROM " + table + " AS t WHERE" + conditions
	values, _ := standQuery(mysqlDB, SQL)
	// rows, err := mysqlDB.Queryx("SELECT * FROM " + table + " AS t WHERE" + conditions)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return values
	// }

	// for rows.Next() {
	// 	rowdata, _ := rows.SliceScan()
	// 	values = append(values, rowdata)
	// }

	return values
}

func AddWildCards(array []string) []string {
	for i, e := range array {
		array[i] = "'%" + strings.ReplaceAll(e, " ", "%") + "%'"
	}
	return array
}

func Insert(mysqlDB *sqlx.DB, table string, columns []interface{}, data [][]interface{}) {
	defer timeTrack(time.Now(), "Insert to "+table)
	combinedColumns, combinedValues := formatParams(columns, data)
	SQL := "INSERT INTO " +
		table +
		" (" + combinedColumns +
		") VALUES (" +
		combinedValues +
		");"
	fmt.Print("Insert SQL:" + SQL)
	// tx := db.MustBegin()
	// tx.MustExec(SQL)
	// err := tx.Commit()
	// if err != nil {
	// 	fmt.Print(err)
	// }
}

func Update(mysqlDB *sqlx.DB, table string, primarykey string, columns []interface{}, data []interface{}) {
	defer timeTrack(time.Now(), "Update of "+table)
	SQL := "UPDATE " + table + " SET " +
		" WHERE " + getPrimaryKeyColumnName(mysqlDB, table) + "=" + primarykey
	fmt.Print("Modify SQL:" + SQL)
}

func Convert(mysqlDB *sqlx.DB, table string, key string, keycolumn string, endcolumn string) []interface{} {
	defer timeTrack(time.Now(), "Convert "+key+" to "+endcolumn)
	// var value []interface{}

	SQL := "SELECT t." + endcolumn + " FROM " + table + " t WHERE t." + keycolumn + "::text='" + key + "'::text"

	value, _ := standQuery(mysqlDB, SQL)
	if value == nil {
		return nil
	}
	// rows, err := mysqlDB.Queryx(SQL)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return value
	// }

	// for rows.Next() {
	// 	value, _ = rows.SliceScan()
	// }

	return value[0]
}

func Exists(mysqlDB *sqlx.DB, table string, key string, keycolumn string) bool {
	defer timeTrack(time.Now(), "Exist check of "+key+" in "+table+"."+keycolumn)
	var exists bool
	rows, err := mysqlDB.Queryx("select exists(select 1 from " + table + " where " + keycolumn + "=" + key)
	for rows.Next() {
		err = rows.Scan(&exists)
	}
	if err != nil {
		fmt.Println(err)
	}
	return exists
}

func formatParams(columns []interface{}, data [][]interface{}) (string, string) {

	combinedColumns := ""
	for _, e := range columns {
		combinedColumns = combinedColumns + fmt.Sprintf("%v", e) + ", "
	}
	combinedColumns = strings.TrimRight(combinedColumns, ", ")

	combinedValues := ""
	for _, rowData := range data {
		for _, val := range rowData {
			combinedValues = combinedValues + fmt.Sprintf("%v", val) + ", "
		}
		combinedValues = strings.TrimRight(combinedValues, ", ") + "),("
	}
	combinedValues = strings.TrimRight(combinedValues, "),(")

	return combinedColumns, combinedValues
}

func getPrimaryKeyColumnName(mysqlDB *sqlx.DB, table string) string {
	defer timeTrack(time.Now(), "Get Primary Key name of "+table)

	SQL := "SELECT a.attname " +
		"FROM pg_index i " +
		"JOIN pg_attribute a ON a.attrelid = i.indrelid " +
		"AND a.attnum = ANY(i.indkey)" +
		"WHERE  i.indrelid = '" + table + "'::regclass" +
		"AND    i.indisprimary;"

	primarykey, _ := standQuery(mysqlDB, SQL)

	return fmt.Sprintf("%v", primarykey)
}

func standQuery(mysqlDB *sqlx.DB, SQL string) ([][]interface{}, error) {
	var values [][]interface{}
	rows, err := mysqlDB.Queryx(SQL)
	if err != nil {
		utility.Log(SQL + " Failed")
		utility.Log(err)
		// fmt.Println(err)
		return values, err
	}

	for rows.Next() {
		rowdata, _ := rows.SliceScan()
		values = append(values, rowdata)
	}
	return values, nil
}

func headerQuery(mysqlDB *sqlx.DB, SQL string, table string) ([][]interface{}, error) { //must be selecting all columns in normal order
	var values [][]interface{}

	values = append(values, GetHeaders(mysqlDB, table))
	rows, err := mysqlDB.Queryx(SQL)
	if err != nil {
		utility.Log(SQL + " Failed")
		utility.Log(err)
		// fmt.Println(err)
		return values, err
	}

	for rows.Next() {
		rowdata, _ := rows.SliceScan()
		values = append(values, rowdata)
	}
	return values, nil
}

func standExec(mysqlDB *sqlx.DB, SQL string) error {
	tx := mysqlDB.MustBegin()
	tx.MustExec(SQL)
	err := tx.Commit()
	if err != nil {
		utility.Log(SQL + " FAILED")
		utility.Log(err)
		// fmt.Print(err)
	}
	return err
}
