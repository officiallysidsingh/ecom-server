package store_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/store"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}

	return mockDB, mock
}

func TestGetAllFromDB(t *testing.T) {
	// Create a new mock database connection
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Wrap mockDB connection with sqlx
	db := sqlx.NewDb(mockDB, "postgres")
	s := store.NewProductStore(db)
	defer db.Close()

	// Create test data
	now := time.Now()
	expectedProducts := []models.Product{
		{
			ProductID:   "prod-1",
			Name:        "Test Product 1",
			Description: "Description 1",
			Price:       99.99,
			Stock:       10,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ProductID:   "prod-2",
			Name:        "Test Product 2",
			Description: "Description 2",
			Price:       149.99,
			Stock:       5,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// Write testcases
	tests := []struct {
		name       string
		mock       func()
		expectErr  bool
		expectRows []models.Product
	}{
		{
			name: "Successful fetch",
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{
						"product_id",
						"name",
						"description",
						"price",
						"stock",
						"created_at",
						"updated_at",
					},
				)

				for _, product := range expectedProducts {
					rows.AddRow(
						product.ProductID,
						product.Name,
						product.Description,
						product.Price,
						product.Stock,
						product.CreatedAt,
						product.UpdatedAt,
					)
				}

				mock.ExpectQuery(regexp.QuoteMeta(`
						SELECT product_id, name, description, price, stock, created_at, updated_at
						FROM products
					`)).
					WillReturnRows(rows)
			},
			expectErr:  false,
			expectRows: expectedProducts,
		},
		{
			name: "No rows found",
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{
						"product_id",
						"name",
						"description",
						"price",
						"stock",
						"created_at",
						"updated_at",
					},
				)

				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT product_id, name, description, price, stock, created_at, updated_at
					FROM products
				`)).
					WillReturnRows(rows)
			},
			expectErr:  false,
			expectRows: []models.Product(nil),
		},
		{
			name: "Query error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT product_id, name, description, price, stock, created_at, updated_at
					FROM products
				`)).
					WillReturnError(errors.New("query error"))
			},
			expectErr:  true,
			expectRows: nil,
		},
	}

	// Run testcases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			products, err := s.GetAllFromDB(context.Background())

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, products)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectRows, products)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByIDFromDB(t *testing.T) {
	// Create a new mock database connection
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Wrap mockDB connection with sqlx
	db := sqlx.NewDb(mockDB, "postgres")
	s := store.NewProductStore(db)
	defer db.Close()

	// Create test data
	now := time.Now()
	expectedProduct := models.Product{
		ProductID:   "prod-1",
		Name:        "Test Product 1",
		Description: "Description 1",
		Price:       99.99,
		Stock:       10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Write testcases
	tests := []struct {
		name      string
		productID string
		mock      func()
		expectErr bool
		expectRow *models.Product
	}{
		{
			name:      "Successful fetch",
			productID: expectedProduct.ProductID,
			mock: func() {
				product := expectedProduct

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
						"name",
						"description",
						"price",
						"stock",
						"created_at",
						"updated_at",
					},
				).AddRow(
					product.ProductID,
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
					product.CreatedAt,
					product.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT product_id, name, description, price, stock, created_at, updated_at
					FROM products
					WHERE product_id = $1
				`)).WithArgs(product.ProductID).WillReturnRows(rows)
			},
			expectErr: false,
			expectRow: &expectedProduct,
		},
		{
			name:      "No rows found",
			productID: "nonexistent-id",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT product_id, name, description, price, stock, created_at, updated_at
					FROM products
					WHERE product_id = $1
				`)).WithArgs("nonexistent-id").WillReturnError(sql.ErrNoRows)
			},
			expectErr: true,
			expectRow: nil,
		},
		{
			name:      "Query error",
			productID: "prod-2",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT product_id, name, description, price, stock, created_at, updated_at
					FROM products
					WHERE product_id = $1
				`)).WithArgs("prod-2").WillReturnError(errors.New("query error"))
			},
			expectErr: true,
			expectRow: nil,
		},
	}

	// Run testcases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			product, err := s.GetByIDFromDB(context.Background(), tt.productID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, product)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectRow, product)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateInDB(t *testing.T) {
	// Create a new mock database connection
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Wrap mockDB connection with sqlx
	db := sqlx.NewDb(mockDB, "postgres")
	s := store.NewProductStore(db)
	defer db.Close()

	// Create test data
	product := models.Product{
		Name:        "Test Product",
		Description: "A test product",
		Price:       100.0,
		Stock:       50,
	}

	// Write testcases
	tests := []struct {
		name      string
		product   *models.Product
		mock      func()
		expectErr bool
		expectID  string
	}{
		{
			name:    "Successful product creation",
			product: &product,
			mock: func() {
				// Mock query for inserting the product
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow("new-product-id")

				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO products (product_id, name, description, price, stock, created_at, updated_at)
					VALUES (gen_random_uuid(), $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
					RETURNING product_id
				`)).WithArgs(
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
				).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			expectErr: false,
			expectID:  "new-product-id",
		},
		{
			name:    "Error starting transaction",
			product: &product,
			mock: func() {
				// Simulate error when starting transaction
				mock.ExpectBegin().WillReturnError(errors.New("failed to start transaction"))
			},
			expectErr: true,
			expectID:  "",
		},
		{
			name:    "Query execution error",
			product: &product,
			mock: func() {
				// Mock the begin and simulate query error
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO products (product_id, name, description, price, stock, created_at, updated_at)
					VALUES (gen_random_uuid(), $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
					RETURNING product_id
				`)).WithArgs(
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
				).WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			expectErr: true,
			expectID:  "",
		},
		{
			name:    "Error committing transaction",
			product: &product,
			mock: func() {
				// Simulate a successful query but an error on commit
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow("new-product-id")

				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO products (product_id, name, description, price, stock, created_at, updated_at)
					VALUES (gen_random_uuid(), $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
					RETURNING product_id
				`)).WithArgs(
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
				).WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectErr: true,
			expectID:  "",
		},
	}

	// Run testcases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			productID, err := s.CreateInDB(context.Background(), tt.product)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectID, productID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectID, productID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPutUpdateInDB(t *testing.T) {
	// Create a new mock database connection
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Wrap mockDB connection with sqlx
	db := sqlx.NewDb(mockDB, "postgres")
	s := store.NewProductStore(db)
	defer db.Close()

	// Create test data
	product := models.Product{
		Name:        "Updated Product",
		Description: "Updated description",
		Price:       150.0,
		Stock:       30,
	}
	productID := "existing-product-id"

	// Write testcases
	tests := []struct {
		name      string
		product   *models.Product
		productID string
		mock      func()
		expectErr bool
	}{
		{
			name:      "Successful put product update",
			product:   &product,
			productID: productID,
			mock: func() {
				// Mock transaction and query for updating product
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow(productID)

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET name=$1, description=$2, price=$3, stock=$4, updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
					productID,
				).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:      "Error starting transaction",
			product:   &product,
			productID: productID,
			mock: func() {
				// Simulate an error when starting the transaction
				mock.ExpectBegin().WillReturnError(errors.New("failed to start transaction"))
			},
			expectErr: true,
		},
		{
			name:      "Product not found (No rows affected)",
			product:   &product,
			productID: "nonexistent-id",
			mock: func() {
				// Simulate an error when no rows are affected by the update
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET name=$1, description=$2, price=$3, stock=$4, updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
					"nonexistent-id",
				).WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Query execution error",
			product:   &product,
			productID: productID,
			mock: func() {
				// Simulate a query error
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET name=$1, description=$2, price=$3, stock=$4, updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
					productID,
				).WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Error committing transaction",
			product:   &product,
			productID: productID,
			mock: func() {
				// Simulate a successful query but an error on commit
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow(productID)

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET name=$1, description=$2, price=$3, stock=$4, updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product.Name,
					product.Description,
					product.Price,
					product.Stock,
					productID,
				).WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectErr: true,
		},
	}

	// Run testcases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := s.PutUpdateInDB(context.Background(), tt.product, tt.productID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPatchUpdateInDB(t *testing.T) {
	// Create a new mock database connection
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Wrap mockDB connection with sqlx
	db := sqlx.NewDb(mockDB, "postgres")
	s := store.NewProductStore(db)
	defer db.Close()

	// Create test data
	productID := "existing-product-id"

	product_all_fields := models.Product{
		Name:        "Updated Product",
		Description: "Updated description",
		Price:       150.0,
		Stock:       30,
	}

	product_missing_fields := models.Product{
		Price: 100.0,
		Stock: 50,
	}

	// Write testcases
	tests := []struct {
		name      string
		product   *models.Product
		productID string
		mock      func()
		expectErr bool
	}{
		{
			name:      "Successful patch update with all fields",
			product:   &product_all_fields,
			productID: productID,
			mock: func() {
				// Mock transaction and query for updating product
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow(productID)

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET
						name = COALESCE(NULLIF($1, ''), name),
						description = COALESCE(NULLIF($2, ''), description),
						price = COALESCE(NULLIF($3, 0), price),
						stock = COALESCE(NULLIF($4, 0), stock),
						updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product_all_fields.Name,
					product_all_fields.Description,
					product_all_fields.Price,
					product_all_fields.Stock,
					productID,
				).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:      "Successful patch update with empty Name and Description fields",
			product:   &product_missing_fields,
			productID: productID,
			mock: func() {
				// Mock transaction and query for updating product
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow(productID)

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET
						name = COALESCE(NULLIF($1, ''), name),
						description = COALESCE(NULLIF($2, ''), description),
						price = COALESCE(NULLIF($3, 0), price),
						stock = COALESCE(NULLIF($4, 0), stock),
						updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					"",
					"",
					product_missing_fields.Price,
					product_missing_fields.Stock,
					productID,
				).WillReturnRows(rows)

				mock.ExpectCommit()
			},

			expectErr: false,
		},
		{
			name:      "Error starting transaction",
			product:   &product_all_fields,
			productID: productID,
			mock: func() {
				// Simulate an error when starting the transaction
				mock.ExpectBegin().WillReturnError(errors.New("failed to start transaction"))
			},
			expectErr: true,
		},
		{
			name:      "Product not found (No rows affected)",
			product:   &product_all_fields,
			productID: "nonexistent-id",
			mock: func() {
				// Simulate an error when no rows are affected by the update
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET
						name = COALESCE(NULLIF($1, ''), name),
						description = COALESCE(NULLIF($2, ''), description),
						price = COALESCE(NULLIF($3, 0), price),
						stock = COALESCE(NULLIF($4, 0), stock),
						updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product_all_fields.Name,
					product_all_fields.Description,
					product_all_fields.Price,
					product_all_fields.Stock,
					"nonexistent-id",
				)

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Query execution error",
			product:   &product_all_fields,
			productID: productID,
			mock: func() {
				// Simulate a query error
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET
						name = COALESCE(NULLIF($1, ''), name),
						description = COALESCE(NULLIF($2, ''), description),
						price = COALESCE(NULLIF($3, 0), price),
						stock = COALESCE(NULLIF($4, 0), stock),
						updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product_all_fields.Name,
					product_all_fields.Description,
					product_all_fields.Price,
					product_all_fields.Stock,
					productID,
				).WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Error committing transaction",
			product:   &product_all_fields,
			productID: productID,
			mock: func() {
				// Simulate a successful query but an error on commit
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow(productID)

				mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE PRODUCTS
					SET
						name = COALESCE(NULLIF($1, ''), name),
						description = COALESCE(NULLIF($2, ''), description),
						price = COALESCE(NULLIF($3, 0), price),
						stock = COALESCE(NULLIF($4, 0), stock),
						updated_at = CURRENT_TIMESTAMP
					WHERE product_id = $5
					RETURNING product_id
				`)).WithArgs(
					product_all_fields.Name,
					product_all_fields.Description,
					product_all_fields.Price,
					product_all_fields.Stock,
					productID,
				).WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectErr: true,
		},
	}

	// Run testcases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := s.PatchUpdateInDB(context.Background(), tt.product, tt.productID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteFromDB(t *testing.T) {
	// Create a new mock database connection
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Wrap mockDB connection with sqlx
	db := sqlx.NewDb(mockDB, "postgres")
	s := store.NewProductStore(db)
	defer db.Close()

	// Create test data
	productID := "existing-product-id"

	// Write testcases
	tests := []struct {
		name      string
		productID string
		mock      func()
		expectErr bool
	}{
		{
			name:      "Successful delete",
			productID: productID,
			mock: func() {
				// Mock transaction and query for updating product
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow(productID)

				mock.ExpectQuery(regexp.QuoteMeta(`
					DELETE FROM products
					WHERE product_id = $1
					RETURNING product_id
				`)).WithArgs(productID).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:      "Error starting transaction",
			productID: productID,
			mock: func() {
				// Simulate an error when starting the transaction
				mock.ExpectBegin().WillReturnError(errors.New("failed to start transaction"))
			},
			expectErr: true,
		},
		{
			name:      "Product not found (No rows affected)",
			productID: "nonexistent-id",
			mock: func() {
				// Simulate an error when no rows are affected by the delete
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
					DELETE FROM products
					WHERE product_id = $1
					RETURNING product_id
				`)).WithArgs("nonexistent-id").WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Query execution error",
			productID: productID,
			mock: func() {
				// Simulate a query error
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`
					DELETE FROM products
					WHERE product_id = $1
					RETURNING product_id
				`)).WithArgs(productID).WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Error committing transaction",
			productID: productID,
			mock: func() {
				// Simulate a successful query but an error on commit
				mock.ExpectBegin()

				rows := sqlmock.NewRows(
					[]string{
						"product_id",
					},
				).AddRow(productID)

				mock.ExpectQuery(regexp.QuoteMeta(`
					DELETE FROM products
					WHERE product_id = $1
					RETURNING product_id
				`)).WithArgs(productID).WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectErr: true,
		},
	}

	// Run testcases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := s.DeleteFromDB(context.Background(), tt.productID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
