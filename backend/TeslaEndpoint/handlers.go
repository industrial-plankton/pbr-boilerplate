package TeslaEndpoint

import (
	// "encoding/json"
	// "fmt"
	"backend/utility"
	"log"
	"net/http"
)

// sampleGet - Accepts a GET request
// func MasterPartsListHandle(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
// 	log.Println("Running get MPL!")
// 	_, err := RefreshMasterPartsList(r.Header)
// 	if err != nil {
// 		// w.WriteHeader(http.StatusNoContent)
// 		w.Write([]byte(err.Error()))
// 		utility.Log(err)
// 		log.Println("No parts found!")
// 		return
// 	}
// 	// log.Println("We found: ", MPL)

// 	// json.NewEncoder(w).Encode(MPL)
// 	// log.Println("Sent JSON response")

// }

func FindShipForEditHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running get Ship!")
	_, err := FindShipForEdit(r.Header)
	if err != nil {
		// w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(err.Error()))
		utility.Log(err)
		log.Println("No parts found!")
		return
	}
	// log.Println("We found: ", part)

}

func SaveTeslaEditHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Running save Tesla!")
	_, err := SaveTeslaEdit(r.Header)
	if err != nil {
		// w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(err.Error()))
		utility.Log(err)
		log.Println("Save failed!")
		return
	}
	// log.Println("We saved: ", part)
}

// func KeywordSearchHandle(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Running Keyword Search!")
// 	_, err := KeywordSearch(r.Header)
// 	if err != nil {
// 		// w.WriteHeader(http.StatusNoContent)
// 		w.Write([]byte(err.Error()))
// 		utility.Log(err)
// 		log.Println("Search failed!")
// 		return
// 	}
// 	// log.Println("We found: ", part)
// }
