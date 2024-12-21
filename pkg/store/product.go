package store

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func GetAllProducts() ([]Product, error) {
	var products []Product
	err := db.Select(&products, "SELECT id, name, price FROM products")
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, err
	}
	return products, nil
}
