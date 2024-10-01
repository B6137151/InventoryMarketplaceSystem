package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"` // Ensure your DB supports this function
	CreatedAt time.Time      `gorm:"autoCreateTime"`                                 // Let GORM handle the timestamp fields automatically
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	StoreName string         `gorm:"size:255;not null"`
	Location  string         `gorm:"size:255"`
	Products  []Product      `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Adds DB-level foreign key constraint
}

func (Store) TableName() string {
	return "store"
}
