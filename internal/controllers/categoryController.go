package controllers

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryController interface {
	CreateCategory(c *fiber.Ctx) error
	GetAllCategories(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type categoryController struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryController(categoryRepository repositories.CategoryRepository) CategoryController {
	return &categoryController{categoryRepository: categoryRepository}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body dtos.CategoryCreateDTO true "Category"
// @Success 201 {object} dtos.CategoryResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /categories [post]
func (h *categoryController) CreateCategory(c *fiber.Ctx) error {
	dto := new(dtos.CategoryCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	category := models.Category{
		Name:    dto.Name,
		StoreID: dto.StoreID,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- h.categoryRepository.CreateCategory(&category)
	}()

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create category"})
	}

	response := dtos.CategoryResponseDTO{
		ID:        category.ID,
		Name:      category.Name,
		StoreID:   category.StoreID,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {array} dtos.CategoryResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /categories [get]
func (h *categoryController) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.categoryRepository.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve categories"})
	}

	var categoryResponses []dtos.CategoryResponseDTO
	for _, category := range categories {
		categoryResponses = append(categoryResponses, dtos.CategoryResponseDTO{
			ID:        category.ID,
			Name:      category.Name,
			StoreID:   category.StoreID,
			CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(categoryResponses)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body dtos.CategoryUpdateDTO true "Category"
// @Success 200 {object} dtos.CategoryResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /categories/{id} [put]
func (h *categoryController) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.CategoryUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	var category *models.Category

	errChan := make(chan error, 1)
	done := make(chan bool)
	go func() {
		defer close(errChan)
		category, err = h.categoryRepository.GetCategoryByID(categoryID)
		if err != nil {
			errChan <- err
			return
		}
		done <- true
	}()

	select {
	case <-done:
		category.Name = dto.Name
		category.StoreID = dto.StoreID

		if updateErr := h.categoryRepository.UpdateCategory(category); updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update category"})
		}
	case err := <-errChan:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "category not found", "details": err.Error()})
	}

	response := dtos.CategoryResponseDTO{
		ID:        category.ID,
		Name:      category.Name,
		StoreID:   category.StoreID,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.JSON(response)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category
// @Tags Categories
// @Param id path string true "Category ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /categories/{id} [delete]
func (h *categoryController) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		errChan <- h.categoryRepository.DeleteCategory(categoryID)
	}()

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete category"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
