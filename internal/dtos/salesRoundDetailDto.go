package dtos

import (
	"time"

	"github.com/google/uuid"
)

// SalesRoundDetailCreateDTO is the structure used for creating a new SalesRoundDetail
type SalesRoundDetailCreateDTO struct {
	RoundID       uuid.UUID `json:"round_id" validate:"required"`
	VariantID     uuid.UUID `json:"variant_id" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required"`
	QuantityLimit int       `json:"quantity_limit" validate:"required"`
	ProductStock  int       `json:"product_stock" validate:"required"` // New field for product stock
}

// SalesRoundDetailUpdateDTO is the structure used for updating an existing SalesRoundDetail
type SalesRoundDetailUpdateDTO struct {
	RoundID       uuid.UUID `json:"round_id" validate:"required"`
	VariantID     uuid.UUID `json:"variant_id" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required"`
	QuantityLimit int       `json:"quantity_limit" validate:"required"`
	Remaining     int       `json:"remaining" validate:"required"`
	ProductStock  int       `json:"product_stock" validate:"required"` // New field for product stock
}

// SalesRoundDetailResponseDTO is the structure used for responding with SalesRoundDetail data
type SalesRoundDetailResponseDTO struct {
	ID             uuid.UUID                  `json:"id"`
	RoundID        uuid.UUID                  `json:"round_id"`
	VariantID      uuid.UUID                  `json:"variant_id"`
	Quantity       int                        `json:"quantity"`
	Remaining      int                        `json:"remaining"`
	QuantityLimit  int                        `json:"quantity_limit"`
	ProductStock   int                        `json:"product_stock"`
	SalesRound     SalesRoundInfoResponse     `json:"sales_round"`     // Nested struct for SalesRound
	ProductVariant ProductVariantInfoResponse `json:"product_variant"` // Nested struct for ProductVariant
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

// SalesRoundInfoResponse is a nested structure for SalesRound details
type SalesRoundInfoResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductVariantInfoResponse is a nested structure for ProductVariant details
type ProductVariantInfoResponse struct {
	VariantID uuid.UUID `json:"variant_id"` // Add VariantID if needed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	SKUCode   string    `json:"sku_code"`
	Price     float64   `json:"price"`
	ImageURL  string    `json:"image_url"`
}
