package main

import (
	"backend/postgres"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	// "backend/sampleEndpoint"
	"backend/DataValidationEndpoint"
	"backend/MPLEndpoint"
	"backend/TeslaEndpoint"
	"backend/utility"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/common/log"
)

type config struct {
	AppEnv string `envconfig:"APP_ENV" default:"development"`
}

var (
	port    = "443"
	mysqlDB *sqlx.DB
	env     config
	entry   string = "signIn/index.html"               //entrypoint to server, combine with static
	static  string = "/home/cameron/Go_Server/static/" //directory to serve static files from
)

func init() {
	envconfig.MustProcess("pbr", &env)
	fmt.Println(os.Getenv("SESSION_KEY"))
	if os.Getenv("SESSION_KEY") == "" {
		os.Setenv("SESSION_KEY", string(securecookie.GenerateRandomKey(32)))
	}
	fmt.Println(os.Getenv("SESSION_KEY"))
}
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://industrialplankton.com"+r.RequestURI, http.StatusMovedPermanently)
}

func indexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, static+entrypoint)
	}
	return http.HandlerFunc(fn)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, static+"favicon.ico")
}

type authenticationMiddleware struct {
}

//AuthMiddleware checks for current session
func (amw *authenticationMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("session-name")

		if err == nil {
			_, err := r.Cookie(token.Value)
			if err == nil {
				// We found the token in our map
				// log.Printf("Authenticated user %s\n", user)
				// Pass down the request to the next middleware (or final handler)
				next.ServeHTTP(w, r)
			} else {
				http.ServeFile(w, r, static+entry)
			}
		} else {
			//Serve the signIn page
			http.ServeFile(w, r, static+entry)
			// http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func main() {
	// Setup our database connection
	mysqlDB = postgres.NewConnection()

	// Setup our router
	main := mux.NewRouter()
	apiSubrouterPath := "/api"
	routerAPI := main.PathPrefix(apiSubrouterPath).Subrouter()
	routerV1 := routerAPI.PathPrefix("/v1").Subrouter()

	//Serve flavor icon
	main.HandleFunc("/favicon.ico", faviconHandler)

	//Industrial Plankton.com redirect
	main.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		redirectTLS(w, r)
	})
	//Website Test Page, Nonessential
	main.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("The server is running.\n"))
	})

	staticFileDirectory := http.Dir(static)
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for "./assets/assets/index.html", and yield an error
	// staticFileHandler := http.StripPrefix("/assets/")
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	main.Handle("/assets/", http.TimeoutHandler(http.Handler(http.FileServer(staticFileDirectory)), 2*time.Second, "Timeout!\n"))
	main.HandleFunc("/bird", getBirdHandler).Methods("GET")
	main.Handle("/tokenSignIn", http.TimeoutHandler(http.HandlerFunc(tokenSignIn), 2*time.Second, "Timeout!\n")).Methods("POST")
	main.Handle("/tokenSignOut", http.TimeoutHandler(http.HandlerFunc(tokenSignOut), 2*time.Second, "Timeout!\n")).Methods("GET")
	//main.HandleFunc("/birdUp", updateBirdHandler).Methods("POST")

	// Catch-all: Serve our JavaScript application's entry-point (index.html).
	main.PathPrefix("/").HandlerFunc(indexHandler(entry))

	amw := authenticationMiddleware{}
	main.Use(amw.AuthMiddleware)

	// Load our endpoints
	// sampleEndpoint.Load(routerV1, mysqlDB)
	MPLEndpoint.Load(routerV1, mysqlDB)
	TeslaEndpoint.Load(routerV1, mysqlDB)
	DataValidationEndpoint.Load(routerV1, mysqlDB)

	//configure HTTPS Settings
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	//configure server settings
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      main, //routerAPI,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	// // Prevents memory leak
	server.SetKeepAlivesEnabled(false)

	log.Info("The server is starting, and it will be listening on port " + port)

	//HTTP redirect to HTTPS
	go func() {
		log.Fatalf("ListenAndServe error: %v", http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)))
	}()

	// HTTPS Rest server
	err := server.ListenAndServeTLS("/etc/letsencrypt/live/industrialplankton.ca/fullchain.pem", "/etc/letsencrypt/live/industrialplankton.ca/privkey.pem")
	if err != nil {
		utility.Log(err)
		log.Fatalf("ListenAndServeTLS error: %v", err)
	}
}
