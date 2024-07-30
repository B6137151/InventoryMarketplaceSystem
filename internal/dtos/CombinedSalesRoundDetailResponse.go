package dtos

import (
	"github.com/google/uuid"
	"time"
)

// CombinedSalesRoundDetailResponse is a combined response DTO for sales round details
type CombinedSalesRoundDetailResponse struct {
	SalesRound
	SalesRoundDetailResponseDTO
	SKUCode         string  `json:"sku_code"`
	VariantPrice    float64 `json:"variant_price"`
	VariantImageURL string  `json:"variant_image_url"`
	ProductName     string  `json:"product_name"`
	Brand           string  `json:"brand"`
	Description     string  `json:"description"`
	Currency        string  `json:"currency"`
	Stock           int     `json:"stock"`
	ProductPrice    float64 `json:"product_price"`
}

// SalesRound represents a sales round in the system
type SalesRound struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
