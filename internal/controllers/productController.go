package controllers

import (
	"log"
	"runtime"

	"sync"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	// "gorm.io/gorm" // Import the gorm package
)

type ProductController interface {
	CreateProduct(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
	UpdateProductStock(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	GetAllProductsWithVariants(c *fiber.Ctx) error // New method
	GetProductByID(c *fiber.Ctx) error             // New method

}

type productController struct {
	productRepository repositories.ProductRepository
	// db                *gorm.DB
}

func NewProductController(productRepository repositories.ProductRepository) ProductController {
	return &productController{productRepository: productRepository}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param product body dtos.ProductCreateDTO true "Product"
// @Success 201 {object} dtos.ProductResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /products [post]
func (h *productController) CreateProduct(c *fiber.Ctx) error {
	dto := new(dtos.ProductCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	product := models.Product{
		StoreID:     dto.StoreID,
		CategoryID:  dto.CategoryID,
		ProductName: dto.ProductName,
		Brand:       dto.Brand,
		Description: dto.Description,
		Currency:    dto.Currency,
		Stock:       dto.Stock,
		Price:       dto.Price,
		ImageURL:    dto.ImageURL,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.productRepository.CreateProduct(&product)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create product"})
	}

	response := dtos.ProductResponseDTO{
		ID:          product.ID,
		StoreID:     product.StoreID,
		CategoryID:  product.CategoryID,
		ProductName: product.ProductName,
		Brand:       product.Brand,
		Description: product.Description,
		Currency:    product.Currency,
		Stock:       product.Stock,
		Price:       product.Price,
		ImageURL:    product.ImageURL,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ProductResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /products [get]
func (h *productController) GetAllProducts(c *fiber.Ctx) error {
	expand := c.Query("expand")

	// Fetch products using the repository method with expand parameter
	products, err := h.productRepository.GetAllProducts(expand)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve products"})
	}

	productResponses := make([]dtos.ProductResponseDTO, len(products))
	for i, product := range products {
		productResponses[i] = dtos.ProductResponseDTO{
			ID:           product.ID,
			StoreID:      product.StoreID,
			StoreName:    product.Store.StoreName,
			CategoryID:   product.CategoryID,
			CategoryName: product.Category.Name,
			ProductName:  product.ProductName,
			Brand:        product.Brand,
			Description:  product.Description,
			Currency:     product.Currency,
			Stock:        product.Stock,
			Price:        product.Price,
			ImageURL:     product.ImageURL,
			CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    product.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// Prepare response with metadata
	response := struct {
		Meta dtos.MetaData             `json:"meta"`
		Data []dtos.ProductResponseDTO `json:"data"`
	}{
		Meta: dtos.MetaData{
			Total: len(products),
			Count: len(productResponses),
		},
		Data: productResponses,
	}

	return c.JSON(response)
}

// UpdateProductStock godoc
// @Summary Update a product's stock
// @Description Update a product's stock
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dtos.ProductUpdateStockDTO true "Product Stock"
// @Success 200 {object} dtos.ProductResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /products/{id}/stock [put]
func (h *productController) UpdateProductStock(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Received request to update stock for product ID:", id)
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Println("Invalid UUID format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.ProductUpdateStockDTO)
	if err := c.BodyParser(dto); err != nil {
		log.Println("Failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	log.Println("Fetching product with ID:", uuid)
	product, err := h.productRepository.GetProductByID(uuid)
	if err != nil {
		log.Println("Product not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}
	log.Println("Fetched product:", product)

	log.Println("Updating product stock to:", dto.Stock)
	product.Stock = dto.Stock

	log.Println("Saving updated product to database")
	if err := h.productRepository.UpdateProduct(product); err != nil {
		log.Println("Failed to update product stock:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update product stock"})
	}
	log.Println("Product stock updated successfully")

	response := dtos.ProductResponseDTO{
		ID:          product.ID,
		StoreID:     product.StoreID,
		CategoryID:  product.CategoryID,
		ProductName: product.ProductName,
		Brand:       product.Brand,
		Description: product.Description,
		Currency:    product.Currency,
		Stock:       product.Stock,
		Price:       product.Price,
		ImageURL:    product.ImageURL,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	log.Println("Prepared response:", response)
	return c.JSON(response)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags Products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /products/{id} [delete]
func (h *productController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.productRepository.DeleteProduct(uuid)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete product"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetAllProductsWithVariants godoc
// @Summary Get all products with their variants
// @Description Get all products with their variants
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} fiber.Map
// @Router /products/variants [get]
func (h *productController) GetAllProductsWithVariants(c *fiber.Ctx) error {
	var products []models.Product
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		products, err = h.productRepository.GetAllProductsWithVariants()
		errChan <- err
	}()
	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve products with variants"})
	}

	// Transform products to response DTOs
	productResponses := make([]dtos.ProductResponseDTO, len(products))
	for i, product := range products {
		productResponses[i] = dtos.ProductResponseDTO{
			ID:           product.ID,
			StoreID:      product.StoreID,
			StoreName:    product.Store.StoreName,
			CategoryID:   product.CategoryID,
			CategoryName: product.Category.Name,
			ProductName:  product.ProductName,
			Brand:        product.Brand,
			Description:  product.Description,
			Currency:     product.Currency,
			Stock:        product.Stock,
			Price:        product.Price,
			ImageURL:     product.ImageURL,
			CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    product.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// Prepare response with metadata
	response := struct {
		Meta dtos.MetaData             `json:"meta"`
		Data []dtos.ProductResponseDTO `json:"data"`
	}{
		Meta: dtos.MetaData{
			Total: len(products),
			Count: len(productResponses),
		},
		Data: productResponses,
	}

	return c.JSON(response)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.ProductResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /products/{id} [get]
func (h *productController) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	product, err := h.productRepository.GetProductByID(uuid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	response := dtos.ProductResponseDTO{
		ID:           product.ID,
		StoreID:      product.StoreID,
		StoreName:    product.Store.StoreName,
		CategoryID:   product.CategoryID,
		CategoryName: product.Category.Name,
		ProductName:  product.ProductName,
		Brand:        product.Brand,
		Description:  product.Description,
		Currency:     product.Currency,
		Stock:        product.Stock,
		Price:        product.Price,
		ImageURL:     product.ImageURL,
		CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Prepare response with metadata
	responseWithMeta := struct {
		Meta dtos.MetaData           `json:"meta"`
		Data dtos.ProductResponseDTO `json:"data"`
	}{
		Meta: dtos.MetaData{
			Total: 1, // Since we're fetching a single product
			Count: 1,
		},
		Data: response,
	}

	return c.JSON(responseWithMeta)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
