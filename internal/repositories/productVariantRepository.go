package repositories

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductVariantRepository interface {
	CreateProductVariant(productVariant *models.ProductVariant) error
	GetAllProductVariants() ([]models.ProductVariant, error)
	GetProductVariantByID(id uuid.UUID) (*models.ProductVariant, error)
	UpdateProductVariant(productVariant *models.ProductVariant) error
	DeleteProductVariant(id uuid.UUID) error
}

type productVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) ProductVariantRepository {
	return &productVariantRepository{db: db}
}

func (r *productVariantRepository) CreateProductVariant(productVariant *models.ProductVariant) error {
	return r.db.Create(productVariant).Error
}

func (r *productVariantRepository) GetAllProductVariants() ([]models.ProductVariant, error) {
	var productVariants []models.ProductVariant
	err := r.db.Find(&productVariants).Error
	return productVariants, err
}

func (r *productVariantRepository) GetProductVariantByID(id uuid.UUID) (*models.ProductVariant, error) {
	var productVariant models.ProductVariant
	err := r.db.First(&productVariant, "variant_id = ?", id).Error
	return &productVariant, err
}

func (r *productVariantRepository) UpdateProductVariant(productVariant *models.ProductVariant) error {
	return r.db.Save(productVariant).Error
}

func (r *productVariantRepository) DeleteProductVariant(id uuid.UUID) error {
	return r.db.Delete(&models.ProductVariant{}, "variant_id = ?", id).Error
}
