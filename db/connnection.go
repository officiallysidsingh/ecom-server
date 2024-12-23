package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Global variable to store the database connection
var DB *sqlx.DB

// Init the DB connection.
func InitDB(connStr string) error {
	var err error

	// Open a connection to the database
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Check if the database is reachable
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("unable to ping the database: %w", err)
	}

	log.Println("DB Connected Successfully!!")
	return nil
}

// Close the DB connection
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		log.Println("DB Connection Closed Successfully!!")
	}
}
