package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SalesRoundDetail struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt      time.Time      `gorm:"type:timestamp with time zone;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"type:timestamp with time zone;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
	RoundID        uuid.UUID      `gorm:"type:uuid;not null;index"` // Foreign key for the SalesRound
	VariantID      uuid.UUID      `gorm:"type:uuid;not null;index"` // Foreign key for the ProductVariant
	Quantity       int            `gorm:"not null"`                 // Quantity of product variants allocated to this sales round
	Remaining      int            `gorm:"not null"`                 // Remaining quantity of product variants available in the sales round
	ProductStock   int            `gorm:"not null"`                 // Product stock available for this sales round detail
	QuantityLimit  int            `gorm:"not null"`                 // Quantity limit for this sales round detail
	SalesRound     SalesRound     `gorm:"foreignKey:RoundID"`       // Many-to-One relationship with SalesRound
	ProductVariant ProductVariant `gorm:"foreignKey:VariantID"`     // Many-to-One relationship with ProductVariant
}

// TableName sets the table name explicitly for the SalesRoundDetail model
func (SalesRoundDetail) TableName() string {
	return "sales-round-detail"
}
