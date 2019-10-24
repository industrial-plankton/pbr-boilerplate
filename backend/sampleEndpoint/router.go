package sampleEndpoint

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

// Load mounts the subrouter on the router and matches each path with a handler
func Load(router *mux.Router, mysqlDB *sqlx.DB) {
	db = mysqlDB
	router.HandleFunc("/sampleEndpoint", sampleGet).Methods(http.MethodGet)
	router.HandleFunc("/sampleEndpoint", samplePost).Methods(http.MethodPost)
}
