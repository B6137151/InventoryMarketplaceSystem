package repositories

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type SalesRoundRepository interface {
	CreateSalesRound(salesRound *models.SalesRound) error
	GetAllSalesRounds() ([]models.SalesRound, error)
	GetSalesRoundByID(id uuid.UUID) (*models.SalesRound, error)
	UpdateSalesRound(salesRound *models.SalesRound) error
	DeleteSalesRound(id uuid.UUID) error
	GetCombinedSalesRoundProductData() ([]dtos.CombinedSalesRoundProductResponse, error) // New method
}

type salesRoundRepository struct {
	db *gorm.DB
}

func NewSalesRoundRepository(db *gorm.DB) SalesRoundRepository {
	return &salesRoundRepository{db: db}
}

func (r *salesRoundRepository) CreateSalesRound(salesRound *models.SalesRound) error {
	return r.db.Create(salesRound).Error
}

func (r *salesRoundRepository) GetAllSalesRounds() ([]models.SalesRound, error) {
	var salesRounds []models.SalesRound
	err := r.db.Find(&salesRounds).Error
	return salesRounds, err
}

func (r *salesRoundRepository) GetSalesRoundByID(id uuid.UUID) (*models.SalesRound, error) {
	var salesRound models.SalesRound
	err := r.db.First(&salesRound, "id = ?", id).Error
	return &salesRound, err
}

func (r *salesRoundRepository) UpdateSalesRound(salesRound *models.SalesRound) error {
	return r.db.Save(salesRound).Error
}

func (r *salesRoundRepository) DeleteSalesRound(id uuid.UUID) error {
	return r.db.Delete(&models.SalesRound{}, "id = ?", id).Error
}

func (r *salesRoundRepository) GetCombinedSalesRoundProductData() ([]dtos.CombinedSalesRoundProductResponse, error) {
	log.Println("Starting to fetch combined sales round product data")
	var results []dtos.CombinedSalesRoundProductResponse
	err := r.db.Table("\"sales-round\"").
		Select("\"sales-round\".id as sales_round_id, \"sales-round\".name as sales_round_name, \"sales-round\".start_date, \"sales-round\".end_date, \"sales-round\".created_at, \"sales-round\".updated_at, " +
			"product.id as product_id, product.store_id, product.category_id, product.product_name, product.brand, product.description, product.currency, product.stock, product.price, product.image_url, product.created_at as product_created_at, product.updated_at as product_updated_at, " +
			"\"product-variant\".id as variant_id, \"product-variant\".sku_code, \"product-variant\".price as variant_price, \"product-variant\".image_url as variant_image_url, \"product-variant\".created_at as variant_created_at, \"product-variant\".updated_at as variant_updated_at").
		Joins("left join product on product.id = \"sales-round\".id").
		Joins("left join \"product-variant\" on \"product-variant\".product_id = product.id").
		Scan(&results).Error
	if err != nil {
		log.Println("Error during fetch:", err)
	} else {
		log.Println("Successfully fetched data:", results)
	}
	return results, err
}
