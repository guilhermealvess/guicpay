package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewConnectionDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./guicpay.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Minute * 10)
	db.SetConnMaxLifetime(time.Minute * 10)

	return db
}
