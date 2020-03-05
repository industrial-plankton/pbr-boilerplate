package IPDatabase

import (
	"backend/utility"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vishalkuo/bimap"
)

//GetHeaders returns just the headers of the table in a row
func GetHeaders(mysqlDB *sqlx.DB, table string) []interface{} {
	// defer timeTrack(time.Now(), "Get Headers of "+table)
	// fetch all places from the db
	var headers []interface{}
	SQL := "SELECT * FROM information_schema.columns WHERE table_schema = 'public' AND table_name='" + table + "';"
	rows, err := mysqlDB.Queryx(SQL)
	if err != nil {
		utility.Log(err)
		return headers
	}
	// iterate over each row
	for rows.Next() {
		rowdata, _ := rows.SliceScan()
		headers = append(headers, rowdata[3]) //column names are stored in column 4 of information_schema.columns
	}
	return headers
}

//GetView returns entire view with headers
func GetView(mysqlDB *sqlx.DB, view string /*, headers []interface{}*/) [][]interface{} {
	// defer timeTrack(time.Now(), "View "+view)
	SQL := "SELECT * FROM " + view
	values, _ := headerQuery(mysqlDB, SQL, view)

	return values
}

//Search a table and return all rows with the key in the keyColumn
func Search(mysqlDB *sqlx.DB, table, key, keyColumn string) ([][]interface{}, error) {
	// defer timeTrack(time.Now(), "Search for "+key)
	SQL := "SELECT * FROM " + table + " t WHERE t." + keyColumn + "::text='" + key + "'::text"
	values, err := standardQuery(mysqlDB, SQL)
	if err != nil {
		return values, err
	}
	return values, nil
}

//Filter returns the keyColumn filtered by the key
func Filter(mysqlDB *sqlx.DB, table, key, keyColumn string) ([][]interface{}, error) {
	// defer timeTrack(time.Now(), "Search for "+key)
	SQL := "SELECT " + keyColumn + " FROM " + table + " t WHERE t." + keyColumn + " ~* '" + key + "' ORDER BY " + keyColumn
	values, err := standardQuery(mysqlDB, SQL)
	if err != nil {
		return values, err
	}
	return values, nil
}

//MultiLIKE returns all rows in a table matching multiple conditions, built from keys, keyColumns and combiners
func MultiLIKE(mysqlDB *sqlx.DB, table string, keys, keyColumns, combiners []string) ([][]interface{}, error) {
	//KeyColumns should be appended with ::text when relevent
	//Keys should be in the form 'key', '%key%' or '_key_', _ is single wildcard, % is multi wildcard
	defer utility.TimeTrack(time.Now(), "MultiSearch")
	conditions := "("
	for i, e := range keyColumns {
		conditions = conditions + combiners[i] + " " + e + " ~* " + keys[i]
	}
	conditions = conditions + ")"
	SQL := "SELECT * FROM " + table + " AS t WHERE" + conditions
	values, err := standardQuery(mysqlDB, SQL)
	if err != nil {
		return values, err
	}

	return values, nil
}

//Insert puts new data into specified columns of a table
//primary keys are stripped from inputs as they are autogenerated
func Insert(mysqlDB *sqlx.DB, table string, columns []interface{}, data [][]interface{}) error {
	// defer timeTrack(time.Now(), "Insert to "+table)
	primarykey := getPrimaryKeyColumnName(mysqlDB, table)
	for i := range columns { //Cut out primary Key index, it should be autogenerated
		if columns[i] == primarykey {
			columns = append(columns[:i], columns[i+1:]...)
			for d := range data {
				data[d] = append(data[d][:i], data[d][i+1:]...)
			}
			break
		}
	}
	combinedColumns, combinedValues := formatParams(columns, data)
	SQL := "INSERT INTO " +
		table +
		" (" + combinedColumns +
		") VALUES (" +
		combinedValues +
		");"
	// fmt.Print("Insert SQL:" + SQL)
	err := standardExecute(mysqlDB, SQL)
	if err != nil {
		return err
	}
	return nil
}

