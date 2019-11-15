package MPLEndpoint

import (
	// "net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var (
	db    *sqlx.DB
	MPLID string = "1Hi0PrHe53q4JhNetcJ_y3WrDIJ9qocVEd4irMunxVyE"
)

// Load mounts the subrouter on the router and matches each path with a handler
func Load(router *mux.Router, mysqlDB *sqlx.DB) {
	db = mysqlDB
	// MPLID = "1Hi0PrHe53q4JhNetcJ_y3WrDIJ9qocVEd4irMunxVyE"
	router.HandleFunc("/MPLEndpoint", MasterPartsListHandle).Headers("RequestType", "MPLrefresh")      //.Methods(http.MethodGet)
	router.HandleFunc("/MPLEndpoint", FindPartForEditHandle).Headers("RequestType", "findPartForEdit") //.Methods(http.MethodGet)
	router.HandleFunc("/MPLEndpoint", SaveMPLEditHandle).Headers("RequestType", "SaveMPLEdit")         //.Methods(http.MethodGet)

}
