package dtos

import "github.com/google/uuid"

type CartItemDTO struct {
	VariantID uuid.UUID `json:"variant_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
}

type CartRequestDTO struct {
	CustomerID    uuid.UUID     `json:"customer_id" validate:"required"`
	RoundID       uuid.UUID     `json:"round_id" validate:"required"`
	StoreID       uuid.UUID     `json:"store_id" validate:"required"`
	Items         []CartItemDTO `json:"items" validate:"required"`
	PaymentSource string        `json:"payment_source" validate:"required"`
}

type CartInteractionResponseDTO struct {
	Message       string        `json:"message,omitempty"`
	CartID        uuid.UUID     `json:"cart_id"`
	CustomerID    uuid.UUID     `json:"customer_id"`
	CustomerName  string        `json:"customer_name"`
	RoundID       uuid.UUID     `json:"round_id"`
	StoreID       uuid.UUID     `json:"store_id"`
	StoreName     string        `json:"store_name"`
	PaymentSource string        `json:"payment_source"`
	Items         []CartItemDTO `json:"items"`
	TotalAmount   float64       `json:"total_amount"`
}

type CartItemDetailDTO struct {
	VariantID   uuid.UUID `json:"variant_id"`
	Quantity    int       `json:"quantity"`
	ProductName string    `json:"product_name"`
	SKUCode     string    `json:"sku_code"`
	Price       float64   `json:"price"`
	TotalPrice  float64   `json:"total_price"`
	ImageURL    string    `json:"image_url"`
}

type CartDetailResponseDTO struct {
	Message       string              `json:"message"`
	CartID        uuid.UUID           `json:"cart_id"`
	CustomerID    uuid.UUID           `json:"customer_id"`
	CustomerName  string              `json:"customer_name"`
	RoundID       uuid.UUID           `json:"round_id"`
	StoreID       uuid.UUID           `json:"store_id"`
	StoreName     string              `json:"store_name"`
	PaymentSource string              `json:"payment_source"`
	Items         []CartItemDetailDTO `json:"items"`
	TotalAmount   float64             `json:"total_amount"`
	MetaData      MetaData            `json:"meta_data"`
}
