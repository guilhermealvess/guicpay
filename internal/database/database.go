package database

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // driver postgres
)

func NewConnectionDB(conn string) *sqlx.DB {
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		log.Panicf("failed to connect on database: %v", err)
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
