package mysql

import (
	"io/ioutil"
	"log"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnection() *sqlx.DB {
	b, err := ioutil.ReadFile("connection.md")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	s := string(b)
	//db, err := sqlx.Connect("mysql", "root:password@(database:3306)/pbr")
	// db, err := sqlx.Connect("mysql", "nickremote:nct.Bcjt@(159.203.63.200:3306)/demodb")
	//db, err := sqlx.Connect("postgres", "cameron:cam@(157.245.171.13:5432)/postgres")
	db, err := sqlx.Connect("postgres", s)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
