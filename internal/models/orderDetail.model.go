package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type OrderDetail struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PurchaseID uuid.UUID      `gorm:"type:uuid;not null;index"` // Added this line
	CreatedAt  time.Time      `gorm:"type:timestamp with time zone"`
	UpdatedAt  time.Time      `gorm:"type:timestamp with time zone"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
	OrderID    uuid.UUID      `gorm:"type:uuid;not null;index"` // Foreign key for the Order
	VariantID  uuid.UUID      `gorm:"type:uuid;not null;index"` // Foreign key for the ProductVariant
	Quantity   int            `gorm:"not null"`                 // Quantity of the product variant ordered
	Price      float64        `gorm:"not null"`                 // Price per unit of the product variant at the time of order
	TotalPrice float64        `gorm:"not null"`                 // Total price for the quantity ordered

	Order          Order          `gorm:"foreignKey:OrderID;references:ID"`
	ProductVariant ProductVariant `gorm:"foreignKey:VariantID;references:ID"`
}

func (OrderDetail) TableName() string {
	return "order-detail"
}
