package config

import (
	"database/sql"
	"log"
	"strconv"
)

const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "password"
	dbname   = "bookstore"
)

func ConnectDB() *sql.DB {
	psqlInfo := "host=" + host + " port=" + strconv.Itoa(port) + " user=" + username + " password=" + password + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to database!")

	return db
}
