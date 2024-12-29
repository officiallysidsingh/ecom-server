package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

type ProductStore interface {
	GetAllFromDB(ctx context.Context) ([]models.Product, error)
	GetByIDFromDB(ctx context.Context, id string) (*models.Product, error)
	CreateInDB(ctx context.Context, product *models.Product) (string, error)
	PutUpdateInDB(ctx context.Context, product *models.Product, id string) error
	PatchUpdateInDB(ctx context.Context, product *models.Product, id string) error
	DeleteFromDB(ctx context.Context, id string) error
}

type productStore struct {
	db *sqlx.DB
}

func NewProductStore(db *sqlx.DB) ProductStore {
	return &productStore{
		db: db,
	}
}

func (s *productStore) GetAllFromDB(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	// SQL query to get all products
	query := `
		SELECT product_id, name, description, price, stock, created_at, updated_at
		FROM products
	`

	if err := utils.ExecSelectQuery(
		s.db,
		query,
		nil,
		&products,
	); err != nil {
		log.Printf("Error fetching products from DB: %v", err)
		return nil, err
	}

	return products, nil
}

func (s *productStore) GetByIDFromDB(ctx context.Context, productID string) (*models.Product, error) {
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

	if err := utils.ExecGetQuery(
		s.db,
		query,
		fields,
		&product,
	); err != nil {
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

func (s *productStore) CreateInDB(ctx context.Context, product *models.Product) (string, error) {
	// Begin a transaction to ensure atomicity
	tx, err := s.db.Beginx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return "", fmt.Errorf("failed to start database transaction: %w", err)
	}

	// For error handling in the deferred rollback
	var txErr error

	// Ensure transaction is properly rolled back in case of failure
	defer func() {
		if txErr != nil {
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

	// Execute the query and return the added product ID
	var productID string
	txErr = utils.ExecGetTransactionQuery(
		s.db,
		tx,
		query,
		fields,
		&productID,
	)
	if txErr != nil {
		log.Printf("Error adding product with Name %s to DB: %v", product.Name, txErr)
		return "", txErr
	}

	// Commit the transaction if update was successful
	txErr = tx.Commit()
	if txErr != nil {
		log.Printf("Error committing transaction for product with ID %s: %v", productID, txErr)
		return "", fmt.Errorf("failed to commit transaction: %w", txErr)
	}

	// Log the success and return the updated product ID
	log.Printf("Product with ID %s added successfully", productID)
	return productID, nil
}

func (s *productStore) PutUpdateInDB(ctx context.Context, product *models.Product, productID string) error {
	// Begin a transaction to ensure atomicity
	tx, err := s.db.Beginx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return fmt.Errorf("failed to start database transaction: %w", err)
	}

	// For error handling in the deferred rollback
	var txErr error

	// Ensure transaction is properly rolled back in case of failure
	defer func() {
		if txErr != nil {
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

	// Execute the query and return the updated product ID
	var updatedProductID string
	txErr = utils.ExecGetTransactionQuery(
		s.db,
		tx,
		query,
		fields,
		&updatedProductID,
	)
	if txErr != nil {
		// If no rows affected (Product Not Found)
		if errors.Is(txErr, sql.ErrNoRows) {
			log.Printf("Product with ID %s not found", productID)
			return fmt.Errorf("product with ID %s not found", productID)
		}
		// General error
		log.Printf("Error updating product in DB: %v", txErr)
		return fmt.Errorf("failed to update product with ID %s: %w", productID, txErr)
	}

	// Commit the transaction if update was successful
	txErr = tx.Commit()
	if txErr != nil {
		log.Printf("Error committing transaction for product with ID %s: %v", productID, txErr)
		return fmt.Errorf("failed to commit transaction: %w", txErr)
	}

	// Log the success and return the updated product ID
	log.Printf("Product with ID %s updated successfully", updatedProductID)
	return nil
}

func (s *productStore) PatchUpdateInDB(ctx context.Context, product *models.Product, productID string) error {
	// Begin a transaction to ensure atomicity
	tx, err := s.db.Beginx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return fmt.Errorf("failed to start database transaction: %w", err)
	}

	// For error handling in the deferred rollback
	var txErr error

	// Ensure transaction is properly rolled back in case of failure
	defer func() {
		if txErr != nil {
			// If error occurs, we rollback the transaction
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
		}
	}()

	// SQL query to update a product using COALESCE
	query := `
		UPDATE PRODUCTS
		SET
			name = COALESCE(NULLIF($1, ''), name),
			description = COALESCE(NULLIF($2, ''), description),
			price = COALESCE(NULLIF($3, 0), price),
			stock = COALESCE(NULLIF($4, 0), stock),
			updated_at = CURRENT_TIMESTAMP
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

	// Execute the query and return the updated product ID
	var updatedProductID string
	txErr = utils.ExecGetTransactionQuery(
		s.db,
		tx,
		query,
		fields,
		&updatedProductID,
	)

	// If no rows affected (Product Not Found)
	if errors.Is(txErr, sql.ErrNoRows) {
		log.Printf("Product with ID %s not found", productID)
		return fmt.Errorf("product with ID %s not found", productID)
	}

	// General error handling if the query failed
	if txErr != nil {
		log.Printf("Error updating product in DB: %v", txErr)
		return fmt.Errorf("failed to update product with ID %s: %w", productID, txErr)
	}

	// Commit the transaction if update was successful
	txErr = tx.Commit()
	if txErr != nil {
		log.Printf("Error committing transaction for product with ID %s: %v", productID, txErr)
		return fmt.Errorf("failed to commit transaction: %w", txErr)
	}

	// Log the success and return the updated product ID
	log.Printf("Product with ID %s updated successfully", updatedProductID)
	return nil
}

func (s *productStore) DeleteFromDB(ctx context.Context, productID string) error {
	// Begin a transaction to ensure atomicity
	tx, err := s.db.Beginx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return fmt.Errorf("failed to start database transaction: %w", err)
	}

	// For error handling in the deferred rollback
	var txErr error

	// Ensure transaction is properly rolled back in case of failure
	defer func() {
		if txErr != nil {
			// If error occurs, we rollback the transaction
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
		}
	}()

	// SQL query to delete a product
	query := `
		DELETE FROM products 
		WHERE product_id = $1
		RETURNING product_id
	`
	fields := []interface{}{
		productID,
	}

	// Execute the query and return the deleted product ID
	var deletedProductID string
	txErr = utils.ExecGetTransactionQuery(
		s.db,
		tx,
		query,
		fields,
		&deletedProductID,
	)

	if txErr != nil {
		// If no rows affected (Product Not Found)
		if errors.Is(txErr, sql.ErrNoRows) {
			log.Printf("Product with ID %s not found", productID)
			return fmt.Errorf("product with ID %s not found", productID)
		}
		// General error
		log.Printf("Error deleting product in DB: %v", txErr)
		return fmt.Errorf("failed to delete product with ID %s: %w", productID, txErr)
	}

	// Commit the transaction if delete was successful
	txErr = tx.Commit()
	if txErr != nil {
		log.Printf("Error committing transaction for product with ID %s: %v", productID, txErr)
		return fmt.Errorf("failed to commit transaction: %w", txErr)
	}

	// Log the success and return the deleted product ID
	log.Printf("Product with ID %s deleted successfully", deletedProductID)
	return nil
}
