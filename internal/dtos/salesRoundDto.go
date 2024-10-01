package dtos

import (
	"time"

	"github.com/google/uuid"
)

// SalesRoundCreateDTO is used for creating a new sales round
type SalesRoundCreateDTO struct {
	Name      string    `json:"name" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

// SalesRoundUpdateDTO is used for updating an existing sales round
type SalesRoundUpdateDTO struct {
	Name      string    `json:"name" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

// SalesRoundResponseDTO is used for returning a sales round response
type SalesRoundResponseDTO struct {
	ID                uuid.UUID                     `json:"id"`
	Name              string                        `json:"name"`
	StartDate         time.Time                     `json:"start_date"`
	EndDate           time.Time                     `json:"end_date"`
	CreatedAt         time.Time                     `json:"created_at"`
	UpdatedAt         time.Time                     `json:"updated_at"`
	Products          []ProductResponseDTO          `json:"products,omitempty"`
	ProductVariants   []ProductVariantResponseDTO   `json:"product_variants,omitempty"`
	SalesRoundDetails []SalesRoundDetailResponseDTO `json:"sales_round_details,omitempty"`
	Categories        []CategoryResponseDTO         `json:"categories,omitempty"`
	Stores            []StoreResponseDTO            `json:"stores,omitempty"`
}

// ProductVariantResponseDTO represents the product variant information
// type ProductVariantResponseDTO struct {
// 	ID        uuid.UUID `json:"id"`
// 	ProductID uuid.UUID `json:"product_id"`
// 	SKUCode   string    `json:"sku_code"`
// 	Price     float64   `json:"price"`
// 	ImageURL  string    `json:"image_url"`
// }

// SalesRoundDetailResponseDTO represents the sales round detail information
// type SalesRoundDetailResponseDTO struct {
// 	ID            uuid.UUID `json:"id"`
// 	RoundID       uuid.UUID `json:"round_id"`
// 	VariantID     uuid.UUID `json:"variant_id"`
// 	Quantity      int       `json:"quantity"`
// 	Remaining     int       `json:"remaining"`
// 	ProductStock  int       `json:"product_stock"`
// 	QuantityLimit int       `json:"quantity_limit"`
// }

// CategoryResponseDTO represents the category information
// type CategoryResponseDTO struct {
// 	ID   uuid.UUID `json:"id"`
// 	Name string    `json:"name"`
// }
