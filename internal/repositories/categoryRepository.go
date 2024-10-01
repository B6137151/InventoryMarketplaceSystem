package repositories

import (
	"log"

	"strings"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetAllCategories(expand string) ([]models.Category, error) // Updated signature
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uuid.UUID) error
	GetCategoriesBySalesRoundID(salesRoundID uuid.UUID) ([]models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *models.Category) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in CreateCategory: %v", r)
		}
	}()
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetAllCategories(expand string) ([]models.Category, error) {
	var categories []models.Category
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in GetAllCategories: %v", r)
		}
	}()

	query := r.db.Model(&models.Category{})

	// Conditionally preload related entities based on the 'expand' parameter
	if expand != "" {
		if strings.Contains(expand, "store") {
			query = query.Preload("Store")
		}
		// Add more conditions if there are other relations you want to support
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in GetCategoryByID: %v", r)
		}
	}()
	err := r.db.First(&category, "id = ?", id).Error
	return &category, err
}

func (r *categoryRepository) UpdateCategory(category *models.Category) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in UpdateCategory: %v", r)
		}
	}()
	return r.db.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(id uuid.UUID) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in DeleteCategory: %v", r)
		}
	}()
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}
func (r *categoryRepository) GetCategoriesBySalesRoundID(salesRoundID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in GetCategoriesBySalesRoundID: %v", r)
		}
	}()

	err := r.db.Joins(`JOIN product ON product.category_id = category.id`).
		Joins(`JOIN "sales-round-detail" ON "sales-round-detail".product_id = product.id`).
		Where(`"sales-round-detail".round_id = ?`, salesRoundID).
		Find(&categories).Error
	return categories, err
}
