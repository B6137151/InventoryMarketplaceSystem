package dtos

import "github.com/google/uuid"

type CheckoutRequestDTO struct {
	CartID          uuid.UUID `json:"CartID" validate:"required"`
	PaymentSource   string    `json:"PaymentSource" validate:"required"`
	DeliveryAddress string    `json:"DeliveryAddress" validate:"required"`
	Currency        string    `json:"Currency" validate:"required"`
}
