package models

import (
	"time"
)

type Order struct {
	OrderID       string    `db:"order_id" json:"order_id"`
	UserID        string    `db:"user_id" json:"user_id"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	TaxPrice      float64   `db:"tax_price" json:"tax_price"`
	ShippingPrice float64   `db:"shipping_price" json:"shipping_price"`
	TotalPrice    float64   `db:"total_price" json:"total_price"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	Items         []OrderItem
}

type OrderItem struct {
	OrderItemID string  `db:"order_item_id" json:"order_item_id"`
	OrderID     string  `db:"order_id" json:"order_id"`
	ProductID   string  `db:"product_id" json:"product_id"`
	Quantity    int     `db:"quantity" json:"quantity"`
	UnitPrice   float64 `db:"unit_price" json:"unit_price"`
	TotalPrice  float64 `db:"total_price" json:"total_price"`
}
