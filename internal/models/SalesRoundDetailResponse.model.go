package models

import "github.com/google/uuid"

type SalesRoundDetailResponse struct {
	VariantID       uuid.UUID `json:"variant_id"`
	ProductID       uuid.UUID `json:"product_id"`
	SKUCode         string    `json:"sku_code"`
	VariantPrice    float64   `json:"variant_price"`
	VariantImageURL string    `json:"variant_image_url"`
	ProductName     string    `json:"product_name"`
	Brand           string    `json:"brand"`
	Description     string    `json:"description"`
	Currency        string    `json:"currency"`
	Stock           int       `json:"stock"`
	ProductPrice    float64   `json:"product_price"`
}
