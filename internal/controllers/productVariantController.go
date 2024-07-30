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

type ProductVariantController interface {
	CreateProductVariant(c *fiber.Ctx) error
	GetAllProductVariants(c *fiber.Ctx) error
	UpdateProductVariant(c *fiber.Ctx) error
	DeleteProductVariant(c *fiber.Ctx) error
}

type productVariantController struct {
	productVariantRepository repositories.ProductVariantRepository
}

func NewProductVariantController(productVariantRepository repositories.ProductVariantRepository) ProductVariantController {
	return &productVariantController{productVariantRepository: productVariantRepository}
}

// CreateProductVariant godoc
// @Summary Create a new product variant
// @Description Create a new product variant
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param productVariant body dtos.ProductVariantCreateDTO true "Product Variant"
// @Success 201 {object} dtos.ProductVariantResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /product-variants [post]
func (h *productVariantController) CreateProductVariant(c *fiber.Ctx) error {
	dto := new(dtos.ProductVariantCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	productVariant := models.ProductVariant{
		ProductID: dto.ProductID,
		SKUCode:   dto.SKUCode,
		Price:     dto.Price,
		ImageURL:  dto.ImageURL,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.productVariantRepository.CreateProductVariant(&productVariant)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create product variant"})
	}

	response := dtos.ProductVariantResponseDTO{
		ID:        productVariant.VariantID,
		ProductID: productVariant.ProductID,
		SKUCode:   productVariant.SKUCode,
		Price:     productVariant.Price,
		ImageURL:  productVariant.ImageURL,
		CreatedAt: productVariant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: productVariant.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllProductVariants godoc
// @Summary Get all product variants
// @Description Get all product variants
// @Tags Product Variants
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ProductVariantResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /product-variants [get]
func (h *productVariantController) GetAllProductVariants(c *fiber.Ctx) error {
	var productVariants []models.ProductVariant

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		productVariants, err = h.productVariantRepository.GetAllProductVariants()
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve product variants"})
	}

	var productVariantResponses []dtos.ProductVariantResponseDTO
	for _, productVariant := range productVariants {
		productVariantResponses = append(productVariantResponses, dtos.ProductVariantResponseDTO{
			ID:        productVariant.VariantID,
			ProductID: productVariant.ProductID,
			SKUCode:   productVariant.SKUCode,
			Price:     productVariant.Price,
			ImageURL:  productVariant.ImageURL,
			CreatedAt: productVariant.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: productVariant.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(productVariantResponses)
}

// UpdateProductVariant godoc
// @Summary Update a product variant
// @Description Update a product variant
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param id path string true "Product Variant ID"
// @Param productVariant body dtos.ProductVariantUpdateDTO true "Product Variant"
// @Success 200 {object} dtos.ProductVariantResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /product-variants/{id} [put]
func (h *productVariantController) UpdateProductVariant(c *fiber.Ctx) error {
	id := c.Params("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.ProductVariantUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	// Initialize the wait group and error channel for thread-safe operations
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	var productVariant *models.ProductVariant

	// Fetch the product variant in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		productVariant, err = h.productVariantRepository.GetProductVariantByID(uuidID)
		if err != nil {
			errChan <- err
		}
	}()

	// Wait for the fetch operation to complete
	wg.Wait()

	// Check for errors from the fetch operation
	if err := <-errChan; err != nil {
		close(errChan) // Close the channel safely after reading the error
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product variant not found"})
	}

	// Update the product variant with new data from the DTO
	productVariant.ProductID = dto.ProductID
	productVariant.SKUCode = dto.SKUCode
	productVariant.Price = dto.Price
	productVariant.ImageURL = dto.ImageURL

	// Reset the wait group and reuse the channel for the update operation
	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- h.productVariantRepository.UpdateProductVariant(productVariant)
	}()

	// Wait for the update operation to complete
	wg.Wait()

	// Handle any errors from the update operation
	if updateErr := <-errChan; updateErr != nil {
		close(errChan) // Ensure to close the channel after all operations are done
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update product variant"})
	}

	close(errChan) // Close the channel as a cleanup action

	// Prepare the response object
	response := dtos.ProductVariantResponseDTO{
		ID:        productVariant.VariantID,
		ProductID: productVariant.ProductID,
		SKUCode:   productVariant.SKUCode,
		Price:     productVariant.Price,
		ImageURL:  productVariant.ImageURL,
		CreatedAt: productVariant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: productVariant.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.JSON(response)
}

// DeleteProductVariant godoc
// @Summary Delete a product variant
// @Description Delete a product variant
// @Tags Product Variants
// @Param id path string true "Product Variant ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /product-variants/{id} [delete]
func (h *productVariantController) DeleteProductVariant(c *fiber.Ctx) error {
	id := c.Params("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.productVariantRepository.DeleteProductVariant(uuidID)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete product variant"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
