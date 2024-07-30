package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductVariant represents the details of product variants in the database.
type ProductVariant struct {
	gorm.Model                           // Includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	VariantID         uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"` // Primary key with auto-generated UUID
	ProductID         uuid.UUID          `gorm:"type:uuid;not null;index"`                       // Foreign key for the Product
	SKUCode           string             `gorm:"size:100;not null;unique"`                       // Stock Keeping Unit code, unique
	Price             float64            `gorm:"not null"`                                       // Price of the product variant
	ImageURL          string             `gorm:"size:255"`                                       // URL to the image of the product variant
	SalesRoundDetails []SalesRoundDetail `gorm:"foreignKey:VariantID"`                           // One-to-many relationship with SalesRoundDetail
	OrderDetails      []OrderDetail      `gorm:"foreignKey:VariantID"`                           // One-to-many relationship with OrderDetail
}

func (ProductVariant) TableName() string {
	return "product-variant" // Ensures the table name is exactly as specified here
}
