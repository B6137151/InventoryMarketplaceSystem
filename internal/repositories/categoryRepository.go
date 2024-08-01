package repositories

import (
	"log"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetAllCategories() ([]models.Category, error)
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uuid.UUID) error
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

func (r *categoryRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in GetAllCategories: %v", r)
		}
	}()
	err := r.db.Find(&categories).Error
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
