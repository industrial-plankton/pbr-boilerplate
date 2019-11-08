package postgres

import (
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnection() *sqlx.DB {
	b, err := ioutil.ReadFile("connection.md")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	s := string(b)
	db, err := sqlx.Connect("postgres", s)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
