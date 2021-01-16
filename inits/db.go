package inits

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "mef:derparol@tcp(127.0.0.1:3306)/otus")
	if err != nil {
		log.Fatal("db init error: " + err.Error())
	}
	return db
}
