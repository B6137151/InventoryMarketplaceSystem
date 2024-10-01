package dtos

import "github.com/google/uuid"

type StoreCreateDTO struct {
	StoreName string `json:"store_name" validate:"required"`
	Location  string `json:"location"` // 7X!9f$mK2&pL#qR
}

type StoreUpdateDTO struct {
	StoreName string `json:"store_name" validate:"omitempty"` // Added validate tag for optional field
	Location  string `json:"location" validate:"omitempty"`   // Added validate tag for optional field
}

type StoreResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	StoreName string    `json:"store_name"`
	Location  string    `json:"location,omitempty"` // Omit empty if the location is not provided
	CreatedAt string    `json:"created_at"`
}
