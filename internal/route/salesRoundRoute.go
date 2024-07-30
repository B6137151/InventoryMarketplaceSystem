package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterSalesRoundRoutes(app *fiber.App, controller controllers.SalesRoundController) {
	app.Post("/sales-rounds", controller.CreateSalesRound)                // Route for creating a new sales round
	app.Get("/sales-rounds", controller.GetAllSalesRounds)                // Route for getting all sales rounds
	app.Get("/sales-rounds/:id", controller.GetSalesRoundDetails)         // Route for getting sales round details by ID
	app.Put("/sales-rounds/:id", controller.UpdateSalesRound)             // Route for updating a sales round by ID
	app.Delete("/sales-rounds/:id", controller.DeleteSalesRound)          // Route for deleting a sales round by ID
	app.Get("/sales-rounds/:id/details", controller.GetSalesRoundDetails) // Specific endpoint for sales round details
	app.Get("/sales-rounds/combined", controller.GetCombinedSalesRoundProductData)
}
