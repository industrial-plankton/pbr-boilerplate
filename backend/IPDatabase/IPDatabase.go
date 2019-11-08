package IPDatabase

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

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

func GetMPL(mysqlDB *sqlx.DB) [][]interface{} {
	// defer timeTrack(time.Now(), "master2")
	// fetch all places from the db
	var values [][]interface{}
	rows, _ := mysqlDB.Queryx("SELECT parts.sku, parts.technical_desc, parts.customer_desc FROM postgres.parts parts ORDER BY parts.index_parts")
	// iterate over each row
	for rows.Next() {
		var tdisc string
		var cdisc string
		var SKU string
		//var mVen string
		// var MPN string
		// use the NullString type if NULLABLE
		_ = rows.Scan(&SKU, &tdisc, &cdisc)         //, &mVen, &MPN)
		rowdata := []interface{}{SKU, tdisc, cdisc} //, mVen, MPN}
		values = append(values, rowdata)
	}

	// var values2 [][]interface{}
	// rows, _ = mysqlDB.Queryx("SELECT `parts`.`IP SKU`, `vendors`.`name`, `parts`.`Secondary supplier PN`, `parts`.`Extra Info`, `Order Type`.`Order Type`, `unit`.`unit` FROM `demodb`.`parts` `parts`, `demodb`.`vendors` `vendors`, `demodb`.`Order Type` `Order Type`, `demodb`.`unit` `unit` WHERE `parts`.`Supplier (Secondary)` = `vendors`.`index_Ven` AND `parts`.`Order Type` = `Order Type`.`index_OType` AND `parts`.`units` = `unit`.`index_unit` ORDER BY `parts`.`index_Parts`")

	// // iterate over each row
	// for rows.Next() {
	// 	var SKU string
	// 	var sVen string
	// 	var SPN string
	// 	var ExtraInf string
	// 	var OT string
	// 	var unit string

	// 	// note that city can be NULL, so we use the NullString type
	// 	_ = rows.Scan(&SKU, &sVen, &SPN, &ExtraInf, &OT, &unit)
	// 	rowdata := []interface{}{sVen, SPN, ExtraInf, OT, unit}
	// 	values2 = append(values2, rowdata)
	// }

	// for index, element := range values2 {
	// 	// index is the index where we are
	// 	// element is the element from someSlice for where we are
	// 	values[index] = append(values[index], element...)
	// }

	return values
}
