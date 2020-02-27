package main

import (
	"backend/IPDatabase"
	"backend/utility"
	"bytes"
	"fmt"
	json "github.com/json-iterator/go" //"encoding/json"
	"net/http"
)

type table struct {
	Data    [][]interface{} `json:"data"`
	Headers []interface{}   `json:"headers"`
}

type auth struct {
	Idtoken string `json:"id_token"`
}

func getBirdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	fmt.Println(r)

	//View := IPDatabase.GetView(mysqlDB, "masterpartslist")
	View := IPDatabase.GetView(mysqlDB, "column_map")

	prejsonObj := table{Headers: View[0], Data: View[1:]}

	// Everything else is the same as before
	jsonObj, err := json.Marshal(prejsonObj)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonObj)
}

func tokensignin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	var token auth
	utility.ParseJSON(w, r, token)
	fmt.Println("ResponseWriter: ")
	fmt.Println(w)
	fmt.Println("Token: ")
	fmt.Println(token)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	newStr := buf.String()

	fmt.Println("r.Body: ", string(newStr))

	fmt.Println("request: ")
	fmt.Println(r)

	w.WriteHeader(http.StatusAccepted)

	fmt.Println("Validation URL: ")
	fmt.Println("https://oauth2.googleapis.com/tokeninfo?id_token=" + token.Idtoken)
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + token.Idtoken)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println("ID validation response: ")
	fmt.Println(resp)

	fmt.Println("ResponseWriter: ")
	fmt.Println(w)
}
