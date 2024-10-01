package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(app *fiber.App, controller controllers.ProductController) {
	app.Post("/products", controller.CreateProduct)
	app.Get("/products", controller.GetAllProducts)
	app.Get("/products/with-variants", controller.GetAllProductsWithVariants)
	app.Get("/products/:id", controller.GetProductByID)
	app.Put("/products/:id/stock", controller.UpdateProductStock)
	app.Delete("/products/:id", controller.DeleteProduct)
}
