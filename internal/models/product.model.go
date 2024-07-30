package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID             uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt      time.Time        `gorm:"type:timestamp with time zone"`
	UpdatedAt      time.Time        `gorm:"type:timestamp with time zone"`
	DeletedAt      gorm.DeletedAt   `gorm:"type:timestamp with time zone;index"`
	StoreID        uuid.UUID        `gorm:"type:uuid;not null;index"`
	CategoryID     uuid.UUID        `gorm:"type:uuid;not null;index"`
	ProductName    string           `gorm:"size:255;not null"`
	Brand          string           `gorm:"size:255;not null"`
	Description    string           `gorm:"type:text"`
	Currency       string           `gorm:"size:3;not null"`
	Stock          int              `gorm:"default:0;check:stock >= 0"`
	Price          float64          `gorm:"not null;default:0"`
	ImageURL       string           `gorm:"size:255"`
	ProductVariant []ProductVariant `gorm:"foreignKey:ProductID"`
	Store          Store            `gorm:"foreignKey:StoreID" json:"-"`
	Category       Category         `gorm:"foreignKey:CategoryID"`
}

// TableName sets the name of the table in the database using snake_case.
func (Product) TableName() string {
	return "product"
}
