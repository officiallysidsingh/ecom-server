package store

import (
	"log"

	"github.com/officiallysidsingh/ecom-server/db"
	"github.com/officiallysidsingh/ecom-server/internal/models"
)

func GetAllProductsFromDB() ([]models.Product, error) {
	var products []models.Product

	query := `SELECT product_id, name, description, price, stock, created_at, updated_at
			  FROM products`

	err := db.DB.Select(&products, query)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, err
	}
	return products, nil
}

func GetProductByIdFromDB(productID string) (*models.Product, error) {
	var product models.Product

	query := `SELECT product_id, name, description, price, stock, created_at, updated_at
              FROM products
              WHERE product_id = $1`

	err := db.DB.Get(&product, query, productID)
	if err != nil {
		log.Printf("Error fetching product with ID %s: %v", productID, err)
		return nil, err
	}

	return &product, nil
}
