package database

import (
	"log"
	"time"

	"github.com/guilhermealvess/guicpay/internal/properties"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // driver postgres
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

func NewConnectionDB() *sqlx.DB {
	db, err := otelsqlx.Open("postgres", properties.Props.DatabaseURL)
	if err != nil {
		log.Panicf("failed to connect on database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(properties.Props.DatabaseMaxConn)
	db.SetMaxIdleConns(properties.Props.DatabaseMaxIdle)
	db.SetConnMaxIdleTime(time.Minute * 10)
	db.SetConnMaxLifetime(time.Minute * 10)

	return db
}
