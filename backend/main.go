package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"backend/postgres"
	// "backend/sampleEndpoint"
	"backend/MPLEndpoint"

	"github.com/gorilla/mux"
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
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(name, " took ", elapsed)
}

func init() {
	envconfig.MustProcess("pbr", &env)
}

func main() {
	// Setup our database connection
	mysqlDB = postgres.NewConnection()

	// Setup our router
	main := mux.NewRouter()
	apiSubrouterPath := "/api"
	routerAPI := main.PathPrefix(apiSubrouterPath).Subrouter()
	routerV1 := routerAPI.PathPrefix("/v1").Subrouter()
	main.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	// Load our endpoints
	// sampleEndpoint.Load(routerV1, mysqlDB)
	MPLEndpoint.Load(routerV1, mysqlDB)

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

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      main, //routerAPI,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.Info("The server is starting, and it will be listening on port " + port)

	// // Prevents memory leak
	server.SetKeepAlivesEnabled(false)

	// HTTP Rest server
	log.Fatal(
		// Serve on the specified port
		server.ListenAndServeTLS("/etc/letsencrypt/live/stuartdehaas.ca/fullchain.pem", "/etc/letsencrypt/live/stuartdehaas.ca/privkey.pem"),
	)
}

//  /etc/letsencrypt/live/stuartdehaas.ca/fullchain.pem
// Your key file has been saved at:
// /etc/letsencrypt/live/stuartdehaas.ca/privkey.pem
// ip 157.245.171.13:443
