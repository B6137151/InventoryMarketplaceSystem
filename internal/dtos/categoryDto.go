package dtos

import "github.com/google/uuid"

// CategoryCreateDTO is used when creating a new category
type CategoryCreateDTO struct {
	Name    string    `json:"name" validate:"required"`
	StoreID uuid.UUID `json:"store_id" validate:"required"` // Store ID for the category
}

// CategoryUpdateDTO is used when updating an existing category
type CategoryUpdateDTO struct {
	Name    string    `json:"name"`
	StoreID uuid.UUID `json:"store_id" validate:"required"` // Store ID for the category
}

// CategoryResponseDTO is used when returning a category response
type CategoryResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StoreID   uuid.UUID `json:"store_id"` // Store ID for the category
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
