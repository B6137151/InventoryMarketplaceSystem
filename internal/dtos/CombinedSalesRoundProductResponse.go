package dtos

import (
	"github.com/google/uuid"
	"time"
)

// CombinedSalesRoundProductResponse is a combined response DTO for sales round and product details
type CombinedSalesRoundProductResponse struct {
	SalesRoundID     uuid.UUID `json:"sales_round_id"`
	SalesRoundName   string    `json:"sales_round_name"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ProductID        uuid.UUID `json:"product_id"`
	StoreID          uuid.UUID `json:"store_id"`
	CategoryID       uuid.UUID `json:"category_id"`
	ProductName      string    `json:"product_name"`
	Brand            string    `json:"brand"`
	Description      string    `json:"description"`
	Currency         string    `json:"currency"`
	Stock            int       `json:"stock"`
	Price            float64   `json:"price"`
	ImageURL         string    `json:"image_url"`
	ProductCreatedAt time.Time `json:"product_created_at"`
	ProductUpdatedAt time.Time `json:"product_updated_at"`
	VariantID        uuid.UUID `json:"variant_id"`
	SKUCode          string    `json:"sku_code"`
	VariantPrice     float64   `json:"variant_price"`
	VariantImageURL  string    `json:"variant_image_url"`
	VariantCreatedAt time.Time `json:"variant_created_at"`
	VariantUpdatedAt time.Time `json:"variant_updated_at"`
}
