package main

import (
	"backend/IPDatabase"
	"backend/utility"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/sessions"
	json "github.com/json-iterator/go" //

	//"encoding/json"
	verify "github.com/futurenda/google-auth-id-token-verifier"
)

//so the package is used no mater what
var _ = verify.Certs{}

// Initialize the Session Store
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type table struct {
	Data    [][]interface{} `json:"data"`
	Headers []interface{}   `json:"headers"`
}

type auth struct {
	Idtoken string `json:"id_token"`
}

type keywords struct {
	Descriptions string `json:"Descriptions"`
	Suppliers    string `json:"Suppliers"`
}

func keyWordSearch(w http.ResponseWriter, r *http.Request) {
	var toSearch keywords
	var keysSlice []string

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&toSearch)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	keysSlice = append(keysSlice, toSearch.Descriptions, toSearch.Suppliers)
	for len(keysSlice) < 2 { //append the empty keywords
		keysSlice = append(keysSlice, "")
	}

	keysSlice = utility.AddWildCards(keysSlice)
	keysSlice = []string{keysSlice[0], keysSlice[0], keysSlice[1], keysSlice[1]} //duplicate entries for multicolumn search

	keycolumns := []string{"technical_desc", "customer_desc", "name", "part_number"}
	combiners := []string{"", " OR ", ") AND (", " OR "}
	SearchResult, err := IPDatabase.MultiLIKE(mysqlDB, "keywordsearch", keysSlice, keycolumns, combiners)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Headers := IPDatabase.GetHeadersPrepared(mysqlDB, "keywordsearch")
	for i, e := range Headers {
		Headers[i] = IPDatabase.Convert(mysqlDB, "column_map", fmt.Sprint(e), "", endcolumn)
	}
	fmt.Println(len(Headers))
	//	//SearchResult = utility.PopColumn(SearchResult, 0)

	//Parse No result
	if len(SearchResult) == 0 {
		var empty []interface{}
		SearchResult = append(SearchResult, empty)
	}
	if len(SearchResult[0]) == 0 {
		SearchResult[0] = append(SearchResult[0], "No Match")
	}

	//Pack JSON
	prejsonObj := table{Headers: Headers, Data: SearchResult}
	jsonObj, err := json.Marshal(prejsonObj)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonObj)
}

func getMPLHandler(w http.ResponseWriter, r *http.Request) {
	utility.TimeTrack(time.Now(), "getMPL")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	Headers := IPDatabase.GetHeadersPrepared(mysqlDB, "mpltext")
	fmt.Println(Headers)
	View := IPDatabase.GetView(mysqlDB, "masterpartslist")
	// View := IPDatabase.GetView(mysqlDB, "column_map")

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

func tokenSignIn(w http.ResponseWriter, r *http.Request) {
	defer utility.TimeTrack(time.Now(), "tokenSignIn")
	fmt.Println(r)

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
		fmt.Println(err)
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	v := verify.Verifier{}
	err = v.VerifyIDToken(token.Idtoken, []string{
		aud,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claimSet, _ := verify.Decode(token.Idtoken)
	// claimSet.Iss,claimSet.Email ... (See claimset.go)

	if claimSet.EmailVerified && strings.Contains(claimSet.Email, authorizedDomain) {
		//You are verified as part of @industrialplankton.com
	} else {
		fmt.Println(err)
		http.Error(w, "You are not part of our domain", http.StatusUnauthorized)
		return
	}

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, claimSet.Sub)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set as Authorized
	session.Values["Auth"] = true

	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func tokenSignOut(w http.ResponseWriter, r *http.Request) {
	// find the cookie that is just a number as it is the session cookie
	cookies := r.Cookies()
	var auth *http.Cookie
	for _, e := range cookies {
		_, err := strconv.ParseFloat(e.Name, 64)
		if err != nil {
			continue
		}
		auth = e
		break
	}
	if auth == nil {
		http.Error(w, "Your already logged out", http.StatusInternalServerError)
		return
	}

	// Get session
	session, err := store.Get(r, auth.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reset some authorization.
	session.Values["Auth"] = false
	session.Options.MaxAge = -1

	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//// addCookie(w, "sessionName", "")
}

// /**/
// // func addCookie(w http.ResponseWriter, name string, value string) {
// // 	expire := time.Now().AddDate(0, 0, 1)
// // 	cookie := &http.Cookie{
// // 		Name:    name,
// // 		Value:   value,
// // 		Expires: expire,
// // 	}
// // 	http.SetCookie(w, cookie)
// // }

//AuthMiddleware checks for current session
func AuthMiddleware(next http.Handler) http.Handler {
	utility.TimeTrack(time.Now(), "Auth Midware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

		// find the cookie that is just a number as it is the session cookie
		cookies := r.Cookies()
		var auth *http.Cookie
		for _, e := range cookies {
			_, err := strconv.ParseFloat(e.Name, 64)
			if err != nil {
				continue
			}
			auth = e
			break
		}

		//Check Authorization
		if auth != nil {
			session, _ := store.Get(r, auth.Name)

			if session.Values["Auth"] == true {
				// Pass down the request to the next middleware (or final handler)
				fmt.Println(r.RemoteAddr + " Authorized for " + r.RequestURI)
				next.ServeHTTP(w, r)
			} else {
				fmt.Println(r.RemoteAddr + " Not Authorized for" + r.RequestURI)
				http.Redirect(w, r, "https://industrialplankton.ca"+signIn, http.StatusForbidden)
			}
		} else {
			//Serve the signIn page
			fmt.Println(r.RemoteAddr + " Not Authorized for" + r.RequestURI)
			http.Redirect(w, r, "https://industrialplankton.ca"+signIn, http.StatusForbidden)
		}
	})
}
