package repositories

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uuid.UUID) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uuid.UUID) error
	GetAllProductsWithVariants() ([]models.Product, error) // New method
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Preload("ProductVariant").Find(&products).Error
	return products, err
}

func (r *productRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("ProductVariant").First(&product, "id = ?", id).Error
	return &product, err
}

func (r *productRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) DeleteProduct(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

func (r *productRepository) GetAllProductsWithVariants() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("ProductVariant").Find(&products).Error
	return products, err
}
