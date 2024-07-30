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
)

type SalesRoundController interface {
	CreateSalesRound(c *fiber.Ctx) error
	GetAllSalesRounds(c *fiber.Ctx) error
	GetSalesRoundDetails(c *fiber.Ctx) error
	UpdateSalesRound(c *fiber.Ctx) error
	DeleteSalesRound(c *fiber.Ctx) error
	GetCombinedSalesRoundProductData(c *fiber.Ctx) error // New method
}

type salesRoundController struct {
	salesRoundRepository       repositories.SalesRoundRepository
	orderRepository            repositories.OrderRepository
	salesRoundDetailRepository repositories.SalesRoundDetailRepository
}

func NewSalesRoundController(salesRoundRepository repositories.SalesRoundRepository, orderRepository repositories.OrderRepository, salesRoundDetailRepository repositories.SalesRoundDetailRepository) SalesRoundController {
	return &salesRoundController{
		salesRoundRepository:       salesRoundRepository,
		orderRepository:            orderRepository,
		salesRoundDetailRepository: salesRoundDetailRepository,
	}
}

// CreateSalesRound godoc
// @Summary Create a new sales round
// @Description Create a new sales round
// @Tags Sales Rounds
// @Accept json
// @Produce json
// @Param salesRound body dtos.SalesRoundCreateDTO true "Sales Round"
// @Success 201 {object} dtos.SalesRoundResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /sales-rounds [post]
func (c *salesRoundController) CreateSalesRound(ctx *fiber.Ctx) error {
	dto := new(dtos.SalesRoundCreateDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	salesRound := models.SalesRound{
		Name:      dto.Name,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- c.salesRoundRepository.CreateSalesRound(&salesRound)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create sales round"})
	}

	response := dtos.SalesRoundResponseDTO{
		ID:        salesRound.ID,
		Name:      salesRound.Name,
		StartDate: salesRound.StartDate,
		EndDate:   salesRound.EndDate,
		CreatedAt: salesRound.CreatedAt,
		UpdatedAt: salesRound.UpdatedAt,
	}
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

// GetAllSalesRounds godoc
// @Summary Get all sales rounds
// @Description Get all sales rounds
// @Tags Sales Rounds
// @Accept json
// @Produce json
// @Success 200 {array} dtos.SalesRoundResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /sales-rounds [get]
func (c *salesRoundController) GetAllSalesRounds(ctx *fiber.Ctx) error {
	salesRounds, err := c.salesRoundRepository.GetAllSalesRounds()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve sales rounds"})
	}

	var responses []dtos.SalesRoundResponseDTO
	for _, round := range salesRounds {
		responses = append(responses, dtos.SalesRoundResponseDTO{
			ID:        round.ID,
			Name:      round.Name,
			StartDate: round.StartDate,
			EndDate:   round.EndDate,
			CreatedAt: round.CreatedAt,
			UpdatedAt: round.UpdatedAt,
		})
	}
	return ctx.JSON(responses)
}

// GetSalesRoundDetails godoc
// @Summary Get details of a sales round
// @Description Get details of a sales round
// @Tags Sales Rounds
// @Accept json
// @Produce json
// @Param id path string true "Sales Round ID"
// @Success 200 {object} models.SalesRoundDetail
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /sales-rounds/{id}/details [get]
func (c *salesRoundController) GetSalesRoundDetails(ctx *fiber.Ctx) error {
	roundID := ctx.Params("id")
	id, err := uuid.Parse(roundID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid round ID"})
	}

	details, err := c.salesRoundDetailRepository.GetSalesRoundDetailsByRoundID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve sales round details"})
	}

	return ctx.JSON(details)
}

// UpdateSalesRound godoc
// @Summary Update a sales round
// @Description Update a sales round
// @Tags Sales Rounds
// @Accept json
// @Produce json
// @Param id path string true "Sales Round ID"
// @Param salesRound body dtos.SalesRoundUpdateDTO true "Sales Round"
// @Success 200 {object} dtos.SalesRoundResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /sales-rounds/{id} [put]
func (c *salesRoundController) UpdateSalesRound(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.SalesRoundUpdateDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var salesRound *models.SalesRound

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		salesRound, err = c.salesRoundRepository.GetSalesRoundByID(uuidID)
		if err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	wg.Wait()

	if err := <-errChan; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "sales round not found"})
	}

	salesRound.Name = dto.Name
	salesRound.StartDate = dto.StartDate
	salesRound.EndDate = dto.EndDate

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- c.salesRoundRepository.UpdateSalesRound(salesRound)
	}()

	wg.Wait()
	close(errChan)

	if updateErr := <-errChan; updateErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update sales round"})
	}

	response := dtos.SalesRoundResponseDTO{
		ID:        salesRound.ID,
		Name:      salesRound.Name,
		StartDate: salesRound.StartDate,
		EndDate:   salesRound.EndDate,
		CreatedAt: salesRound.CreatedAt,
		UpdatedAt: salesRound.UpdatedAt,
	}
	return ctx.JSON(response)
}

// DeleteSalesRound godoc
// @Summary Delete a sales round
// @Description Delete a sales round
// @Tags Sales Rounds
// @Param id path string true "Sales Round ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /sales-rounds/{id} [delete]
func (c *salesRoundController) DeleteSalesRound(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	err = c.salesRoundRepository.DeleteSalesRound(uuidID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete sales round"})
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

// GetCombinedSalesRoundProductData godoc
// @Summary Get combined sales round and product data
// @Description Get combined sales round and product data
// @Tags Sales Rounds
// @Accept json
// @Produce json
// @Success 200 {object} []models.CombinedSalesRoundProductData
// @Failure 500 {object} fiber.Map
// @Router /sales-rounds/combined-data [get]
func (h *salesRoundController) GetCombinedSalesRoundProductData(c *fiber.Ctx) error {
	log.Println("Fetching combined sales round product data")
	data, err := h.salesRoundRepository.GetCombinedSalesRoundProductData()
	if err != nil {
		log.Println("Error fetching data:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not fetch data"})
	}
	log.Println("Successfully fetched data:", data)
	return c.JSON(data)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
