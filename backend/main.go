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
	mysqlDB *sqlx.DB
	env     config
	// store   *sessions.CookieStore
)

const (
	port             string = "443"
	aud              string = "1095332051856-mgt08ppg80t5je1co4h388kujqu43ia8.apps.googleusercontent.com"
	entry            string = "/signIn/index.html"             //entrypoint to server, combine with static
	static           string = "/home/cameron/Go_Server/static" //directory to serve static files from
	signIn           string = "/signIn/"
	authorizedDomain string = "@industrialplankton.com"
)

func init() {
	envconfig.MustProcess("pbr", &env)
	fmt.Println(os.Getenv("SESSION_KEY"))
	if os.Getenv("SESSION_KEY") == "" {
		os.Setenv("SESSION_KEY", string(securecookie.GenerateRandomKey(32)))
	}
	fmt.Println(os.Getenv("SESSION_KEY"))
	// store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	utility.TimeTrack(time.Now(), "http redirect")
	http.Redirect(w, r, "https://industrialplankton.com"+r.RequestURI, http.StatusMovedPermanently)
}

func indexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	utility.TimeTrack(time.Now(), "Catch all handle")
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://industrialplankton.ca"+signIn, http.StatusNotFound)
	}
	return http.HandlerFunc(fn)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	utility.TimeTrack(time.Now(), "Serve Icon")
	http.ServeFile(w, r, static+"favicon.ico")
}

func timedHandler(Handler http.HandlerFunc, timeOut int) http.Handler {
	return http.TimeoutHandler(http.HandlerFunc(Handler), time.Duration(timeOut)*time.Second, "Timeout!\n")
}

func main() {
	// Setup our database connection
	mysqlDB = postgres.NewConnection()

	// Setup our router
	main := mux.NewRouter()
	apiSubrouterPath := "/api"
	routerAPI := main.PathPrefix(apiSubrouterPath).Subrouter()
	routerV1 := routerAPI.PathPrefix("/v1").Subrouter()
	authRouter := main.Host("industrialplankton.ca").PathPrefix("/auth").Subrouter()
	Protected := main.Host("industrialplankton.ca").PathPrefix("/prot").Subrouter()
	Protected.Use(AuthMiddleware)
	ProtectedAPI := Protected.PathPrefix(apiSubrouterPath).PathPrefix("/v1").Subrouter()
	//catchAll := main.PathPrefix("/").Subrouter()

	//Serve flavor icon
	main.HandleFunc("/favicon.ico", faviconHandler)

	//IndustrialPlankton.com redirect
	main.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(".com redirect")
		redirectTLS(w, r)
	})
	//Website Test Page, Nonessential
	main.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("The server is running.\n"))
	})
	main.Handle("/signIn/", http.TimeoutHandler(http.Handler(http.FileServer(http.Dir(static))), 2*time.Second, "Timeout!\n"))

	// //staticFileDirectory := http.Dir(static)
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for "./assets/assets/index.html", and yield an error
	// staticFileHandler := http.StripPrefix("/assets/")
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself

	//	// Protected.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
	//	// 	(http.ServeFile(w, r, static+r.RequestURI))
	//	// })
	Protected.Handle("/assets/", http.TimeoutHandler(http.Handler(http.FileServer(http.Dir(static))), 2*time.Second, "Timeout!\n"))
	ProtectedAPI.Handle("/MPL", timedHandler(getMPLHandler, 5)).Methods("GET")
	authRouter.Handle("/tokenSignIn", timedHandler(tokenSignIn, 2)).Methods("POST")
	authRouter.Handle("/tokenSignOut", timedHandler(tokenSignOut, 2)).Methods("GET")

	////main.HandleFunc("/birdUp", updateBirdHandler).Methods("POST")

	// * Catch-all: Serve our JavaScript application's entry-point (index.html).
	main.PathPrefix("/").HandlerFunc(indexHandler(entry))

	// Load our endpoints
	//// sampleEndpoint.Load(routerV1, mysqlDB)
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
	// Prevents memory leak
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
