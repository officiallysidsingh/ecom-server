package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID    uuid.UUID `db:"order_id" json:"order_id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	OrderDate  time.Time `db:"order_date" json:"order_date"`
	TotalPrice float64   `db:"total_price" json:"total_price"`
	Status     string    `db:"status" json:"status"`
	Address    string    `db:"address" json:"address"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	OrderItemID uuid.UUID `db:"order_item_id" json:"order_item_id"`
	OrderID     uuid.UUID `db:"order_id" json:"order_id"`
	ProductID   uuid.UUID `db:"product_id" json:"product_id"`
	Quantity    int       `db:"quantity" json:"quantity"`
	UnitPrice   float64   `db:"unit_price" json:"unit_price"`
	TotalPrice  float64   `db:"total_price" json:"total_price"`
}
