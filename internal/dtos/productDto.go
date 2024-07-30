package dtos

import "github.com/google/uuid"

// ProductCreateDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการสร้าง Product ใหม่
type ProductCreateDTO struct {
	StoreID     uuid.UUID `json:"store_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	ProductName string    `json:"product_name" validate:"required"`
	Brand       string    `json:"brand" validate:"required"`
	Description string    `json:"description"`
	Currency    string    `json:"currency" validate:"required"`
	Stock       int       `json:"stock" validate:"required"`
	Price       float64   `json:"price" validate:"required"`     // Add Price field
	ImageURL    string    `json:"image_url" validate:"required"` // Add ImageURL field
}

// ProductUpdateDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการอัปเดต Product
type ProductUpdateDTO struct {
	StoreID     uuid.UUID `json:"store_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	ProductName string    `json:"product_name" validate:"required"`
	Brand       string    `json:"brand" validate:"required"`
	Description string    `json:"description"`
	Currency    string    `json:"currency" validate:"required"`
	Stock       int       `json:"stock" validate:"required"`
	Price       float64   `json:"price" validate:"required"`     // Add Price field
	ImageURL    string    `json:"image_url" validate:"required"` // Add ImageURL field
}

// ProductResponseDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการตอบกลับข้อมูล Product
type ProductResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	StoreID     uuid.UUID `json:"store_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	ProductName string    `json:"product_name"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	Stock       int       `json:"stock"`
	Price       float64   `json:"price"`     // Add Price field
	ImageURL    string    `json:"image_url"` // Add ImageURL field
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
