package main

import (
	"backend/IPDatabase"
	//"backend/utility"
	"fmt"
	json "github.com/json-iterator/go" //"encoding/json"
	"net/http"
)

type Table struct {
	Data    [][]interface{} `json:"data"`
	Headers []interface{}   `json:"headers"`
}

func getBirdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	View := IPDatabase.GetView(mysqlDB, "masterpartslist")
	//View := IPDatabase.GetView(mysqlDB, "column_map")

	prejsonObj := Table{Headers: View[0], Data: View[1:]}

	// Everything else is the same as before
	jsonObj, err := json.Marshal(prejsonObj)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonObj)
}
