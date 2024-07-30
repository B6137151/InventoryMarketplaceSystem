package models

//
//import (
//	"time"
//
//	"github.com/google/uuid"
//)
//
//// Purchase represents a customer purchase
//type Purchase struct {
//	ID              uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
//	CustomerID      uuid.UUID     `gorm:"type:uuid;not null;index"`
//	RoundID         uuid.UUID     `gorm:"type:uuid;not null;index"`
//	PurchaseDate    time.Time     `gorm:"type:timestamp with time zone;not null;default:current_timestamp"`
//	TotalPrice      float64       `gorm:"type:numeric;not null"`
//	DeliveryAddress string        `gorm:"type:varchar(255);not null"`
//	PaymentSource   string        `gorm:"type:varchar(100);not null"`
//	CreatedAt       time.Time     `gorm:"type:timestamp with time zone;autoCreateTime"`
//	UpdatedAt       time.Time     `gorm:"type:timestamp with time zone;autoUpdateTime"`
//	OrderDetails    []OrderDetail `gorm:"foreignKey:PurchaseID"`
//}
//
//// TableName sets the table name for the Purchase struct
//func (Purchase) TableName() string {
//	return "purchase"
//}
