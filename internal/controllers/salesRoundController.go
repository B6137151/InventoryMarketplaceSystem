package controllers

import (
	"log"
	"runtime"
	"strings"
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
	productRepository          repositories.ProductRepository
	productVariantRepository   repositories.ProductVariantRepository
	salesRoundDetailRepository repositories.SalesRoundDetailRepository
	categoryRepository         repositories.CategoryRepository
	storeRepository            repositories.StoreRepository
}

func NewSalesRoundController(
	salesRoundRepository repositories.SalesRoundRepository,
	productRepository repositories.ProductRepository,
	productVariantRepository repositories.ProductVariantRepository,
	salesRoundDetailRepository repositories.SalesRoundDetailRepository,
	categoryRepository repositories.CategoryRepository,
	storeRepository repositories.StoreRepository,
) SalesRoundController {
	return &salesRoundController{
		salesRoundRepository:       salesRoundRepository,
		productRepository:          productRepository,
		productVariantRepository:   productVariantRepository,
		salesRoundDetailRepository: salesRoundDetailRepository,
		categoryRepository:         categoryRepository,
		storeRepository:            storeRepository,
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
// @Description Get all sales rounds, optionally expanding related entities like products, product variants, sales round details, categories, and stores
// @Tags Sales Rounds
// @Accept json
// @Produce json
// @Param expand query string false "Fields to expand, e.g. expand=product,product-variant,sales-round-detail,category,store"
// @Success 200 {array} dtos.SalesRoundResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /sales-rounds [get]
func (c *salesRoundController) GetAllSalesRounds(ctx *fiber.Ctx) error {
	expand := ctx.Query("expand")

	// Fetch all sales rounds
	salesRounds, err := c.salesRoundRepository.GetAllSalesRounds()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve sales rounds"})
	}

	var responses []dtos.SalesRoundResponseDTO
	for _, round := range salesRounds {
		response := dtos.SalesRoundResponseDTO{
			ID:        round.ID,
			Name:      round.Name,
			StartDate: round.StartDate,
			EndDate:   round.EndDate,
			CreatedAt: round.CreatedAt, // Direct assignment
			UpdatedAt: round.UpdatedAt, // Direct assignment
		}

		if expand != "" {
			expansions := strings.Split(expand, ",")
			for _, exp := range expansions {
				switch strings.TrimSpace(exp) {
				case "product":
					products, err := c.productRepository.GetProductsBySalesRoundID(round.ID)
					if err != nil {
						return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve products"})
					}
					var productDTOs []dtos.ProductResponseDTO
					for _, product := range products {
						productDTOs = append(productDTOs, dtos.ProductResponseDTO{
							ID:          product.ID,
							ProductName: product.ProductName,
							Brand:       product.Brand,
							Description: product.Description,
							Price:       product.Price,
							ImageURL:    product.ImageURL,
							CategoryID:  product.CategoryID,
							StoreID:     product.StoreID,
						})
					}
					response.Products = productDTOs

				case "product-variant":
					productVariants, err := c.productVariantRepository.GetProductVariantsBySalesRoundID(round.ID)
					if err != nil {
						return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve product variants"})
					}
					var productVariantDTOs []dtos.ProductVariantResponseDTO
					for _, variant := range productVariants {
						productVariantDTOs = append(productVariantDTOs, dtos.ProductVariantResponseDTO{
							// ID:        variant.ID,
							ID:        variant.VariantID,
							ProductID: variant.ProductID,
							SKUCode:   variant.SKUCode,
							Price:     variant.Price,
							ImageURL:  variant.ImageURL,
						})
					}
					response.ProductVariants = productVariantDTOs

				case "sales-round-detail":
					salesRoundDetails, err := c.salesRoundDetailRepository.GetSalesRoundDetailsByRoundID(round.ID)
					if err != nil {
						return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve sales round details"})
					}
					var salesRoundDetailDTOs []dtos.SalesRoundDetailResponseDTO
					for _, detail := range salesRoundDetails {
						// Log the specific ID
						// log.Printf("Detail ID: %s", detail.ID) // Assuming `detail.ID` is the correct reference

						salesRoundDetailDTOs = append(salesRoundDetailDTOs, dtos.SalesRoundDetailResponseDTO{
							// ID:            detail.ID,
							RoundID:       detail.RoundID,
							VariantID:     detail.VariantID,
							Quantity:      detail.Quantity,
							Remaining:     detail.Remaining,
							ProductStock:  detail.ProductStock,
							QuantityLimit: detail.QuantityLimit,
						})
					}
					response.SalesRoundDetails = salesRoundDetailDTOs

				case "category":
					categories, err := c.categoryRepository.GetCategoriesBySalesRoundID(round.ID)
					if err != nil {
						return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve categories"})
					}
					var categoryDTOs []dtos.CategoryResponseDTO
					for _, category := range categories {
						categoryDTOs = append(categoryDTOs, dtos.CategoryResponseDTO{
							ID:   category.ID,
							Name: category.Name,
						})
					}
					response.Categories = categoryDTOs

				case "store":
					stores, err := c.storeRepository.GetStoresBySalesRoundID(round.ID)
					if err != nil {
						return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve stores"})
					}
					var storeDTOs []dtos.StoreResponseDTO
					for _, store := range stores {
						storeDTOs = append(storeDTOs, dtos.StoreResponseDTO{
							ID:        store.ID,
							StoreName: store.StoreName,
							Location:  store.Location,
						})
					}
					response.Stores = storeDTOs
				}
			}
		}

		responses = append(responses, response)
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
