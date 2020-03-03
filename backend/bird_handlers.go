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
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type table struct {
	Data    [][]interface{} `json:"data"`
	Headers []interface{}   `json:"headers"`
}

type auth struct {
	Idtoken string `json:"id_token"`
}

func getMPLHandler(w http.ResponseWriter, r *http.Request) {
	utility.TimeTrack(time.Now(), "getMPL")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	// fmt.Println(r)

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
	// aud := "1095332051856-mgt08ppg80t5je1co4h388kujqu43ia8.apps.googleusercontent.com"
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
	// fmt.Println(session.Values[42])
	// fmt.Println(session.Values[claimSet])

	// Set some session values.
	session.Values["Auth"] = true
	// session.Values[42] = 43

	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(session)

	addCookie(w, "sessionName", claimSet.Sub)
	fmt.Println(w)
	w.WriteHeader(http.StatusOK)
}

func tokenSignOut(w http.ResponseWriter, r *http.Request) {
	//Sub, err := r.Cookie("sessionName")
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
	session, err := store.Get(r, auth.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(session.Values[42])
	// fmt.Println(session.Values[claimSet])

	// Set some session values.
	session.Values["Auth"] = false
	// session.Values[42] = 43
	session.Options.MaxAge = -1
	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// addCookie(w, "sessionName", "")
}

func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, cookie)
}

//AuthMiddleware checks for current session
func AuthMiddleware(next http.Handler) http.Handler {
	utility.TimeTrack(time.Now(), "Auth Midware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		//token, err := r.Cookie("sessionName")
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

		if auth != nil {
			fmt.Println("token.Value:")
			fmt.Println(auth.Value)

			session, _ := store.Get(r, auth.Name)
			if session.Values["Auth"] == true {

				// _, err := r.Cookie(auth.Value)
				//	if err == nil {
				//// We found the token in our map
				//// log.Printf("Authenticated user %s\n", user)
				// Pass down the request to the next middleware (or final handler)
				fmt.Println("Authorized")
				next.ServeHTTP(w, r)
			} else {
				fmt.Println("Not Authorized: ")
				fmt.Println(r)
				http.Redirect(w, r, "https://industrialplankton.ca"+signIn, http.StatusForbidden)
				// http.ServeFile(w, r, static+entry)
			}
		} else {
			//Serve the signIn page
			fmt.Println("Not Authorized: ")
			fmt.Println(r)
			// http.ServeFile(w, r, static+entry)
			http.Redirect(w, r, "https://industrialplankton.ca"+signIn, http.StatusForbidden)
			//// http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
