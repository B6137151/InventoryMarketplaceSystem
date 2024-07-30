package dtos

import (
	"github.com/google/uuid"
	"time"
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
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
