package main

import (
	"backend/IPDatabase"
	"backend/utility"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/gddo/httputil/header"
	json "github.com/json-iterator/go" //

	//"encoding/json"
	verify "github.com/futurenda/google-auth-id-token-verifier"
)

//so the package is used no mater what
var _ = verify.Certs{}

type table struct {
	Data    [][]interface{} `json:"data"`
	Headers []interface{}   `json:"headers"`
}

type auth struct {
	Idtoken string `json:"id_token"`
}

//AuthorizedDomain , the only sign ins we should give authorization for are emails wth  @industrialplankton.com
const AuthorizedDomain = "@industrialplankton.com"

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
	defer utility.TimeTrack(time.Now(), "tokensignin")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

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

	var token auth
	err := dec.Decode(&token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check that the request body only contained a single JSON object.
	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	v := verify.Verifier{}
	aud := "1095332051856-mgt08ppg80t5je1co4h388kujqu43ia8.apps.googleusercontent.com"
	err = v.VerifyIDToken(token.Idtoken, []string{
		aud,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claimSet, _ := verify.Decode(token.Idtoken)
	// claimSet.Iss,claimSet.Email ... (See claimset.go)

	if claimSet.EmailVerified && strings.Contains(claimSet.Email, AuthorizedDomain) {
		//You are verifiyed as part of @industrialplankton.com
	} else {
		http.Error(w, "You are not part of our domain", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

}
