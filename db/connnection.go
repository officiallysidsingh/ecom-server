package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Init the DB connection.
func InitDB(connStr string) (*sqlx.DB, error) {
	// Open a connection to the database
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Check if the database is reachable
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %w", err)
	}

	log.Println("DB Connected Successfully!!")
	return db, nil
}

// Close the DB connection
func CloseDB(db *sqlx.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		log.Println("DB Connection Closed Successfully!!")
	}
}
