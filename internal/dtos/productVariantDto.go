package dtos

import "github.com/google/uuid"

type ProductVariantCreateDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	SKUCode   string    `json:"sku_code" validate:"required"`
	Price     float64   `json:"price" validate:"required"`
	ImageURL  string    `json:"image_url"`
}

type ProductVariantUpdateDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	SKUCode   string    `json:"sku_code" validate:"required"`
	Price     float64   `json:"price" validate:"required"`
	ImageURL  string    `json:"image_url"`
}

type ProductVariantResponseDTO struct {
	ID        uuid.UUID           `json:"id"`
	ProductID uuid.UUID           `json:"product_id"`
	SKUCode   string              `json:"sku_code"`
	VariantID uuid.UUID           `json:"variant_id"` // เพิ่ม VariantID
	Price     float64             `json:"price"`
	ImageURL  string              `json:"image_url"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
	Product   *ProductResponseDTO `json:"product,omitempty"`
}
type ProductVariantsResponse struct {
	Meta MetaData                    `json:"meta"`
	Data []ProductVariantResponseDTO `json:"data"`
}
