package dtos

import "github.com/google/uuid"

// OrderDetailCreateDTO is the structure used for creating a new OrderDetail
type OrderDetailCreateDTO struct {
	OrderID    uuid.UUID `json:"order_id" validate:"required"`
	VariantID  uuid.UUID `json:"variant_id" validate:"required"`
	Quantity   int       `json:"quantity" validate:"required"`
	Price      float64   `json:"price" validate:"required"`
	TotalPrice float64   `json:"total_price" validate:"required"`
}

// OrderDetailUpdateDTO is the structure used for updating an existing OrderDetail
type OrderDetailUpdateDTO struct {
	Quantity   int     `json:"quantity" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	TotalPrice float64 `json:"total_price" validate:"required"`
}

// OrderDetailResponseDTO is the structure used for responding with OrderDetail data
type OrderDetailResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	VariantID  uuid.UUID `json:"variant_id"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}
