package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrderHistoryRoutes(app *fiber.App, controller controllers.OrderHistoryController) {
	app.Post("/order-histories", controller.CreateOrderHistory)
	app.Get("/order-histories", controller.GetAllOrderHistories)
	app.Put("/order-histories/:id", controller.UpdateOrderHistory)
	app.Delete("/order-histories/:id", controller.DeleteOrderHistory)
}
