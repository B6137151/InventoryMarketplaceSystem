package repositories

import (
	"fmt"
	"log"

	"strings"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalesRoundDetailRepository interface {
	CreateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) (*dtos.SalesRoundDetailResponseDTO, error) // เปลี่ยนแปลง: แก้ไขให้รีเทิร์นค่าประเภท *dtos.SalesRoundDetailResponseDTO
	GetAllSalesRoundDetails(expand string) ([]models.SalesRoundDetail, error)
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
	CheckPurchaseLimit(roundID uuid.UUID, variantID uuid.UUID, quantity int) (bool, error)
}

type salesRoundDetailRepository struct {
	db *gorm.DB
}

func NewSalesRoundDetailRepository(db *gorm.DB) SalesRoundDetailRepository {
	return &salesRoundDetailRepository{db: db}
}

func (r *salesRoundDetailRepository) CreateSalesRoundDetail(salesRoundDetail *models.SalesRoundDetail) (*dtos.SalesRoundDetailResponseDTO, error) {
	log.Printf("Fetching product by variant ID: %v", salesRoundDetail.VariantID)

	// ดึงข้อมูล product โดยใช้ variant_id แทน id
	product, err := r.GetProductByVariantID(salesRoundDetail.VariantID)
	if err != nil {
		log.Printf("Error fetching product by variant ID: %v, error: %v", salesRoundDetail.VariantID, err)
		return nil, err
	}
	log.Printf("Product found: %v", product)

	var existingDetail models.SalesRoundDetail
	log.Printf("Fetching existing sales round detail for round ID: %v and variant ID: %v", salesRoundDetail.RoundID, salesRoundDetail.VariantID)

	// แก้ไขให้ใช้ variant_id ในการตรวจสอบว่ามี SalesRoundDetail อยู่แล้วหรือไม่
	err = r.db.Where("round_id = ? AND variant_id = ?", salesRoundDetail.RoundID, salesRoundDetail.VariantID).First(&existingDetail).Error

	if err == nil {
		log.Printf("Existing sales round detail found: %v", existingDetail)

		totalQuantity := existingDetail.Quantity + salesRoundDetail.Quantity

		if totalQuantity > product.Stock+existingDetail.Quantity {
			log.Printf("Quantity exceeds available stock: %d > %d", totalQuantity, product.Stock+existingDetail.Quantity)
			return nil, fmt.Errorf("quantity exceeds available stock")
		}

		product.Stock -= salesRoundDetail.Quantity
		existingDetail.Quantity = totalQuantity
		existingDetail.Remaining = product.Stock
		existingDetail.ProductStock = product.Stock

		if err := r.UpdateProductStock(product); err != nil {
			log.Printf("Error updating product stock: %v", err)
			return nil, err
		}

		log.Printf("Updating existing sales round detail: %v", existingDetail)
		if err := r.db.Save(&existingDetail).Error; err != nil {
			return nil, err
		}

		// เพิ่มการตรวจสอบผลการ preload ข้อมูล SalesRound และ ProductVariant
		log.Println("Preloading SalesRound and ProductVariant...")
		if err := r.db.Preload("SalesRound").
			Preload("ProductVariant").
			First(&existingDetail, existingDetail.ID).Error; err != nil {
			log.Printf("Error reloading SalesRoundDetail: %v", err)
			return nil, err
		}

		// ตรวจสอบค่าที่ถูก preload เข้ามา
		log.Printf("Loaded ProductVariant: %+v", existingDetail.ProductVariant)

		// สร้างและรีเทิร์น SalesRoundDetailResponseDTO
		response := &dtos.SalesRoundDetailResponseDTO{
			ID:            existingDetail.ID,
			RoundID:       existingDetail.RoundID,
			VariantID:     existingDetail.VariantID,
			Quantity:      existingDetail.Quantity,
			Remaining:     existingDetail.Remaining,
			QuantityLimit: existingDetail.QuantityLimit,
			ProductStock:  existingDetail.ProductStock,
			CreatedAt:     existingDetail.CreatedAt,
			UpdatedAt:     existingDetail.UpdatedAt,
			SalesRound: dtos.SalesRoundInfoResponse{
				ID:        existingDetail.SalesRound.ID,
				Name:      existingDetail.SalesRound.Name,
				StartDate: existingDetail.SalesRound.StartDate,
				EndDate:   existingDetail.SalesRound.EndDate,
				CreatedAt: existingDetail.SalesRound.CreatedAt,
				UpdatedAt: existingDetail.SalesRound.UpdatedAt,
			},
			ProductVariant: dtos.ProductVariantInfoResponse{
				CreatedAt: existingDetail.ProductVariant.CreatedAt,
				UpdatedAt: existingDetail.ProductVariant.UpdatedAt,
				ID:        existingDetail.ProductVariant.VariantID,
				ProductID: existingDetail.ProductVariant.ProductID,
				SKUCode:   existingDetail.ProductVariant.SKUCode,
				// VariantID: existingDetail.ProductVariant.VariantID,
				Price:    existingDetail.ProductVariant.Price,
				ImageURL: existingDetail.ProductVariant.ImageURL,
			},
		}

		return response, nil
	}

	// ตรวจสอบว่ามี stock เพียงพอสำหรับการเพิ่มข้อมูลใหม่หรือไม่
	if salesRoundDetail.Quantity > product.Stock {
		log.Printf("Quantity exceeds available stock: %d > %d", salesRoundDetail.Quantity, product.Stock)
		return nil, fmt.Errorf("quantity exceeds available stock")
	}

	// จัดสรร stock สำหรับการเพิ่มข้อมูลใหม่
	product.Stock -= salesRoundDetail.Quantity

	// อัปเดต stock ของสินค้า
	if err := r.UpdateProductStock(product); err != nil {
		log.Printf("Error updating product stock: %v", err)
		return nil, err
	}

	// กำหนดค่า stock ที่เหลือใน sales round detail
	salesRoundDetail.ProductStock = product.Stock
	salesRoundDetail.Remaining = product.Stock

	log.Printf("Creating new sales round detail: %v", salesRoundDetail)
	// สร้าง sales round detail
	if err := r.db.Create(salesRoundDetail).Error; err != nil {
		return nil, err
	}

	// เพิ่มการตรวจสอบผลการ preload ข้อมูล SalesRound และ ProductVariant
	log.Println("Preloading SalesRound and ProductVariant for new SalesRoundDetail...")
	if err := r.db.Preload("SalesRound").
		Preload("ProductVariant").
		First(salesRoundDetail, salesRoundDetail.ID).Error; err != nil {
		log.Printf("Error reloading SalesRoundDetail: %v", err)
		return nil, err
	}

	// ตรวจสอบค่าที่ถูก preload เข้ามา
	log.Printf("Loaded ProductVariant: %+v", salesRoundDetail.ProductVariant)

	// สร้างและรีเทิร์น SalesRoundDetailResponseDTO
	response := &dtos.SalesRoundDetailResponseDTO{
		ID:            salesRoundDetail.ID,
		RoundID:       salesRoundDetail.RoundID,
		VariantID:     salesRoundDetail.VariantID,
		Quantity:      salesRoundDetail.Quantity,
		Remaining:     salesRoundDetail.Remaining,
		QuantityLimit: salesRoundDetail.QuantityLimit,
		ProductStock:  salesRoundDetail.ProductStock,
		CreatedAt:     salesRoundDetail.CreatedAt,
		UpdatedAt:     salesRoundDetail.UpdatedAt,
		SalesRound: dtos.SalesRoundInfoResponse{
			ID:        salesRoundDetail.SalesRound.ID,
			Name:      salesRoundDetail.SalesRound.Name,
			StartDate: salesRoundDetail.SalesRound.StartDate,
			EndDate:   salesRoundDetail.SalesRound.EndDate,
			CreatedAt: salesRoundDetail.SalesRound.CreatedAt,
			UpdatedAt: salesRoundDetail.SalesRound.UpdatedAt,
		},
		ProductVariant: dtos.ProductVariantInfoResponse{
			CreatedAt: salesRoundDetail.ProductVariant.CreatedAt,
			UpdatedAt: salesRoundDetail.ProductVariant.UpdatedAt,
			ID:        salesRoundDetail.ProductVariant.VariantID,
			ProductID: salesRoundDetail.ProductVariant.ProductID,
			SKUCode:   salesRoundDetail.ProductVariant.SKUCode,
			// VariantID: salesRoundDetail.ProductVariant.VariantID,
			Price:    salesRoundDetail.ProductVariant.Price,
			ImageURL: salesRoundDetail.ProductVariant.ImageURL,
		},
	}

	return response, nil
}

