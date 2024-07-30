package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Order represents an order placed by a customer
type Order struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt       time.Time      `gorm:"type:timestamp with time zone"`
	UpdatedAt       time.Time      `gorm:"type:timestamp with time zone"`
	DeletedAt       gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
	OrderID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()"`
	CustomerID      uuid.UUID      `gorm:"type:uuid;not null"`
	RoundID         uuid.UUID      `gorm:"type:uuid;not null"`
	OrderDate       time.Time      `gorm:"not null"`
	Status          string         `gorm:"type:varchar(100);not null"`
	Code            string         `gorm:"type:varchar(100);not null"`
	TotalPrice      float64        `gorm:"not null"`
	DeliveryAddress string         `gorm:"type:varchar(255);not null"`
	PaymentSource   string         `gorm:"type:varchar(100);not null"`

	Customer     Customer       `gorm:"foreignKey:CustomerID;references:ID"`
	SalesRound   SalesRound     `gorm:"foreignKey:RoundID;references:ID"`
	OrderDetail  []OrderDetail  `gorm:"foreignKey:OrderID"`
	OrderHistory []OrderHistory `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "order"
}
