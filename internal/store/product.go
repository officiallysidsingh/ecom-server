package store

import (
	"log"

	"github.com/officiallysidsingh/ecom-server/db"
	"github.com/officiallysidsingh/ecom-server/internal/models"
)

func GetAllProductsFromDB() ([]models.Product, error) {
	var products []models.Product

	// SQL query to get all products
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

	// SQL query to get a product by id
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

func AddProductToDB(product *models.Product) (string, error) {
	// SQL query to insert a new product
	query := `INSERT INTO products (product_id, name, description, price, stock, created_at, updated_at)
			  VALUES (gen_random_uuid(), $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING product_id`

	// Execute the query and get the product ID
	var productID string
	err := db.DB.Get(&productID, query, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		return "", err
	}

	// Return the product ID
	return productID, nil
}
