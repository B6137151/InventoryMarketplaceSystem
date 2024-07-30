package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterProductVariantRoutes(app *fiber.App, controller controllers.ProductVariantController) {
	app.Post("/product-variants", controller.CreateProductVariant)
	app.Get("/product-variants", controller.GetAllProductVariants)
	app.Put("/product-variants/:id", controller.UpdateProductVariant)
	app.Delete("/product-variants/:id", controller.DeleteProductVariant)
}
