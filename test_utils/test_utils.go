package test_utils

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "modernc.org/sqlite"
)

func DBSetup() *sql.DB {
	// Memory-only database only lasts for duration of 'db' variable
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Database connected")

	query, err := ioutil.ReadFile("../sql/schemas.sql")
	if err != nil {
		log.Println("HERE")
		log.Fatalln(err)
	}

	if _, err := db.Exec(string(query)); err != nil {
		log.Println("HERE")
		log.Fatalln(err)
	}

	return db
}
