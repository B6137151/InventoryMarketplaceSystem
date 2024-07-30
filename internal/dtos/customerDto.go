package dtos

import "github.com/google/uuid"

// CustomerCreateDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการสร้าง Customer ใหม่
type CustomerCreateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// CustomerUpdateDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการอัปเดต Customer
type CustomerUpdateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// CustomerResponseDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการตอบกลับข้อมูล Customer
type CustomerResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
