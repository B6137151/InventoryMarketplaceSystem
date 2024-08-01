package repositories

import (
	"fmt"
	"log"

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
	GetSalesRoundDetailsByVariantID(variantID uuid.UUID) ([]models.SalesRoundDetail, error)
	GetSalesRoundDetailByRoundIDAndVariantID(roundID uuid.UUID, variantID uuid.UUID) (*models.SalesRoundDetail, error)
	UpdateSalesRoundDetailByRoundIDAndVariantID(roundID uuid.UUID, variantID uuid.UUID, salesRoundDetail *models.SalesRoundDetail) error
}

type salesRoundDetailRepository struct {
	db *gorm.DB
}

func NewSalesRoundDetailRepository(db *gorm.DB) SalesRoundDetailRepository {
	return &salesRoundDetailRepository{db: db}
}

func (r *salesRoundDetailRepository) CreateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) error {
	// Fetch the product associated with the product variant
	log.Printf("Fetching product by variant ID: %v", salesRoundDetail.VariantID)
	product, err := r.GetProductByVariantID(salesRoundDetail.VariantID)
	if err != nil {
		log.Printf("Error fetching product by variant ID: %v, error: %v", salesRoundDetail.VariantID, err)
		return err
	}
	log.Printf("Product found: %v", product)

	var existingDetail models.SalesRoundDetail
	log.Printf("Fetching existing sales round detail for round ID: %v and variant ID: %v", salesRoundDetail.RoundID, salesRoundDetail.VariantID)
	err = r.db.Where("round_id = ? AND variant_id = ?", salesRoundDetail.RoundID, salesRoundDetail.VariantID).First(&existingDetail).Error

	// If the sales round detail already exists
	if err == nil {
		log.Printf("Existing sales round detail found: %v", existingDetail)
		// Calculate the total quantity to be updated
		totalQuantity := existingDetail.Quantity + salesRoundDetail.Quantity

		// Check if there is enough stock for the total quantity
		if totalQuantity > product.Stock+existingDetail.Quantity {
			log.Printf("Quantity exceeds available stock: %d > %d", totalQuantity, product.Stock+existingDetail.Quantity)
			return fmt.Errorf("quantity exceeds available stock")
		}

		// Update the existing sales round detail
		product.Stock -= salesRoundDetail.Quantity
		existingDetail.Quantity = totalQuantity
		existingDetail.Remaining = product.Stock
		existingDetail.ProductStock = product.Stock

		// Update the product stock
		if err := r.UpdateProductStock(product); err != nil {
			log.Printf("Error updating product stock: %v", err)
			return err
		}

		log.Printf("Updating existing sales round detail: %v", existingDetail)
		return r.db.Save(&existingDetail).Error
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Error fetching existing sales round detail: %v", err)
		return err
	}

	// Check if there is enough stock for the new entry
	if salesRoundDetail.Quantity > product.Stock {
		log.Printf("Quantity exceeds available stock: %d > %d", salesRoundDetail.Quantity, product.Stock)
		return fmt.Errorf("quantity exceeds available stock")
	}

	// Allocate the stock for the new entry
	product.Stock -= salesRoundDetail.Quantity

	// Update the product stock
	if err := r.UpdateProductStock(product); err != nil {
		log.Printf("Error updating product stock: %v", err)
		return err
	}

	// Set the remaining stock in the sales round detail
	salesRoundDetail.ProductStock = product.Stock
	salesRoundDetail.Remaining = product.Stock

	log.Printf("Creating new sales round detail: %v", salesRoundDetail)
	// Create the sales round detail
	return r.db.Create(salesRoundDetail).Error
}

func (r *salesRoundDetailRepository) GetAllSalesRoundDetails() ([]models.SalesRoundDetail, error) {
	var salesRoundDetails []models.SalesRoundDetail
	err := r.db.Find(&salesRoundDetails).Error
	log.Printf("Fetched all sales round details: %v, error: %v", salesRoundDetails, err)
	return salesRoundDetails, err
}

func (r *salesRoundDetailRepository) GetSalesRoundDetailByID(id uuid.UUID) (*models.SalesRoundDetail, error) {
	var salesRoundDetail models.SalesRoundDetail
	err := r.db.First(&salesRoundDetail, "id = ?", id).Error
	log.Printf("Fetched sales round detail by ID %v: %v, error: %v", id, salesRoundDetail, err)
	return &salesRoundDetail, err
}

func (r *salesRoundDetailRepository) UpdateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) error {
	log.Printf("Updating sales round detail: %v", salesRoundDetail)
	return r.db.Save(salesRoundDetail).Error
}

func (r *salesRoundDetailRepository) DeleteSalesRoundDetail(id uuid.UUID) error {
	log.Printf("Deleting sales round detail with ID %v", id)
	return r.db.Delete(&models.SalesRoundDetail{}, "id = ?", id).Error
}

