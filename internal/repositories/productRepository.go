package repositories

import (
	"log"

	"strings"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetAllProducts(expand string) ([]models.Product, error)
	GetProductByID(id uuid.UUID) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uuid.UUID) error
	GetAllProductsWithVariants() ([]models.Product, error) // New method
	GetProductsBySalesRoundID(salesRoundID uuid.UUID) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	log.Println("Creating product:", product)
	err := r.db.Create(product).Error
	if err != nil {
		log.Println("Error creating product:", err)
	}
	return err
}

func (r *productRepository) GetAllProducts(expand string) ([]models.Product, error) {
	var products []models.Product

	log.Println("Starting GetAllProducts...")

	query := r.db.Model(&models.Product{})
	log.Println("Initial query prepared")

	if expand != "" {
		log.Printf("Expand parameter found: %s\n", expand)

		if strings.Contains(expand, "product-varian") {
			log.Println("Preloading ProductVariants...")
			query = query.Preload("product-variant")
		}

		if strings.Contains(expand, "store") {
			log.Println("Preloading Store...")
			query = query.Preload("Store")
		}

		if strings.Contains(expand, "category") {
			log.Println("Preloading Category...")
			query = query.Preload("Category")
		}
	} else {
		log.Println("No expand parameter provided")
	}

	err := query.Find(&products).Error
	if err != nil {
		log.Printf("Error retrieving products: %v\n", err)
		return nil, err
	}

	log.Printf("Successfully retrieved %d products\n", len(products))
	return products, nil
}

func (r *productRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	log.Println("Fetching product by ID:", id)
	err := r.db.Preload("Category").Preload("ProductVariant").First(&product, "id = ?", id).Error
	if err != nil {
		log.Println("Error fetching product:", err)
	}
	log.Println("Fetched product:", product)
	return &product, err
}

func (r *productRepository) UpdateProduct(product *models.Product) error {
	log.Println("Updating product:", product)
	err := r.db.Save(product).Error
	if err != nil {
		log.Println("Error updating product:", err)
	}
	return err
}

func (r *productRepository) DeleteProduct(id uuid.UUID) error {
	log.Println("Deleting product by ID:", id)
	err := r.db.Delete(&models.Product{}, "id = ?", id).Error
	if err != nil {
		log.Println("Error deleting product:", err)
	}
	return err
}

func (r *productRepository) GetAllProductsWithVariants() ([]models.Product, error) {
	var products []models.Product
	log.Println("Fetching all products with variants")
	err := r.db.Preload("ProductVariant").Find(&products).Error
	if err != nil {
		log.Println("Error fetching products with variants:", err)
	}
	log.Println("Fetched products with variants:", products)
	return products, err
}
func (r *productRepository) GetProductsBySalesRoundID(salesRoundID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Joins(`JOIN "sales-round-detail" ON "sales-round-detail".product_id = product.id`).
		Where(`"sales-round-detail".round_id = ?`, salesRoundID).
		Find(&products).Error
	if err != nil {
		log.Println("Error fetching products by sales round ID:", err)
	}
	log.Println("Fetched products by sales round ID:", products)
	return products, err
}

func (r *productRepository) GetProductVariantsBySalesRoundID(salesRoundID uuid.UUID) ([]models.ProductVariant, error) {
	var productVariants []models.ProductVariant
	err := r.db.Joins(`JOIN "sales-round-detail" ON "sales-round-detail".variant_id = "product-variant".id`).
		Where(`"sales-round-detail".round_id = ?`, salesRoundID).
		Find(&productVariants).Error
	if err != nil {
		log.Println("Error fetching product variants by sales round ID:", err)
	}
	log.Println("Fetched product variants by sales round ID:", productVariants)
	return productVariants, err
}
