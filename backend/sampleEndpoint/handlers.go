package sampleEndpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// samplePost - Accepts a POST request
func samplePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req sampleRequest
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("You got an error: %v", err)))
		return
	}

	pet, err := getPetById(req.ID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Respond with some the sampleRequest struct converted to JSON
	json.NewEncoder(w).Encode(pet)
}

// sampleGet - Accepts a GET request
func sampleGet(w http.ResponseWriter, r *http.Request) {
	log.Println("Running get pets!")
	pets, err := getPet()
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println("No pets found!")
		return
	}
	log.Println("We found: ", pets)

	json.NewEncoder(w).Encode(pets)
	log.Println("Sent JSON response")
}
