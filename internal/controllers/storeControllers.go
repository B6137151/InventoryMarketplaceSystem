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

type StoreController interface {
	CreateStore(c *fiber.Ctx) error
	GetAllStores(c *fiber.Ctx) error
	GetStoreByID(c *fiber.Ctx) error
	UpdateStore(c *fiber.Ctx) error
	DeleteStore(c *fiber.Ctx) error
}

type storeController struct {
	storeRepository repositories.StoreRepository
}

func NewStoreController(storeRepository repositories.StoreRepository) StoreController {
	return &storeController{storeRepository: storeRepository}
}

// CreateStore godoc
// @Summary Create a new store
// @Description Create a new store
// @Tags Stores
// @Accept json
// @Produce json
// @Param store body dtos.StoreCreateDTO true "Store"
// @Success 201 {object} dtos.StoreResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /stores [post]
func (h *storeController) CreateStore(c *fiber.Ctx) error {
	dto := new(dtos.StoreCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	store := models.Store{StoreName: dto.StoreName, Location: dto.Location}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := h.storeRepository.CreateStore(&store); err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create store"})
	}

	response := dtos.StoreResponseDTO{
		ID:        store.ID,
		StoreName: store.StoreName,
		Location:  store.Location,
		CreatedAt: store.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllStores godoc
// @Summary Get all stores
// @Description Get all stores
// @Tags Stores
// @Accept json
// @Produce json
// @Success 200 {array} models.Store
// @Failure 500 {object} fiber.Map
// @Router /stores [get]
func (h *storeController) GetAllStores(c *fiber.Ctx) error {
	var stores []models.Store

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		stores, err = h.storeRepository.GetAllStores()
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve stores"})
	}
	return c.JSON(stores)
}

// GetStoreByID godoc
// @Summary Get store by ID
// @Description Get store by ID
// @Tags Stores
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} models.Store
// @Failure 500 {object} fiber.Map
// @Router /stores/{id} [get]
func (h *storeController) GetStoreByID(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	var store *models.Store
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		store, err = h.storeRepository.GetStoreByID(uuid)
		if err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Store not found"})
	}

	return c.JSON(store)
}

// UpdateStore godoc
// @Summary Update a store
// @Description Update a store
// @Tags Stores
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param store body dtos.StoreUpdateDTO true "Store"
// @Success 200 {object} dtos.StoreResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /stores/{id} [put]
func (h *storeController) UpdateStore(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.StoreUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	var store *models.Store
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// Fetch the existing store
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		store, err = h.storeRepository.GetStoreByID(uuid)
		if err != nil {
			errChan <- err
		}
	}()

	// Wait for the fetch to complete before attempting to close the channel
	wg.Wait()

	// Handle potential fetch error

	select {
	case err := <-errChan:
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "store not found"})
		}
	default:
		// No error was sent; continue
	}

	// Proceed to update the store details if found
	store.StoreName = dto.StoreName
	store.Location = dto.Location

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := h.storeRepository.UpdateStore(store); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update store"})
		}
	default:

	}

	response := dtos.StoreResponseDTO{
		ID:        store.ID,
		StoreName: store.StoreName,
		Location:  store.Location,
		CreatedAt: store.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.JSON(response)
}

// DeleteStore godoc
// @Summary Delete a store
// @Description Delete a store
// @Tags Stores
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Success 204 {object} nil
// @Failure 500 {object} fiber.Map
// @Router /stores/{id} [delete]
func (h *storeController) DeleteStore(c *fiber.Ctx) error {
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
		if err := h.storeRepository.DeleteStore(uuid); err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete store"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
