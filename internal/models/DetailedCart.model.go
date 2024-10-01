// DetailCart.model
package models

import (
	"github.com/google/uuid"
)

// DetailedCart represents a cart with its detailed information.
type DetailedCart struct {
	CustomerID uuid.UUID          `json:"customer_id"`
	RoundID    uuid.UUID          `json:"round_id"`
	RoundName  string             `json:"round_name"`
	StoreID    uuid.UUID          `json:"store_id"`
	StoreName  string             `json:"store_name"`
	Items      []DetailedCartItem `json:"items" gorm:"-"`
}

// DetailedCartItem represents an item in the detailed cart.
type DetailedCartItem struct {
	VariantID uuid.UUID `json:"variant_id"`
	SKUCode   string    `json:"sku_code"`
	Price     float64   `json:"price"`
	ImageURL  string    `json:"image_url"`
	Quantity  int       `json:"quantity"`
}
