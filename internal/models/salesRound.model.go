package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalesRound struct {
	ID        uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string             `gorm:"size:255;not null"`
	StartDate time.Time          `gorm:"not null"`
	EndDate   time.Time          `gorm:"not null"`
	CreatedAt time.Time          `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"type:timestamp with time zone;index"`
	Details   []SalesRoundDetail `gorm:"foreignKey:RoundID;constraint:OnDelete:CASCADE;"` // One-to-many relationship with SalesRoundDetail
	Orders    []Order            `gorm:"foreignKey:RoundID;constraint:OnDelete:CASCADE;"`
}

// TableName sets the table name explicitly for the SalesRound model
func (SalesRound) TableName() string {
	return "sales-round" // Ensures the table name is exactly as specified here
}
