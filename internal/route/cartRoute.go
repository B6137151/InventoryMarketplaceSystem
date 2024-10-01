package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func CartRoute(app *fiber.App, controller *controllers.CartController) {
	app.Post("/cart", controller.CreateCart)
	app.Put("/cart", controller.UpdateCart)
	app.Delete("/cart", controller.DeleteCart)
	app.Get("/cart/:cart_id/details", controller.GetCartDetails) // Update to use cart_id
	app.Post("/cart/checkout", controller.Checkout)              // Add this line to define the checkout route
}
