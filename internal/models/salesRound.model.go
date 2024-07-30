package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SalesRound struct {
	gorm.Model
	ID        uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time          `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"type:timestamp with time zone;index"`
	Name      string             `gorm:"size:100;not null" json:"name"`
	StartDate time.Time          `gorm:"type:timestamp with time zone;not null" json:"start_date"`
	EndDate   time.Time          `gorm:"type:timestamp with time zone;not null" json:"end_date"`
	Details   []SalesRoundDetail `gorm:"foreignKey:RoundID"` // One-to-many relationship with SalesRoundDetail
	Orders    []Order            `gorm:"foreignKey:RoundID"`
}

// TableName sets the table name explicitly for the SalesRound model
func (SalesRound) TableName() string {
	return "sales-round" // Ensures the table name is exactly as specified here
}
