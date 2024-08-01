// internal/repositories/purchaseRepository.go
package repositories

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseRepository interface {
	CreateOrder(order *models.Order) error
	CreateOrderDetail(orderDetail *models.OrderDetail) error
	GetProductVariantByID(id uuid.UUID) (*models.ProductVariant, error)
	UpdateProductStock(product *models.Product) error
}

type purchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *purchaseRepository) CreateOrderDetail(orderDetail *models.OrderDetail) error {
	return r.db.Create(orderDetail).Error
}

func (r *purchaseRepository) GetProductVariantByID(id uuid.UUID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Preload("Product").First(&variant, "variant_id = ?", id).Error
	return &variant, err
}

func (r *purchaseRepository) UpdateProductStock(product *models.Product) error {
	return r.db.Save(product).Error
}
