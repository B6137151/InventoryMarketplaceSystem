package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrderDetailRoutes(app *fiber.App, controller controllers.OrderDetailController) {
	app.Post("/order-details", controller.CreateOrderDetail)
	app.Get("/order-details", controller.GetAllOrderDetails)
	app.Put("/order-details/:id", controller.UpdateOrderDetail)
	app.Delete("/order-details/:id", controller.DeleteOrderDetail)
}
