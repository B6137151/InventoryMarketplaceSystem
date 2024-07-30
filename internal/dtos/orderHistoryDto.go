package dtos

import "time"

// OrderHistoryCreateDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการสร้าง OrderHistory ใหม่
type OrderHistoryCreateDTO struct {
	OrderID     uint      `json:"order_id" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	ChangedAt   time.Time `json:"changed_at" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

// OrderHistoryUpdateDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการอัปเดต OrderHistory
type OrderHistoryUpdateDTO struct {
	Status      string    `json:"status" validate:"required"`
	ChangedAt   time.Time `json:"changed_at" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

// OrderHistoryResponseDTO เป็นโครงสร้างข้อมูลที่ใช้สำหรับการตอบกลับข้อมูล OrderHistory
type OrderHistoryResponseDTO struct {
	ID          uint      `json:"id"`
	OrderID     uint      `json:"order_id"`
	Status      string    `json:"status"`
	ChangedAt   time.Time `json:"changed_at"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
