package main

import (
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
	port    = "3030"
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
	main.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("The server is running.\n"))
	})

	// Load our endpoints
	// sampleEndpoint.Load(routerV1, mysqlDB)
	MPLEndpoint.Load(routerV1, mysqlDB)

	log.Info("The server is starting, and it will be listening on port " + port)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      main, //routerAPI,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Prevents memory leak
	server.SetKeepAlivesEnabled(false)

	// HTTP Rest server
	log.Fatal(
		// Serve on the specified port
		server.ListenAndServe(),
	)
}
