package DataValidationEndpoint

import (
	// "encoding/json"
	// "fmt"
	"backend/utility"
	"log"
	"net/http"
)

//RefreshHandler refreshes the MPL on googel sheets
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Running get Ship!")
	_, err := Refresh(r.Header)
	if err != nil {
		// w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(err.Error()))
		utility.Log(err)
		log.Println("No parts found!")
		return
	}
	// log.Println("We found: ", part)
}
