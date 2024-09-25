package functions

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
	database "github.com/nicoletavoinea/GolangProducerConsumer/database/sqlc"
)

var Queries *database.Queries
var Db *sql.DB

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "../database/tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := RunSchema(db, "../database/schema.sql"); err != nil {
		log.Fatalf("Failed to run schema: %v", err)
	}

	Queries = database.New(db)
	return db
}

func RunSchema(db *sql.DB, schemaFilePath string) error {
	schema, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	return err
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}
