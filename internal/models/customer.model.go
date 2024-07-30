package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Customer represents a customer in the system
type Customer struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
	Name      string         `gorm:"size:255;not null"`
	Email     string         `gorm:"size:255;unique;not null"`
	Orders    []Order        `gorm:"foreignKey:CustomerID"` // One-to-many relationship with orders
}

func (Customer) TableName() string {
	return "customer"
}
