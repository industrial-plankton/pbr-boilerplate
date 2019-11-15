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

func FindPartForEditHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running get Part!")
	part, err := FindPartForEdit(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("No parts found!")
		return
	}
	log.Println("We found: ", part)
}

func SaveMPLEditHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running save MPL!")
	part, err := SaveMPLEdit(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("Save failed!")
		return
	}
	log.Println("We saved: ", part)
}