//Update changes the existing data
func Update(mysqlDB *sqlx.DB, table, primarykey string, columns, data []interface{}) error {
	// defer timeTrack(time.Now(), "Update of "+table)
	setString := ""
	data = utility.ParseNulls(utility.MatchSizes([][]interface{}{data}, columns))[0]
	for i := range columns {
		setString += fmt.Sprint(columns[i]) + "=" + fmt.Sprint(data[i]) + ", "
	}
	setString = strings.TrimSuffix(setString, ", ")
	SQL := "UPDATE " + table + " SET " + setString +
		" WHERE " + getPrimaryKeyColumnName(mysqlDB, table) + "=" + primarykey
	// fmt.Print("Modify SQL:" + SQL)
	err := standardExecute(mysqlDB, SQL)
	if err != nil {
		return err
	}
	return nil
}

//Convert Finds Key in Keycolumn and returns the data in endcolumn
//used to get another columns data
func Convert(mysqlDB *sqlx.DB, table, key, keycolumn, endcolumn string) []interface{} {
	// defer timeTrack(time.Now(), "Convert "+key+" to "+endcolumn)
	// var value []interface{}

	SQL := "SELECT t." + endcolumn + " FROM " + table + " t WHERE t." + keycolumn + "::text='" + key + "'::text"

	value, _ := standardQuery(mysqlDB, SQL)
	if value == nil {
		return nil
	}

	return value[0]
}

//Exists checks if key exists in keycolumn, returns false immediately if no key
func Exists(mysqlDB *sqlx.DB, table, key, keycolumn string) bool {
	// defer timeTrack(time.Now(), "Exist check of "+key+" in "+table+"."+keycolumn)
	if key == "" {
		return false
	}
	var exists bool
	SQL := "select exists(select 1 from " + table + " where " + keycolumn + "=" + key + ")"
	rows, err := mysqlDB.Queryx(SQL)
	if err != nil {
		utility.Log(SQL + " Failed")
		utility.Log(err)
		return exists
	}
	for rows.Next() {
		err = rows.Scan(&exists)
	}
	if err != nil {
		fmt.Println(err)
	}
	return exists
}

//Delete row matching given primaryIndex
func Delete(mysqlDB *sqlx.DB, table, primaryIndex string) error {
	PrimaryKeyColumnName := getPrimaryKeyColumnName(mysqlDB, table)
	// if Exists(mysqlDB, table, primaryIndex, PrimaryKeyColumnName) {
	SQL := "DELETE FROM " + table +
		" WHERE " + PrimaryKeyColumnName + "=" + primaryIndex + ";"
	err := standardExecute(mysqlDB, SQL)
	if err != nil {
		return err
	}
	// }
	return nil
}

