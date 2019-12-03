package DataValidationEndpoint

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
	IDMap["cam@industrialplankton.com"] = "1yDBlOAAGNPeHSCFMOghWvDCa6crJrgFdNyf6JsvbS10" //Databased DataValidation
}

// Load mounts the subrouter on the router and matches each path with a handler
func Load(router *mux.Router, mysqlDB *sqlx.DB) {
	db = mysqlDB
	router.HandleFunc("/Validation", RefreshHandler).Headers("RequestType", "Refresh") //.Methods(http.MethodGet)

}
