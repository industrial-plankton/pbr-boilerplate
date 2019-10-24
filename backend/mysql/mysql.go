package mysql

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewConnection() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root:password@(database:3306)/pbr")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
