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

type OrderDetailController interface {
	CreateOrderDetail(c *fiber.Ctx) error
	GetAllOrderDetails(c *fiber.Ctx) error
	UpdateOrderDetail(c *fiber.Ctx) error
	DeleteOrderDetail(c *fiber.Ctx) error
}

type orderDetailController struct {
	orderDetailRepository repositories.OrderDetailRepository
}

func NewOrderDetailController(orderDetailRepository repositories.OrderDetailRepository) OrderDetailController {
	return &orderDetailController{orderDetailRepository: orderDetailRepository}
}

// CreateOrderDetail godoc
// @Summary Create a new order detail
// @Description Create a new order detail
// @Tags OrderDetails
// @Accept json
// @Produce json
// @Param orderDetail body dtos.OrderDetailCreateDTO true "Order Detail"
// @Success 201 {object} dtos.OrderDetailResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /order-details [post]
func (h *orderDetailController) CreateOrderDetail(c *fiber.Ctx) error {
	dto := new(dtos.OrderDetailCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	orderDetail := models.OrderDetail{
		OrderID:    dto.OrderID,
		VariantID:  dto.VariantID,
		Quantity:   dto.Quantity,
		Price:      dto.Price,
		TotalPrice: dto.TotalPrice,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.orderDetailRepository.CreateOrderDetail(&orderDetail)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create order detail"})
	}

	response := dtos.OrderDetailResponseDTO{
		ID:         orderDetail.ID,
		OrderID:    orderDetail.OrderID,
		VariantID:  orderDetail.VariantID,
		Quantity:   orderDetail.Quantity,
		Price:      orderDetail.Price,
		TotalPrice: orderDetail.TotalPrice,
		CreatedAt:  orderDetail.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  orderDetail.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllOrderDetails godoc
// @Summary Get all order details
// @Description Get all order details
// @Tags OrderDetails
// @Accept json
// @Produce json
// @Success 200 {array} dtos.OrderDetailResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /order-details [get]
func (h *orderDetailController) GetAllOrderDetails(c *fiber.Ctx) error {
	var orderDetails []models.OrderDetail

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		orderDetails, err = h.orderDetailRepository.GetAllOrderDetails()
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve order details"})
	}

	var orderDetailResponses []dtos.OrderDetailResponseDTO
	for _, orderDetail := range orderDetails {
		orderDetailResponses = append(orderDetailResponses, dtos.OrderDetailResponseDTO{
			ID:         orderDetail.ID,
			OrderID:    orderDetail.OrderID,
			VariantID:  orderDetail.VariantID,
			Quantity:   orderDetail.Quantity,
			Price:      orderDetail.Price,
			TotalPrice: orderDetail.TotalPrice,
			CreatedAt:  orderDetail.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  orderDetail.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(orderDetailResponses)
}

// UpdateOrderDetail godoc
// @Summary Update an order detail
// @Description Update an order detail
// @Tags OrderDetails
// @Accept json
// @Produce json
// @Param id path string true "Order Detail ID"
// @Param orderDetail body dtos.OrderDetailUpdateDTO true "Order Detail"
// @Success 200 {object} dtos.OrderDetailResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /order-details/{id} [put]
func (h *orderDetailController) UpdateOrderDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	detailUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.OrderDetailUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	// Initialize error channel and wait group
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	var orderDetail *models.OrderDetail
	var fetchErr error

	// Adding a goroutine to fetch the order detail
	wg.Add(1)
	go func() {
		defer wg.Done()
		orderDetail, fetchErr = h.orderDetailRepository.GetOrderDetailByID(detailUUID)
		if fetchErr != nil {
			errChan <- fetchErr
		}
	}()

	// Wait for the goroutine to finish
	wg.Wait()

	// Handle fetch error after the wait group is done
	if fetchErr != nil {
		close(errChan) // Safe to close the channel as it is read only after the goroutine is done
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order detail not found", "details": fetchErr.Error()})
	}

	// Proceed to update the order detail if the fetch was successful
	orderDetail.Quantity = dto.Quantity
	orderDetail.Price = dto.Price
	orderDetail.TotalPrice = dto.TotalPrice

	// Reset the wait group for the update operation
	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- h.orderDetailRepository.UpdateOrderDetail(orderDetail)
	}()

	// Wait for the update goroutine to finish
	wg.Wait()

	// Handle the update error
	if updateErr := <-errChan; updateErr != nil {
		close(errChan) // Close the channel as no more errors will be sent
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update order detail"})
	}

	close(errChan) // Ensure to close the channel safely after all operations are complete

	response := dtos.OrderDetailResponseDTO{
		ID:         orderDetail.ID,
		OrderID:    orderDetail.OrderID,
		VariantID:  orderDetail.VariantID,
		Quantity:   orderDetail.Quantity,
		Price:      orderDetail.Price,
		TotalPrice: orderDetail.TotalPrice,
		CreatedAt:  orderDetail.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  orderDetail.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.JSON(response)
}

// DeleteOrderDetail godoc
// @Summary Delete an order detail
// @Description Delete an order detail
// @Tags OrderDetails
// @Param id path string true "Order Detail ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /order-details/{id} [delete]
func (h *orderDetailController) DeleteOrderDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	detailUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.orderDetailRepository.DeleteOrderDetail(detailUUID)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete order detail"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
