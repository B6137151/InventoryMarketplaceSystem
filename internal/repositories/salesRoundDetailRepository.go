package repositories

import (
	"fmt"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalesRoundDetailRepository interface {
	CreateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) error
	GetAllSalesRoundDetails() ([]models.SalesRoundDetail, error)
	GetSalesRoundDetailByID(id uuid.UUID) (*models.SalesRoundDetail, error)
	UpdateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) error
	DeleteSalesRoundDetail(id uuid.UUID) error
	GetSalesRoundDetailsByRoundID(roundID uuid.UUID) ([]dtos.CombinedSalesRoundDetailResponse, error)
	UpdateSalesRoundDetailQuantity(id uuid.UUID, quantity int) error
	GetProductVariantByID(id uuid.UUID) (*models.ProductVariant, error)
	GetProductByVariantID(variantID uuid.UUID) (*models.Product, error)
	UpdateProductStock(product *models.Product) error
}

type salesRoundDetailRepository struct {
	db *gorm.DB
}

func NewSalesRoundDetailRepository(db *gorm.DB) SalesRoundDetailRepository {
	return &salesRoundDetailRepository{db: db}
}

func (r *salesRoundDetailRepository) CreateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) error {
	// Fetch the product associated with the product variant
	product, err := r.GetProductByVariantID(salesRoundDetail.VariantID)
	if err != nil {
		return err
	}

	var existingDetail models.SalesRoundDetail
	err = r.db.Where("round_id = ? AND variant_id = ?", salesRoundDetail.RoundID, salesRoundDetail.VariantID).First(&existingDetail).Error

	// If the sales round detail already exists
	if err == nil {
		// Calculate the total quantity to be updated
		totalQuantity := existingDetail.Quantity + salesRoundDetail.Quantity

		// Check if there is enough stock for the total quantity
		if totalQuantity > product.Stock+existingDetail.Quantity {
			return fmt.Errorf("quantity exceeds available stock")
		}

		// Update the existing sales round detail
		product.Stock -= salesRoundDetail.Quantity
		existingDetail.Quantity = totalQuantity
		existingDetail.Remaining = product.Stock
		existingDetail.ProductStock = product.Stock

		// Update the product stock
		if err := r.UpdateProductStock(product); err != nil {
			return err
		}

		return r.db.Save(&existingDetail).Error
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// Check if there is enough stock for the new entry
	if salesRoundDetail.Quantity > product.Stock {
		return fmt.Errorf("quantity exceeds available stock")
	}

	// Allocate the stock for the new entry
	product.Stock -= salesRoundDetail.Quantity

	// Update the product stock
	if err := r.UpdateProductStock(product); err != nil {
		return err
	}

	// Set the remaining stock in the sales round detail
	salesRoundDetail.ProductStock = product.Stock
	salesRoundDetail.Remaining = product.Stock

	// Create the sales round detail
	return r.db.Create(salesRoundDetail).Error
}

func (r *salesRoundDetailRepository) GetAllSalesRoundDetails() ([]models.SalesRoundDetail, error) {
	var salesRoundDetails []models.SalesRoundDetail
	err := r.db.Find(&salesRoundDetails).Error
	return salesRoundDetails, err
}

func (r *salesRoundDetailRepository) GetSalesRoundDetailByID(id uuid.UUID) (*models.SalesRoundDetail, error) {
	var salesRoundDetail models.SalesRoundDetail
	err := r.db.First(&salesRoundDetail, "id = ?", id).Error
	return &salesRoundDetail, err
}

func (r *salesRoundDetailRepository) UpdateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) error {
	return r.db.Save(salesRoundDetail).Error
}

func (r *salesRoundDetailRepository) DeleteSalesRoundDetail(id uuid.UUID) error {
	return r.db.Delete(&models.SalesRoundDetail{}, "id = ?", id).Error
}

func (r *salesRoundDetailRepository) GetSalesRoundDetailsByRoundID(roundID uuid.UUID) ([]dtos.CombinedSalesRoundDetailResponse, error) {
	var details []dtos.CombinedSalesRoundDetailResponse
	err := r.db.Table("\"sales-round-detail\"").
		Select("\"sales-round\".*, \"sales-round-detail\".*, \"product-variant\".sku_code, \"product-variant\".price AS variant_price, \"product-variant\".image_url AS variant_image_url, product.product_name, product.brand, product.description, product.currency, product.stock, product.price AS product_price").
		Joins("JOIN \"sales-round\" ON \"sales-round-detail\".round_id = \"sales-round\".id").
		Joins("JOIN \"product-variant\" ON \"sales-round-detail\".variant_id = \"product-variant\".variant_id").
		Joins("JOIN product ON \"product-variant\".product_id = product.id").
		Where("\"sales-round-detail\".round_id = ?", roundID).
		Scan(&details).Error
	return details, err
}

func (r *salesRoundDetailRepository) UpdateSalesRoundDetailQuantity(id uuid.UUID, quantity int) error {
	var detail models.SalesRoundDetail
	if err := r.db.First(&detail, "id = ?", id).Error; err != nil {
		return err
	}

	product, err := r.GetProductByVariantID(detail.VariantID)
	if err != nil {
		return err
	}

	if quantity > product.Stock+detail.Quantity {
		return fmt.Errorf("cannot update quantity beyond available stock")
	}

	// Adjust the product stock based on the new quantity
	stockChange := detail.Quantity - quantity
	product.Stock += stockChange

	// Update the product stock
	if err := r.UpdateProductStock(product); err != nil {
		return err
	}

	// Update the sales round detail quantity
	detail.Quantity = quantity
	detail.Remaining = product.Stock
	return r.db.Save(&detail).Error
}

func (r *salesRoundDetailRepository) GetProductVariantByID(id uuid.UUID) (*models.ProductVariant, error) {
	var productVariant models.ProductVariant
	err := r.db.First(&productVariant, "id = ?", id).Error
	return &productVariant, err
}

func (r *salesRoundDetailRepository) GetProductByVariantID(variantID uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Table("product").
		Select("product.*").
		Joins("JOIN \"product-variant\" ON \"product-variant\".product_id = product.id").
		Where("\"product-variant\".variant_id = ?", variantID).
		First(&product).Error
	return &product, err
}

func (r *salesRoundDetailRepository) UpdateProductStock(product *models.Product) error {
	return r.db.Save(product).Error
}
