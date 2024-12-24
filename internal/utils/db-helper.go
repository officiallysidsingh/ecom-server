package utils

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/db"
)

func ExecSelectQuery(query string, fields []interface{}, model interface{}) error {
	if err := db.DB.Select(model, query, fields...); err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	return nil
}

func ExecGetQuery(query string, fields []interface{}, model interface{}) error {
	if err := db.DB.Get(model, query, fields...); err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	return nil
}

func ExecGetTransactionQuery(tx *sqlx.Tx, query string, fields []interface{}, model interface{}) error {
	if err := tx.Get(model, query, fields...); err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	return nil
}
