package inits

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func NewDB() *sql.DB {
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	pass := os.Getenv("db_pass")
	user := os.Getenv("db_user")
	dbName := os.Getenv("db_name")
	if len(host) == 0 {
		log.Fatal("db host is empty")
	}
	if len(port) == 0 {
		log.Fatal("db port is empty")
	}
	if len(pass) == 0 {
		log.Fatal("db pass is empty")
	}
	if len(user) == 0 {
		log.Fatal("db user is empty")
	}
	if len(pass) == 0 {
		log.Fatal("db pass is empty")
	}
	if len(dbName) == 0 {
		log.Fatal("db name is empty")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("db init error: " + err.Error())
	}
	return db
}
