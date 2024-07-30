package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderHistory struct {
	gorm.Model            // Includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	HistoryID   uint      `gorm:"primaryKey;autoIncrement"` // Primary key with auto-increment
	OrderID     uint      `gorm:"not null;index"`           // Foreign key for the Order
	Status      string    `gorm:"size:100;not null"`        // Status of the order at this history point
	ChangedAt   time.Time `gorm:"not null"`                 // Timestamp when the status change occurred
	Description string    `gorm:"type:text;not null"`       // Description of the status change

}

func (OrderHistory) TableName() string {
	return "order-history"
}
