package utils

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func ExecSelectQuery(db *sqlx.DB, query string, fields []interface{}, model interface{}) error {
	if err := db.Select(model, query, fields...); err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	return nil
}

func ExecGetQuery(db *sqlx.DB, query string, fields []interface{}, model interface{}) error {
	if err := db.Get(model, query, fields...); err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	return nil
}

func ExecGetTransactionQuery(db *sqlx.DB, tx *sqlx.Tx, query string, fields []interface{}, model interface{}) error {
	if err := tx.Get(model, query, fields...); err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	return nil
}