func (r *salesRoundDetailRepository) GetAllSalesRoundDetails(expand string) ([]models.SalesRoundDetail, error) {
	var salesRoundDetails []models.SalesRoundDetail

	query := r.db.Model(&models.SalesRoundDetail{})

	// Always preload ProductVariant and its associated Product
	query = query.Preload("ProductVariant", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped() // This will include soft-deleted records
	}).Preload("ProductVariant.Product")

	// Preload related entities based on the expand parameter
	if strings.Contains(expand, "sales-rounds") {
		query = query.Preload("SalesRound")
	}
	if strings.Contains(expand, "categories") {
		query = query.Preload("ProductVariant.Category")
		query = query.Preload("ProductVariant.Product.Category")
	}
	if strings.Contains(expand, "stores") {
		query = query.Preload("ProductVariant.Store")
		query = query.Preload("ProductVariant.Product.Category.Store")
	}

	err := query.Find(&salesRoundDetails).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching sales round details: %w", err)
	}

	// Check if data was actually loaded
	for i, detail := range salesRoundDetails {
		if detail.ProductVariant.ID == uuid.Nil {
			log.Printf("Warning: ProductVariant not loaded for SalesRoundDetail ID: %s", detail.ID)
		}
		if detail.SalesRound.ID == uuid.Nil {
			log.Printf("Warning: SalesRound not loaded for SalesRoundDetail ID: %s", detail.ID)
		}
		// Add similar checks for other preloaded relations

		// Manually fetch ProductVariant if it's not loaded
		if detail.ProductVariant.ID == uuid.Nil {
			var productVariant models.ProductVariant
			if err := r.db.Unscoped().First(&productVariant, "variant_id = ?", detail.VariantID).Error; err != nil {
				log.Printf("Error fetching ProductVariant for SalesRoundDetail ID %s: %v", detail.ID, err)
			} else {
				salesRoundDetails[i].ProductVariant = productVariant
			}
		}

		// Manually fetch SalesRound if it's not loaded
		if detail.SalesRound.ID == uuid.Nil {
			var salesRound models.SalesRound
			if err := r.db.First(&salesRound, "id = ?", detail.RoundID).Error; err != nil {
				log.Printf("Error fetching SalesRound for SalesRoundDetail ID %s: %v", detail.ID, err)
			} else {
				salesRoundDetails[i].SalesRound = salesRound
			}
		}

		// Add similar manual fetching for other relations if needed
	}

	return salesRoundDetails, nil
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
		Select("\"sales-round\".*, \"sales-round-detail\".*, \"product-variant\".sku_code, \"product-variant\".price AS variant_price, \"product-variant\".image_url AS variant_image_url, \"product\".product_name, \"product\".brand, \"product\".description, \"product\".currency, \"product\".stock, \"product\".price AS product_price").
		Joins("JOIN \"sales-round\" ON \"sales-round-detail\".round_id = \"sales-round\".id").
		Joins("JOIN \"product-variant\" ON \"sales-round-detail\".variant_id = \"product-variant\".variant_id").
		Joins("JOIN \"product\" ON \"product-variant\".product_id = \"product\".id").
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

func (r *salesRoundDetailRepository) CheckPurchaseLimit(roundID uuid.UUID, variantID uuid.UUID, quantity int) (bool, error) {
	var detail models.SalesRoundDetail
	err := r.db.Where("round_id = ? AND variant_id = ?", roundID, variantID).First(&detail).Error
	if err != nil {
		return false, err
	}
	return detail.Quantity >= quantity, nil
}