//FormatParams arranges columns and data for Inserting
func formatParams(columns []interface{}, data [][]interface{}) (string, string) {
	combinedColumns := ""
	for _, e := range columns {
		combinedColumns = combinedColumns + fmt.Sprintf("%v", e) + ", "
	}
	combinedColumns = strings.TrimRight(combinedColumns, ", ")
	data = utility.MatchSizes(data, columns)
	data = utility.ParseNulls(data)
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

//GetPrimaryKeyColumnName returns primary key for table
func getPrimaryKeyColumnName(mysqlDB *sqlx.DB, table string) string {
	// defer timeTrack(time.Now(), "Get Primary Key name of "+table)
	SQL := "SELECT c.column_name " +
		"FROM information_schema.key_column_usage AS c " +
		"LEFT JOIN information_schema.table_constraints AS t " +
		"ON t.constraint_name = c.constraint_name " +
		"WHERE t.table_name = '" + table + "' AND t.constraint_type = 'PRIMARY KEY';"

	primarykey, _ := standardQuery(mysqlDB, SQL)
	return fmt.Sprint(primarykey[0][0])
}

//StandardQuery executes SQL string as a query and returns the data
func standardQuery(mysqlDB *sqlx.DB, SQL string) ([][]interface{}, error) {
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

//HeaderQuery executes SQL string as a query and returns the data but with the table headers as the first row
//must be selecting all columns, in normal order
func headerQuery(mysqlDB *sqlx.DB, SQL, table string) ([][]interface{}, error) {
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

//StandardExecute executes the SQL string. Use for non-query SQL comands
func standardExecute(mysqlDB *sqlx.DB, SQL string) error { //for not query SQl statements
	utility.Log(SQL)
	tx := mysqlDB.MustBegin()
	// tx.MustExec(SQL)
	_, err := tx.Exec(SQL)
	if err != nil {
		utility.Log(SQL + " FAILED")
		utility.Log(err)
		// fmt.Print(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		utility.Log(SQL + " FAILED")
		utility.Log(err)
		// fmt.Print(err)
		return err
	}
	return nil
}

//TranslateIndexs converts database index values to their corresponding tables data value
//if primary index is not the furthest right column the table name needs to be added to the switch statement and specify the correct column in order to build the map correctly
func TranslateIndexs(mysqlDB *sqlx.DB, translationTables []string, columns []interface{}, keycolumn interface{}, headerMap *bimap.BiMap, data [][]interface{}) [][]interface{} { //turns readable information back into database indexes, and the keycolumn to its proper column name
	maps := make([]*bimap.BiMap, len(translationTables))
	datalocations := make([]int, len(translationTables))
	for i := range maps {
		switch translationTables[i] {
		case "shipments":
			maps[i] = utility.BuildMap(GetView(mysqlDB, translationTables[i]), []int{4, 0}) // build translation maps //if index column not furthest right will need to re-logic this
		default:
			maps[i] = utility.BuildMap(GetView(mysqlDB, translationTables[i]), []int{0}) // build translation maps //if index column not furthest right will need to re-logic this
		}
	}

	for i, e := range columns { //go through columns and save index locations, convert keycolumn to non-index header
		for m := range datalocations {
			if columns[i] == translationTables[m]+"_index" {
				datalocations[m] = i
				continue
			}
		}
		if e == keycolumn {
			columns[i], _ = headerMap.GetInverse(keycolumn)
		}
	}
	for i := range data { //translate index parts back to index value
		for d := range datalocations {
			data[i][datalocations[d]], _ = maps[d].Get(data[i][datalocations[d]])
		}
	}
	return data
}

//UpdateOrAdd adds or updates depending on the primary keys existing
func UpdateOrAdd(mysqlDB *sqlx.DB, table string, headerMap *bimap.BiMap, data [][]interface{}, headerValues []interface{}) error {
	columns := utility.RearrangeHeaders(headerMap, headerValues) //relabel the sheet headers to the database headers, in order of sheet
	keycolumnLoc := utility.FindUnIndexedLocation(table, columns)
	var keycolumn interface{}
	if keycolumnLoc == -1 {
		keycolumn = nil //if table doesn't have a keycolumn -> nil
	} else {
		keycolumn = columns[keycolumnLoc]
	}
	indexLoc := utility.FindPrimIndexLocation(columns)
	translationTables := utility.FindTranslationTables(table, columns)

	if len(translationTables) > 0 { //only translate if needed
		data = TranslateIndexs(mysqlDB, translationTables, columns, keycolumn, headerMap, data)
	} else {
		for i := range columns { //go through columns and convert keycolumn to non-index header
			if columns[i] == keycolumn {
				columns[i], _ = headerMap.GetInverse(keycolumn)
			}
		}
	}

	for _, row := range data { //go through each line
		primaryIndex := fmt.Sprint(row[indexLoc]) //grab index_parts from the collected data
		var newData [][]interface{}
		if Exists(mysqlDB, table, primaryIndex, fmt.Sprint(columns[indexLoc])) { //check that the index exists
			err := Update(mysqlDB, table, primaryIndex, columns, row) // update the database entry
			if err != nil {
				return err
			}
		} else {
			newData = append(newData, row) //add data to newdata collection
		}
		if len(newData) > 0 { //if there is newData
			err := Insert(mysqlDB, table, columns, newData) //commit new data
			if err != nil {
				return err
			}
		}
	}
	return nil
}
