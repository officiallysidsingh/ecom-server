package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// Global variable to store the database connection
var DB *sqlx.DB

// InitDB initializes the database connection.
func InitDB(connStr string) error {
	// Open a connection to the database
	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Check if the database is reachable
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("unable to ping the database: %w", err)
	}

	log.Println("Successfully connected to the database")

	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}
