package dtos

import "github.com/google/uuid"

type StoreCreateDTO struct {
	StoreName string `json:"store_name" validate:"required"`
	Location  string `json:"location"`
}

type StoreUpdateDTO struct {
	StoreName string `json:"store_name"`
	Location  string `json:"location"`
}

type StoreResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	StoreName string    `json:"store_name"`
	Location  string    `json:"location"`
	CreatedAt string    `json:"created_at"`
}
