package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterCategoryRoutes(app *fiber.App, controller controllers.CategoryController) {
	app.Post("/categories", controller.CreateCategory)
	app.Get("/categories", controller.GetAllCategories)
	app.Put("/categories/:id", controller.UpdateCategory)
	app.Delete("/categories/:id", controller.DeleteCategory)
}
