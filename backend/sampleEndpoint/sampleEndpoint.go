package sampleEndpoint

import (
	"errors"
	"log"
)

type sampleRequest struct {
	ID string `json:"ID"`
}

type Pet struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Owner   string `json:"owner" db:"owner"`
	Species string `json:"species" db:"species"`
	Sex     string `json:"sex" db:"sex"`
}

func getPets() (Pet, error) {
	pet := Pet{}

	// this will pull the first pet directly into the pet variable
	err := db.Get(&pet, "SELECT * FROM pet LIMIT 1")

	if err != nil {
		log.Println("Error: ", err)
		return pet, err
	}

	if pet.Name == "" {
		return pet, errors.New("No pet found")
	}

	return pet, nil
}
