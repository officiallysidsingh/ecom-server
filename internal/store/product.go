package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/officiallysidsingh/ecom-server/db"
	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

func GetAllProductsFromDB() ([]models.Product, error) {
	var products []models.Product

	// SQL query to get all products
	query := `
		SELECT product_id, name, description, price, stock, created_at, updated_at
		FROM products
	`

	if err := utils.ExecSelectQuery(query, nil, &products); err != nil {
		log.Printf("Error fetching products from DB: %v", err)
		return nil, err
	}

	return products, nil
}

func GetProductByIdFromDB(productID string) (*models.Product, error) {
	var product models.Product

	// SQL query to get a product by id
	query := `
		SELECT product_id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE product_id = $1
	`

	fields := []interface{}{
		productID,
	}

	if err := utils.ExecGetQuery(query, fields, &product); err != nil {
		// If no rows found
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Product with ID %s not found", productID)
			return nil, fmt.Errorf("product with ID %s not found", productID)
		}
		log.Printf("Error fetching product with ID %s from DB: %v", productID, err)
		return nil, err
	}

	return &product, nil
}

func AddProductToDB(product *models.Product) (string, error) {
	// Begin a transaction to ensure atomicity
	tx, err := db.DB.Beginx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return "", fmt.Errorf("failed to start database transaction: %w", err)
	}

	// Ensure transaction is properly rolled back in case of failure
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
		}
	}()

	// SQL query to insert a new product
	query := `
		INSERT INTO products (product_id, name, description, price, stock, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING product_id
	`

	fields := []interface{}{
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	}

	// Execute the query and get the product ID
	var productID string
	if err = utils.ExecGetTransactionQuery(
		tx,
		query,
		fields,
		&productID,
	); err != nil {
		// If no rows found
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Product with ID %s not found", productID)
			return "", fmt.Errorf("product with ID %s not found", productID)
		}
		log.Printf("Error adding product with ID %s to DB: %v", productID, err)
		return "", err
	}

	// Commit the transaction if update was successful
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction for product with ID %s: %v", productID, err)
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log the success and return the updated product ID
	log.Printf("Product with ID %s added successfully", productID)
	return productID, nil
}

func PutUpdateProductInDB(product *models.Product, productID string) error {
	// Begin a transaction to ensure atomicity
	tx, err := db.DB.Beginx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return fmt.Errorf("failed to start database transaction: %w", err)
	}

	// Ensure transaction is properly rolled back in case of failure
	defer func() {
		if err != nil {
			// If error occurs, we rollback the transaction
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
		}
	}()

	// SQL query to update a product
	query := `
		UPDATE PRODUCTS
		SET name=$1, description=$2, price=$3, stock=$4, updated_at = CURRENT_TIMESTAMP
		WHERE product_id = $5
		RETURNING product_id
	`

	fields := []interface{}{
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		productID,
	}

	// Variable to capture the updated product ID
	var updatedProductID string
	if err := utils.ExecGetTransactionQuery(
		tx,
		query,
		fields,
		&updatedProductID,
	); err != nil {
		// If no rows affected (Product Not Found)
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Product with ID %s not found", productID)
			return fmt.Errorf("product with ID %s not found", productID)
		}
		// General error
		log.Printf("Error updating product in DB: %v", err)
		return fmt.Errorf("failed to update product with ID %s: %w", productID, err)
	}

	// Commit the transaction if update was successful
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction for product with ID %s: %v", productID, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log the success and return the updated product ID
	log.Printf("Product with ID %s updated successfully", updatedProductID)
	return nil
}
