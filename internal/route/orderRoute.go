package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRoutes(app *fiber.App, controller controllers.OrderController) {
	app.Post("/orders", controller.CreateOrder)
	app.Get("/orders", controller.GetAllOrders)
	app.Get("/orders/:id", controller.GetOrderByID) // Add this line to get an order by ID
	app.Put("/orders/:id", controller.UpdateOrder)
	app.Delete("/orders/:id", controller.DeleteOrder)
	app.Post("/orders/complete/:cartID", controller.CompleteOrder)
}
