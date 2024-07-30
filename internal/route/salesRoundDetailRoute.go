package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterSalesRoundDetailRoutes(app *fiber.App, controller controllers.SalesRoundDetailController) {
	app.Post("/sales-round-details", controller.CreateSalesRoundDetail)                     // Route for creating a new sales round detail
	app.Get("/sales-round-details", controller.GetAllSalesRoundDetails)                     // Route for getting all sales round details
	app.Put("/sales-round-details/:id", controller.UpdateSalesRoundDetail)                  // Route for updating a specific sales round detail by ID
	app.Delete("/sales-round-details/:id", controller.DeleteSalesRoundDetail)               // Route for deleting a specific sales round detail by ID
	app.Put("/sales-round-details/:id/quantity", controller.UpdateSalesRoundDetailQuantity) // New route for updating the quantity of a specific sales round detail by ID
}
