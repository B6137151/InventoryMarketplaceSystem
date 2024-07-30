package dtos

import (
	"github.com/google/uuid"
	"time"
)

// OrderCreateDTO is used when creating a new order
type OrderCreateDTO struct {
	CustomerID      uuid.UUID `json:"customer_id" validate:"required"`
	RoundID         uuid.UUID `json:"round_id" validate:"required"`
	OrderDate       time.Time `json:"order_date" validate:"required"`
	Status          string    `json:"status" validate:"required"`
	Code            string    `json:"code" validate:"required"`
	TotalPrice      float64   `json:"total_price" validate:"required"`
	DeliveryAddress string    `json:"delivery_address" validate:"required"`
	PaymentSource   string    `json:"payment_source" validate:"required"`
}

// OrderUpdateDTO is used when updating an existing order
type OrderUpdateDTO struct {
	Status          string  `json:"status" validate:"required"`
	TotalPrice      float64 `json:"total_price" validate:"required"`
	DeliveryAddress string  `json:"delivery_address" validate:"required"`
	PaymentSource   string  `json:"payment_source" validate:"required"`
}

// OrderResponseDTO is used when returning an order response
type OrderResponseDTO struct {
	ID              uuid.UUID `json:"id"`
	CustomerID      uuid.UUID `json:"customer_id"`
	RoundID         uuid.UUID `json:"round_id"`
	OrderDate       time.Time `json:"order_date"`
	Status          string    `json:"status"`
	Code            string    `json:"code"`
	TotalPrice      float64   `json:"total_price"`
	DeliveryAddress string    `json:"delivery_address"`
	PaymentSource   string    `json:"payment_source"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}
