package controllers

import (
	"runtime"
	"sync"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type OrderHistoryController interface {
	CreateOrderHistory(c *fiber.Ctx) error
	GetAllOrderHistories(c *fiber.Ctx) error
	UpdateOrderHistory(c *fiber.Ctx) error
	DeleteOrderHistory(c *fiber.Ctx) error
}

type orderHistoryController struct {
	orderHistoryRepository repositories.OrderHistoryRepository
}

func NewOrderHistoryController(orderHistoryRepository repositories.OrderHistoryRepository) OrderHistoryController {
	return &orderHistoryController{orderHistoryRepository: orderHistoryRepository}
}

// CreateOrderHistory godoc
// @Summary Create a new order history
// @Description Create a new order history
// @Tags OrderHistories
// @Accept json
// @Produce json
// @Param orderHistory body dtos.OrderHistoryCreateDTO true "Order History"
// @Success 201 {object} dtos.OrderHistoryResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /order-histories [post]
func (h *orderHistoryController) CreateOrderHistory(c *fiber.Ctx) error {
	dto := new(dtos.OrderHistoryCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	orderHistory := models.OrderHistory{
		OrderID:     dto.OrderID,
		Status:      dto.Status,
		ChangedAt:   dto.ChangedAt,
		Description: dto.Description,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.orderHistoryRepository.CreateOrderHistory(&orderHistory)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create order history"})
	}

	response := dtos.OrderHistoryResponseDTO{
		ID:          orderHistory.ID,
		OrderID:     orderHistory.OrderID,
		Status:      orderHistory.Status,
		ChangedAt:   orderHistory.ChangedAt,
		Description: orderHistory.Description,
		CreatedAt:   orderHistory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   orderHistory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllOrderHistories godoc
// @Summary Get all order histories
// @Description Get all order histories
// @Tags OrderHistories
// @Accept json
// @Produce json
// @Success 200 {array} dtos.OrderHistoryResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /order-histories [get]
func (h *orderHistoryController) GetAllOrderHistories(c *fiber.Ctx) error {
	var orderHistories []models.OrderHistory

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		orderHistories, err = h.orderHistoryRepository.GetAllOrderHistories()
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve order histories"})
	}

	var orderHistoryResponses []dtos.OrderHistoryResponseDTO
	for _, orderHistory := range orderHistories {
		orderHistoryResponses = append(orderHistoryResponses, dtos.OrderHistoryResponseDTO{
			ID:          orderHistory.ID,
			OrderID:     orderHistory.OrderID,
			Status:      orderHistory.Status,
			ChangedAt:   orderHistory.ChangedAt,
			Description: orderHistory.Description,
			CreatedAt:   orderHistory.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   orderHistory.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(orderHistoryResponses)
}

// UpdateOrderHistory godoc
// @Summary Update an order history
// @Description Update an order history
// @Tags OrderHistories
// @Accept json
// @Produce json
// @Param id path string true "Order History ID"
// @Param orderHistory body dtos.OrderHistoryUpdateDTO true "Order History"
// @Success 200 {object} dtos.OrderHistoryResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /order-histories/{id} [put]
func (h *orderHistoryController) UpdateOrderHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	dto := new(dtos.OrderHistoryUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	// Initialize the error channel and wait group for synchronized goroutine management
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	var orderHistory *models.OrderHistory
	var fetchErr error

	// Start a goroutine to fetch the order history
	wg.Add(1)
	go func() {
		defer wg.Done()
		orderHistory, fetchErr = h.orderHistoryRepository.GetOrderHistoryByID(id)
		if fetchErr != nil {
			errChan <- fetchErr
		}
	}()

	// Wait for the fetch goroutine to complete
	wg.Wait()

	// Check for fetch errors
	if fetchErr != nil {
		close(errChan) // Safe to close here as no more writes to errChan
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order history not found", "details": fetchErr.Error()})
	}

	// Proceed to update the order history details
	orderHistory.Status = dto.Status
	orderHistory.ChangedAt = dto.ChangedAt
	orderHistory.Description = dto.Description

	// Reset the wait group for the update operation
	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- h.orderHistoryRepository.UpdateOrderHistory(orderHistory)
	}()

	// Wait for the update operation to complete
	wg.Wait()

	// Handle potential update errors
	if updateErr := <-errChan; updateErr != nil {
		close(errChan) // Close the channel as no more writes will occur
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update order history"})
	}

	close(errChan) // Ensure the channel is closed safely after all operations are complete

	response := dtos.OrderHistoryResponseDTO{
		ID:          orderHistory.ID,
		OrderID:     orderHistory.OrderID,
		Status:      orderHistory.Status,
		ChangedAt:   orderHistory.ChangedAt,
		Description: orderHistory.Description,
		CreatedAt:   orderHistory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   orderHistory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.JSON(response)
}

// DeleteOrderHistory godoc
// @Summary Delete an order history
// @Description Delete an order history
// @Tags OrderHistories
// @Param id path string true "Order History ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /order-histories/{id} [delete]
func (h *orderHistoryController) DeleteOrderHistory(c *fiber.Ctx) error {
	id := c.Params("id")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.orderHistoryRepository.DeleteOrderHistory(id)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete order history"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
