package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Store struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
	StoreName string         `gorm:"size:255;not null"`
	Location  string         `gorm:"size:255"`
	Products  []Product      `gorm:"foreignKey:StoreID"` // Ensure StoreID in Product points to this ID
}

func (Store) TableName() string {
	return "store"
}
