package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseCreateDTO struct {
	CustomerID      uuid.UUID         `json:"customer_id"`
	RoundID         uuid.UUID         `json:"round_id"`
	OrderDate       time.Time         `json:"order_date"`
	Code            string            `json:"code"`
	TotalPrice      float64           `json:"total_price"`
	DeliveryAddress string            `json:"delivery_address"`
	PaymentSource   string            `json:"payment_source"`
	Items           []PurchaseItemDTO `json:"items"`
}

type PurchaseItemDTO struct {
	VariantID uuid.UUID `json:"variant_id"`
	Quantity  int       `json:"quantity"`
}

type PurchaseResponseDTO struct {
	ID              uuid.UUID         `json:"id"`
	CustomerID      uuid.UUID         `json:"customer_id"`
	RoundID         uuid.UUID         `json:"round_id"`
	OrderDate       time.Time         `json:"order_date"`
	Status          string            `json:"status"`
	Code            string            `json:"code"`
	TotalPrice      float64           `json:"total_price"`
	DeliveryAddress string            `json:"delivery_address"`
	PaymentSource   string            `json:"payment_source"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
	Items           []PurchaseItemDTO `json:"items"` // Add items to response DTO

}
