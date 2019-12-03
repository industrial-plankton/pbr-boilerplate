package TeslaEndpoint

import (
	// "net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var (
	db    *sqlx.DB
	IDMap map[string]string
)

func init() {
	//map user emails to their personal sheet ID
	IDMap = make(map[string]string)
	IDMap["cam@industrialplankton.com"] = "1x4hcf5J0haQN_1MrxUPR1T2gCxXEpwpK5RWMTrXJLbY" //Databased ZERO
}

// Load mounts the subrouter on the router and matches each path with a handler
func Load(router *mux.Router, mysqlDB *sqlx.DB) {
	db = mysqlDB
	router.HandleFunc("/Tesla", FindShipForEditHandle).Headers("RequestType", "findShipForEdit") //.Methods(http.MethodGet)
	router.HandleFunc("/Tesla", SaveTeslaEditHandle).Headers("RequestType", "SaveTeslaEdit")     //.Methods(http.MethodGet)
	router.HandleFunc("/Tesla", SearchPOsHandle).Headers("RequestType", "SearchPOs")
}
