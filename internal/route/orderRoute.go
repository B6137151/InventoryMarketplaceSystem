package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRoutes(app *fiber.App, controller controllers.OrderController) {
	app.Post("/orders", controller.CreateOrder)
	app.Get("/orders", controller.GetAllOrders)
	app.Put("/orders/:id", controller.UpdateOrder)
	app.Delete("/orders/:id", controller.DeleteOrder)
}
