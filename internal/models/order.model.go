package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type PaymentInfo struct {
	PaymentMethod string  `gorm:"size:100;not null"`
	PaymentStatus string  `gorm:"size:100;not null"`
	TransactionID string  `gorm:"size:100;not null"`
	PaidAmount    float64 `gorm:"not null"`         // Added field for the amount paid
	PaidCurrency  string  `gorm:"size:10;not null"` // Added field for the currency used in payment
}

type ShippingInfo struct {
	ShippingMethod        string    `gorm:"size:100;not null"`
	EstimatedDeliveryDate time.Time `gorm:"not null"`
	TrackingNumber        string    `gorm:"size:100"` // Added field for tracking number
	CarrierName           string    `gorm:"size:100"` // Added field for carrier name
	ShippingStatus        string    `gorm:"size:100"` // Added field for shipping status
	EstimatedArrival      time.Time `gorm:"not null"` // Added field for estimated arrival date
}

type Order struct {
	ID              uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CartID          uuid.UUID    `gorm:"type:uuid;not null;index"` // Foreign key to the Cart
	CustomerID      uuid.UUID    `gorm:"type:uuid;not null"`
	RoundID         uuid.UUID    `gorm:"type:uuid;not null"`
	OrderDate       time.Time    `gorm:"type:timestamp with time zone;not null"`
	Status          OrderStatus  `gorm:"type:varchar(20);not null"`
	Code            string       `gorm:"type:varchar(100);not null"`
	TotalPrice      float64      `gorm:"not null"`
	DeliveryAddress string       `gorm:"type:varchar(255);not null"`
	PaymentSource   string       `gorm:"type:varchar(100);not null"` // Now a regular string
	PaymentInfo     PaymentInfo  `gorm:"embedded"`                   // Embedding PaymentInfo
	ShippingInfo    ShippingInfo `gorm:"embedded"`                   // Embedding ShippingInfo

	CreatedAt time.Time      `gorm:"type:timestamp with time zone"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`

	Customer     Customer       `gorm:"foreignKey:CustomerID;references:ID"`
	SalesRound   SalesRound     `gorm:"foreignKey:RoundID;references:ID"`
	OrderDetail  []OrderDetail  `gorm:"foreignKey:OrderID"`
	OrderHistory []OrderHistory `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}