func (r *salesRoundDetailRepository) GetSalesRoundDetailsByRoundID(roundID uuid.UUID) ([]dtos.CombinedSalesRoundDetailResponse, error) {
	var details []dtos.CombinedSalesRoundDetailResponse
	log.Printf("Fetching sales round details by round ID: %v", roundID)
	err := r.db.Table("\"sales-round-detail\"").
		Select("\"sales-round\".*, \"sales-round-detail\".*, \"product-variant\".sku_code, \"product-variant\".price AS variant_price, \"product-variant\".image_url AS variant_image_url, product.product_name, product.brand, product.description, product.currency, product.stock, product.price AS product_price").
		Joins("JOIN \"sales-round\" ON \"sales-round-detail\".round_id = \"sales-round\".id").
		Joins("JOIN \"product-variant\" ON \"sales-round-detail\".variant_id = \"product-variant\".variant_id").
		Joins("JOIN product ON \"product-variant\".product_id = product.id").
		Where("\"sales-round-detail\".round_id = ?", roundID).
		Scan(&details).Error
	log.Printf("Fetched sales round details: %v, error: %v", details, err)
	return details, err
}

func (r *salesRoundDetailRepository) UpdateSalesRoundDetailQuantity(id uuid.UUID, quantity int) error {
	var detail models.SalesRoundDetail
	if err := r.db.First(&detail, "id = ?", id).Error; err != nil {
		log.Printf("Error fetching sales round detail by ID %v: %v", id, err)
		return err
	}

	product, err := r.GetProductByVariantID(detail.VariantID)
	if err != nil {
		log.Printf("Error fetching product by variant ID %v: %v", detail.VariantID, err)
		return err
	}

	if quantity > product.Stock+detail.Quantity {
		log.Printf("Cannot update quantity beyond available stock: %d > %d", quantity, product.Stock+detail.Quantity)
		return fmt.Errorf("cannot update quantity beyond available stock")
	}

	// Adjust the product stock based on the new quantity
	stockChange := detail.Quantity - quantity
	product.Stock += stockChange

	// Update the product stock
	if err := r.UpdateProductStock(product); err != nil {
		log.Printf("Error updating product stock: %v", err)
		return err
	}

	// Update the sales round detail quantity
	detail.Quantity = quantity
	detail.Remaining = product.Stock
	log.Printf("Updating sales round detail quantity: %v", detail)
	return r.db.Save(&detail).Error
}

func (r *salesRoundDetailRepository) GetProductVariantByID(id uuid.UUID) (*models.ProductVariant, error) {
	var productVariant models.ProductVariant
	err := r.db.First(&productVariant, "id = ?", id).Error
	log.Printf("Fetched product variant by ID %v: %v, error: %v", id, productVariant, err)
	return &productVariant, err
}

func (r *salesRoundDetailRepository) GetProductByVariantID(variantID uuid.UUID) (*models.Product, error) {
	var product models.Product
	log.Printf("Fetching product by variant ID: %v", variantID)
	err := r.db.Table("product").
		Select("product.*").
		Joins("JOIN \"product-variant\" ON \"product-variant\".product_id = product.id").
		Where("\"product-variant\".variant_id = ?", variantID).
		First(&product).Error
	log.Printf("Fetched product by variant ID %v: %v, error: %v", variantID, product, err)
	return &product, err
}

func (r *salesRoundDetailRepository) UpdateProductStock(product *models.Product) error {
	log.Printf("Updating product stock: %v", product)
	return r.db.Save(product).Error
}

func (r *salesRoundDetailRepository) GetSalesRoundDetailsByVariantID(variantID uuid.UUID) ([]models.SalesRoundDetail, error) {
	var details []models.SalesRoundDetail
	log.Printf("Fetching sales round details by variant ID: %v", variantID)
	err := r.db.Where("variant_id = ?", variantID).Find(&details).Error
	log.Printf("Fetched sales round details by variant ID %v: %v, error: %v", variantID, details, err)
	return details, err
}

func (r *salesRoundDetailRepository) UpdateSalesRoundDetailByRoundIDAndVariantID(roundID uuid.UUID, variantID uuid.UUID, salesRoundDetail *models.SalesRoundDetail) error {
	var detail models.SalesRoundDetail
	err := r.db.Where("round_id = ? AND variant_id = ?", roundID, variantID).First(&detail).Error
	if err != nil {
		log.Printf("Error fetching sales round detail by round ID %v and variant ID %v: %v", roundID, variantID, err)
		return err
	}

	log.Printf("Updating sales round detail: %v", detail)
	return r.db.Model(&detail).Updates(salesRoundDetail).Error
}

func (r *salesRoundDetailRepository) GetSalesRoundDetailByRoundIDAndVariantID(roundID uuid.UUID, variantID uuid.UUID) (*models.SalesRoundDetail, error) {
	var salesRoundDetail models.SalesRoundDetail
	err := r.db.Where("round_id = ? AND variant_id = ?", roundID, variantID).First(&salesRoundDetail).Error
	return &salesRoundDetail, err
}
