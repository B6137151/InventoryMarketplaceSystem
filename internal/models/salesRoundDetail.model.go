package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalesRoundDetail struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"` // UUID จะถูกสร้างโดยใช้ฟังก์ชัน gen_random_uuid() ของ PostgreSQL
	CreatedAt      time.Time      `gorm:"type:timestamp with time zone;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"type:timestamp with time zone;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`
	RoundID        uuid.UUID      `gorm:"type:uuid;not null;index"` // Foreign key สำหรับ SalesRound
	VariantID      uuid.UUID      `gorm:"type:uuid;not null;index"` // Foreign key สำหรับ ProductVariant
	Quantity       int            `gorm:"not null"`                 // จำนวน product variants ที่จัดสรรให้กับ sales round นี้
	Remaining      int            `gorm:"not null"`                 // จำนวน product variants ที่เหลืออยู่ใน sales round
	ProductStock   int            `gorm:"not null"`                 // จำนวนสินค้าคงเหลือที่สามารถใช้ได้ใน sales round นี้
	QuantityLimit  int            `gorm:"not null"`                 // จำนวนสูงสุดที่สามารถซื้อได้ใน sales round นี้
	SalesRound     SalesRound     `gorm:"foreignKey:RoundID"`       // ความสัมพันธ์แบบ Many-to-One กับ SalesRound
	ProductVariant ProductVariant `gorm:"foreignKey:VariantID"`     // ความสัมพันธ์แบบ Many-to-One กับ ProductVariant
}

// TableName sets the table name explicitly for the SalesRoundDetail model
func (SalesRoundDetail) TableName() string {
	return "sales-round-detail" // กำหนดชื่อ table ในฐานข้อมูลเป็น sales-round-detail
}
