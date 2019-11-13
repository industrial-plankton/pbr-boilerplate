package MPLEndpoint

import (
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"
)

// sampleGet - Accepts a GET request
func MasterPartsListHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running get MPL!")
	_, err := RefreshMasterPartsList(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("No parts found!")
		return
	}
	// log.Println("We found: ", MPL)

	// json.NewEncoder(w).Encode(MPL)
	log.Println("Sent JSON response")
}
