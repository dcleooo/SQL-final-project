package src 

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func dbConn() (db *sql.DB) {
    db, err := sql.Open("sqlite3", "SQL/data.db")
    if err != nil {
        log.Fatal(err)
    }
    return db
}