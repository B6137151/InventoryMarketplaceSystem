package controllers

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SalesRoundDetailController interface {
	CreateSalesRoundDetail(c *fiber.Ctx) error
	GetAllSalesRoundDetails(c *fiber.Ctx) error
	UpdateSalesRoundDetail(c *fiber.Ctx) error
	DeleteSalesRoundDetail(c *fiber.Ctx) error
	GetSalesRoundDetailsByRoundID(c *fiber.Ctx) error
	UpdateSalesRoundDetailQuantity(c *fiber.Ctx) error
}

type salesRoundDetailController struct {
	salesRoundDetailRepository repositories.SalesRoundDetailRepository
}

func NewSalesRoundDetailController(salesRoundDetailRepository repositories.SalesRoundDetailRepository) SalesRoundDetailController {
	return &salesRoundDetailController{salesRoundDetailRepository: salesRoundDetailRepository}
}

func (h *salesRoundDetailController) CreateSalesRoundDetail(c *fiber.Ctx) error {
	dto := new(dtos.SalesRoundDetailCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	// Fetch the product associated with the product variant
	product, err := h.salesRoundDetailRepository.GetProductByVariantID(dto.VariantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	// Check if there is enough stock
	if product.Stock < dto.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "not enough stock"})
	}

	// Allocate the stock
	product.Stock -= dto.Quantity

	// Create the sales round detail
	salesRoundDetail := &models.SalesRoundDetail{
		RoundID:       dto.RoundID,
		VariantID:     dto.VariantID,
		Quantity:      dto.Quantity,
		QuantityLimit: dto.QuantityLimit,
		Remaining:     product.Stock, // Set the remaining stock
		ProductStock:  product.Stock,
	}

	// Update the product stock
	if err := h.salesRoundDetailRepository.UpdateProductStock(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update product stock"})
	}

	// เปลี่ยนแปลง: รับค่ากลับจาก CreateSalesRoundDetail เป็น 2 ค่า (response และ error)
	response, err := h.salesRoundDetailRepository.CreateSalesRoundDetail(salesRoundDetail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create sales round detail"})
	}

	// เปลี่ยนแปลง: รีเทิร์น response แทนการรีเทิร์น salesRoundDetail โดยตรง
	return c.JSON(response)
}

func (h *salesRoundDetailController) GetAllSalesRoundDetails(c *fiber.Ctx) error {
	expand := c.Query("expand")
	salesRoundDetails, err := h.salesRoundDetailRepository.GetAllSalesRoundDetails(expand)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not fetch sales round details"})
	}
	return c.JSON(salesRoundDetails)
}

// // ข้อผิดพลาดนี้เกิดจากฟังก์ชัน CreateSalesRoundDetail ใน repository มีการแก้ไขให้รีเทิร์น 2 ค่า คือ response และ error
// แต่ใน controller ยังรีเทิร์นค่าจากฟังก์ชัน CreateSalesRoundDetail โดยคาดหวังแค่ 1 ค่าเท่านั้น

func (h *salesRoundDetailController) UpdateSalesRoundDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	detailUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.SalesRoundDetailUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	salesRoundDetail, err := h.salesRoundDetailRepository.GetSalesRoundDetailByID(detailUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "sales round detail not found"})
	}

	salesRoundDetail.RoundID = dto.RoundID
	salesRoundDetail.VariantID = dto.VariantID
	salesRoundDetail.Quantity = dto.Quantity
	salesRoundDetail.QuantityLimit = dto.QuantityLimit
	salesRoundDetail.Remaining = dto.Remaining
	salesRoundDetail.ProductStock = dto.ProductStock

	if err := h.salesRoundDetailRepository.UpdateSalesRoundDetail(salesRoundDetail); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update sales round detail"})
	}

	return c.JSON(salesRoundDetail)
}

func (h *salesRoundDetailController) DeleteSalesRoundDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	detailUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	if err := h.salesRoundDetailRepository.DeleteSalesRoundDetail(detailUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete sales round detail"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *salesRoundDetailController) GetSalesRoundDetailsByRoundID(c *fiber.Ctx) error {
	roundID := c.Params("round_id")
	roundUUID, err := uuid.Parse(roundID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid round ID"})
	}

	salesRoundDetails, err := h.salesRoundDetailRepository.GetSalesRoundDetailsByRoundID(roundUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not fetch sales round details"})
	}
	return c.JSON(salesRoundDetails)
}

func (h *salesRoundDetailController) UpdateSalesRoundDetailQuantity(c *fiber.Ctx) error {
	id := c.Params("id")
	detailUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	var updateQuantityDTO dtos.SalesRoundDetailUpdateQuantityDTO
	if err := c.BodyParser(&updateQuantityDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.salesRoundDetailRepository.UpdateSalesRoundDetailQuantity(detailUUID, updateQuantityDTO.Quantity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update sales round detail quantity"})
	}

	return c.JSON(fiber.Map{
		"message":  "Quantity updated successfully",
		"quantity": updateQuantityDTO.Quantity,
	})
}
