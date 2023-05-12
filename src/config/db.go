package config

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	devPort     = 5432
	testPort = 5433
	username = "irfan"
	password = "123456"
	devDB    = "bookstore"
	testDB   = "bookstore_test"
)

func ConnectDB() *sql.DB {
	env := os.Getenv("APP_ENV")

	// Set the appropriate database connection string based on the environment
	var dbname string
	var port int
	switch env {
	case "development":
		dbname = devDB
		port = devPort
	case "test":
		dbname = testDB
		port = testPort
	default:
		log.Fatal("Invalid or missing APP_ENV environment variable")
	}
	psqlInfo := "host=" + host + " port=" + strconv.Itoa(port) + " user=" + username + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
