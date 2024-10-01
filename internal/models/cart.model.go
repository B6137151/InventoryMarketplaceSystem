package models

import (
	"github.com/google/uuid"
)

// CartItem represents an item in the cart.
type CartItem struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CartID      uuid.UUID `json:"cart_id" gorm:"type:uuid;not null"`
	VariantID   uuid.UUID `json:"variant_id" gorm:"type:uuid;not null"`
	ProductName string    `gorm:"type:varchar(255);not null"` // New field for product name
	SKUCode     string    `gorm:"type:varchar(100);not null"` // New field for SKU code
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"` // New field for price

	Cart           Cart           `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;"`
	ProductVariant ProductVariant `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE;"`
}

// Cart represents the shopping cart.
type Cart struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID      uuid.UUID  `gorm:"type:uuid;not null"`
	RoundID         uuid.UUID  `gorm:"type:uuid;not null"`
	StoreID         uuid.UUID  `gorm:"type:uuid;not null"`
	Items           []CartItem `gorm:"foreignKey:CartID"`
	TotalAmount     float64    `gorm:"not null"`                   // New field for total amount
	PaymentSource   string     `gorm:"type:varchar(100);not null"` // New field for payment source
	DeliveryAddress string     `gorm:"type:varchar(255);not null"` // New field for delivery address
	Currency        string     `gorm:"type:varchar(10);not null"`  // Add this field for currency
	// Relationships
	Customer   Customer   `gorm:"foreignKey:CustomerID" json:"customer"`
	SalesRound SalesRound `gorm:"foreignKey:RoundID" json:"sales_round"`
	Store      Store      `gorm:"foreignKey:StoreID" json:"store"`
}

func (Cart) TableName() string {
	return "cart"
}

func (CartItem) TableName() string {
	return "cart_item"
}
