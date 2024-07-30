package dtos

import "github.com/google/uuid"

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
	ID            uuid.UUID `json:"id"`
	RoundID       uuid.UUID `json:"round_id"`
	VariantID     uuid.UUID `json:"variant_id"`
	Quantity      int       `json:"quantity"`
	Remaining     int       `json:"remaining"`
	QuantityLimit int       `json:"quantity_limit"`
	ProductStock  int       `json:"product_stock"` // New field for product stock
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}
