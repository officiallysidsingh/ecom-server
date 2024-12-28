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
