package controllers

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/services"
	"github.com/gofiber/fiber/v2"
)

type PurchaseController interface {
	MakePurchase(c *fiber.Ctx) error
}

type purchaseController struct {
	PurchaseService services.PurchaseService
}

func NewPurchaseController(purchaseService services.PurchaseService) PurchaseController {
	return &purchaseController{PurchaseService: purchaseService}
}

// MakePurchase godoc
// @Summary Make a new purchase
// @Description Make a new purchase
// @Tags Purchases
// @Accept json
// @Produce json
// @Param purchase body dtos.PurchaseCreateDTO true "Purchase"
// @Success 201 {object} dtos.OrderResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /purchases [post]
func (h *purchaseController) MakePurchase(c *fiber.Ctx) error {
	dto := new(dtos.PurchaseCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	response, err := h.PurchaseService.MakePurchase(*dto)
	if err != nil {
		if err.Error() == "not enough stock" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "not enough stock"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
