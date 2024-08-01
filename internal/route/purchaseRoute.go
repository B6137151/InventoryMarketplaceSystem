package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterPurchaseRoutes(app *fiber.App, controller controllers.PurchaseController) {
	app.Post("/purchases", controller.MakePurchase)
}
