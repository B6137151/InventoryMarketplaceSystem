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

type OrderController interface {
	CreateOrder(c *fiber.Ctx) error
	GetAllOrders(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
	DeleteOrder(c *fiber.Ctx) error
}

type orderController struct {
	orderRepository repositories.OrderRepository
}

func NewOrderController(orderRepository repositories.OrderRepository) OrderController {
	return &orderController{orderRepository: orderRepository}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body dtos.OrderCreateDTO true "Order"
// @Success 201 {object} dtos.OrderResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /orders [post]
func (h *orderController) CreateOrder(c *fiber.Ctx) error {
	dto := new(dtos.OrderCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	order := models.Order{
		CustomerID:      dto.CustomerID,
		RoundID:         dto.RoundID,
		OrderDate:       dto.OrderDate,
		Status:          dto.Status,
		Code:            dto.Code,
		TotalPrice:      dto.TotalPrice,
		DeliveryAddress: dto.DeliveryAddress,
		PaymentSource:   dto.PaymentSource,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.orderRepository.CreateOrder(&order)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create order"})
	}

	response := dtos.OrderResponseDTO{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		RoundID:         order.RoundID,
		OrderDate:       order.OrderDate,
		Status:          order.Status,
		Code:            order.Code,
		TotalPrice:      order.TotalPrice,
		DeliveryAddress: order.DeliveryAddress,
		PaymentSource:   order.PaymentSource,
		CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Get all orders
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} dtos.OrderResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /orders [get]
func (h *orderController) GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		orders, err = h.orderRepository.GetAllOrders()
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve orders"})
	}

	var orderResponses []dtos.OrderResponseDTO
	for _, order := range orders {
		orderResponses = append(orderResponses, dtos.OrderResponseDTO{
			ID:              order.ID,
			CustomerID:      order.CustomerID,
			RoundID:         order.RoundID,
			OrderDate:       order.OrderDate,
			Status:          order.Status,
			Code:            order.Code,
			TotalPrice:      order.TotalPrice,
			DeliveryAddress: order.DeliveryAddress,
			PaymentSource:   order.PaymentSource,
			CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       order.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(orderResponses)
}

// UpdateOrder godoc
// @Summary Update an order
// @Description Update an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param order body dtos.OrderUpdateDTO true "Order"
// @Success 200 {object} dtos.OrderResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /orders/{id} [put]
func (h *orderController) UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.OrderUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	var order *models.Order

	// Use a goroutine to fetch order with synchronization
	errChan := make(chan error, 1)
	done := make(chan bool)
	go func() {
		defer close(errChan)
		order, err = h.orderRepository.GetOrderByID(uuid)
		if err != nil {
			errChan <- err
			return
		}
		done <- true // signal completion of fetch
	}()

	// Wait for the fetch to complete or an error
	select {
	case <-done:
		// Continue with update if fetch was successful
		order.Status = dto.Status
		order.TotalPrice = dto.TotalPrice
		order.DeliveryAddress = dto.DeliveryAddress
		order.PaymentSource = dto.PaymentSource

		if updateErr := h.orderRepository.UpdateOrder(order); updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update order"})
		}
	case err := <-errChan:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order not found", "details": err.Error()})
	}

	response := dtos.OrderResponseDTO{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		RoundID:         order.RoundID,
		OrderDate:       order.OrderDate,
		Status:          order.Status,
		Code:            order.Code,
		TotalPrice:      order.TotalPrice,
		DeliveryAddress: order.DeliveryAddress,
		PaymentSource:   order.PaymentSource,
		CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.JSON(response)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order
// @Tags Orders
// @Param id path string true "Order ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /orders/{id} [delete]
func (h *orderController) DeleteOrder(c *fiber.Ctx) error {
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
		errChan <- h.orderRepository.DeleteOrder(uuid)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete order"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
