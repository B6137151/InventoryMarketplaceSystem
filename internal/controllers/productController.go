package controllers

import (
	"runtime"
	"sync"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductController interface {
	CreateProduct(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	GetAllProductsWithVariants(c *fiber.Ctx) error // New method
}

type productController struct {
	productRepository repositories.ProductRepository
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
	var products []models.Product

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		products, err = h.productRepository.GetAllProducts()
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve products"})
	}

	var productResponses []dtos.ProductResponseDTO
	for _, product := range products {
		productResponses = append(productResponses, dtos.ProductResponseDTO{
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
		})
	}

	return c.JSON(productResponses)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dtos.ProductUpdateDTO true "Product"
// @Success 200 {object} dtos.ProductResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /products/{id} [put]
func (h *productController) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.ProductUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	var product *models.Product

	// Start a goroutine to fetch the product
	wg.Add(1)
	go func() {
		defer wg.Done()
		var fetchErr error
		product, fetchErr = h.productRepository.GetProductByID(uuid)
		if fetchErr != nil {
			errChan <- fetchErr
		}
	}()

	// Wait for the fetch goroutine to finish
	wg.Wait()

	// Handle potential fetch error
	if err := <-errChan; err != nil {
		close(errChan) // Safe to close here, as no further writes to errChan
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	// Update product details if the fetch was successful
	product.StoreID = dto.StoreID
	product.CategoryID = dto.CategoryID
	product.ProductName = dto.ProductName
	product.Brand = dto.Brand
	product.Description = dto.Description
	product.Currency = dto.Currency
	product.Stock = dto.Stock
	product.Price = dto.Price
	product.ImageURL = dto.ImageURL

	// Reset the wait group for the update operation
	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- h.productRepository.UpdateProduct(product)
	}()

	// Wait for the update operation to complete
	wg.Wait()

	// Check for update errors and close the channel
	if updateErr := <-errChan; updateErr != nil {
		close(errChan) // Close the channel after reading the error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update product"})
	}

	close(errChan) // Ensure the channel is closed safely after all operations are complete

	// Prepare the response
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
	var wg sync.WaitGroup
	var products []models.Product
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

	return c.JSON(products)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
