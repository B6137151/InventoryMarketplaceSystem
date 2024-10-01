package dtos

import (
	"time"

	"github.com/google/uuid"
)

// OrderItemDTO is used to represent an item in an order
type OrderItemDTO struct {
	ID          uuid.UUID `json:"id,omitempty"`
	VariantID   uuid.UUID `json:"variant_id" validate:"required"`
	ProductName string    `json:"product_name,omitempty"` // Included in the response
	SKUCode     string    `json:"sku_code,omitempty"`     // Included in the response
	Quantity    int       `json:"quantity" validate:"required"`
	Price       float64   `json:"price,omitempty"` // Included in the response
}

// OrderCreateDTO is used when creating a new order
type OrderCreateDTO struct {
	CustomerID      uuid.UUID      `json:"customer_id" validate:"required"`
	RoundID         uuid.UUID      `json:"round_id" validate:"required"`
	Status          string         `json:"status" validate:"required"` // Keep as string for JSON handling
	Code            string         `json:"code" validate:"required"`
	TotalPrice      float64        `json:"total_price" validate:"required"`
	DeliveryAddress string         `json:"delivery_address" validate:"required"`
	PaymentSource   string         `json:"payment_source" validate:"required"` // Keep as string for JSON handling
	Items           []OrderItemDTO `json:"items" validate:"required,dive"`
}

// OrderUpdateDTO is used when updating an existing order
type OrderUpdateDTO struct {
	Status          string  `json:"status" validate:"required"` // Keep as string for JSON handling
	TotalPrice      float64 `json:"total_price" validate:"required"`
	DeliveryAddress string  `json:"delivery_address" validate:"required"`
	PaymentSource   string  `json:"payment_source" validate:"required"` // Keep as string for JSON handling
}

// CustomerDTO is used to represent the customer details in an order response
type CustomerDTO struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// SalesRoundDTO is used to represent the sales round details in an order response
type SalesRoundDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// StoreDTO is used to represent the store details in an order response
type StoreDTO struct {
	ID        uuid.UUID `json:"id"`
	StoreName string    `json:"store_name"`
	Location  string    `json:"location"`
}

// PaymentInfoDTO is used to represent payment details in an order response
type PaymentInfoDTO struct {
	PaymentMethod string `json:"payment_method"`
	PaymentStatus string `json:"payment_status"`
	TransactionID string `json:"transaction_id"`
}

// ShippingInfoDTO is used to represent shipping details in an order response
type ShippingInfoDTO struct {
	ShippingMethod        string    `json:"shipping_method"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
}

// OrderResponseDTO is used when returning an order response
type OrderResponseDTO struct {
	ID              uuid.UUID       `json:"id"`
	OrderDate       time.Time       `json:"order_date"`
	Status          string          `json:"status"`
	Code            string          `json:"code"`
	TotalPrice      float64         `json:"total_price"`
	DeliveryAddress string          `json:"delivery_address"`
	PaymentSource   string          `json:"payment_source"`
	Customer        CustomerDTO     `json:"customer"`
	SalesRound      SalesRoundDTO   `json:"sales_round"`
	Store           StoreDTO        `json:"store"`
	Items           []OrderItemDTO  `json:"items"` // Add items to response DTO
	PaymentInfo     PaymentInfoDTO  `json:"payment_info"`
	ShippingInfo    ShippingInfoDTO `json:"shipping_info"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}
