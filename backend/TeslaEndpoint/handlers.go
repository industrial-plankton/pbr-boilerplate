package TeslaEndpoint

import (
	// "encoding/json"
	// "fmt"
	"backend/utility"
	"log"
	"net/http"
)

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
