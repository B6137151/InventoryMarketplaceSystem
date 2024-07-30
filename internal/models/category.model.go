package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Category represents a product category
type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
	Name      string         `gorm:"size:255;not null"`
	StoreID   uuid.UUID      `gorm:"type:uuid;not null;index"` // Store ID for the category
	Products  []Product      `gorm:"foreignKey:CategoryID"`    // One-to-many relationship with products
}

func (Category) TableName() string {
	return "category"
}
