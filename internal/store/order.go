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

type OrderStore interface {
	GetAllFromDB(ctx context.Context, userID string) ([]models.Order, error)
	GetByIDFromDB(ctx context.Context, orderID string, userID string) (*models.Order, error)
	// PutUpdateInDB(ctx context.Context, order *models.Order, orderID string) error
	// PatchUpdateInDB(ctx context.Context, order *models.Order, orderID string) error
	// CreateInDB(ctx context.Context, order *models.Order) (string, error)
	// DeleteFromDB(ctx context.Context, orderID string) error
}

type orderStore struct {
	db *sqlx.DB
}

func NewOrderStore(db *sqlx.DB) OrderStore {
	return &orderStore{
		db: db,
	}
}

func (s *orderStore) GetAllFromDB(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order

	// SQL query to get all orders
	query := `
		SELECT order_id, user_id, payment_method, tax_price, shipping_price, total_price, created_at, updated_at
		FROM orders
		WHERE user_id = $1
	`

	fields := []interface{}{
		userID,
	}

	if err := utils.ExecSelectQuery(
		s.db,
		query,
		fields,
		&orders,
	); err != nil {
		log.Printf("Error fetching orders from DB: %v", err)
		return nil, err
	}

	return orders, nil
}

func (s *orderStore) GetByIDFromDB(ctx context.Context, orderID string, userID string) (*models.Order, error) {
	var order models.Order

	// SQL query to get an order by id
	query := `
		SELECT order_id, user_id, payment_method, tax_price, shipping_price, total_price, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		AND order_id = $2
	`

	fields := []interface{}{
		userID,
		orderID,
	}

	if err := utils.ExecSelectQuery(
		s.db,
		query,
		fields,
		&order,
	); err != nil {
		// If no rows found
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Order with userID %s and orderID %s not found", userID, orderID)
			return nil, fmt.Errorf("order with userID %s and orderID %s not found", userID, orderID)
		}
		log.Printf("Error fetching order with userID %s and orderID %s from DB: %v", userID, orderID, err)
		return nil, err
	}

	return &order, nil
}
