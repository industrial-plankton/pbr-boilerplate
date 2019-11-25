package MPLEndpoint

import (
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"
)

// sampleGet - Accepts a GET request
func MasterPartsListHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	log.Println("Running get MPL!")
	_, err := RefreshMasterPartsList(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("No parts found!")
		return
	}
	// log.Println("We found: ", MPL)

	// json.NewEncoder(w).Encode(MPL)
	// log.Println("Sent JSON response")
}

func FindPartForEditHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running get Part!")
	_, err := FindPartForEdit(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("No parts found!")
		return
	}
	// log.Println("We found: ", part)
}

func SaveMPLEditHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running save MPL!")
	_, err := SaveMPLEdit(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("Save failed!")
		return
	}
	// log.Println("We saved: ", part)
}

func KeywordSearchHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running Keyword Search!")
	_, err := KeywordSearch(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("Search failed!")
		return
	}
	// log.Println("We found: ", part)
}
