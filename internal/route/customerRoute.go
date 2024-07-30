package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterCustomerRoutes(app *fiber.App, controller controllers.CustomerController) {
	app.Post("/customers", controller.CreateCustomer)
	app.Get("/customers", controller.GetAllCustomers)
	app.Put("/customers/:id", controller.UpdateCustomer)
	app.Delete("/customers/:id", controller.DeleteCustomer)
}
