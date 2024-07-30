package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
	"sync"
)

func RegisterProductRoutes(app *fiber.App, controller controllers.ProductController) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Post("/products", controller.CreateProduct) // Route for creating a new product
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Get("/products", controller.GetAllProducts) // Route for getting all products
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Get("/products/with-variants", controller.GetAllProductsWithVariants) // Route for getting all products with variants
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Put("/products/:id", controller.UpdateProduct) // Route for updating a specific product by ID
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Delete("/products/:id", controller.DeleteProduct) // Route for deleting a specific product by ID
	}()

	wg.Wait()
}
