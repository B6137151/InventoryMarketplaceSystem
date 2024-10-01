package models

import (
	"github.com/google/uuid"
)

// OrderItem represents an item in an order.
type OrderItem struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID      uuid.UUID `gorm:"type:uuid;not null"`
	VariantID    uuid.UUID `gorm:"type:uuid;not null"` // Adding the VariantID field
	ProductName  string    `gorm:"type:varchar(255);not null"`
	SKUCode      string    `gorm:"type:varchar(100);not null"`
	Quantity     int       `gorm:"not null"`
	PricePerUnit float64   `gorm:"not null"`
	TotalPrice   float64   `gorm:"not null"`
}
