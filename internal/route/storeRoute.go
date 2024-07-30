package route

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterStoreRoutes(app *fiber.App, controller controllers.StoreController) {
	app.Post("/stores", controller.CreateStore)                     // Create a new store
	app.Get("/stores", controller.GetAllStores)                     // Get all stores
	app.Get("/stores/:id", validateUUID, controller.GetStoreByID)   // Get a store by ID
	app.Put("/stores/:id", validateUUID, controller.UpdateStore)    // Update a store by ID
	app.Delete("/stores/:id", validateUUID, controller.DeleteStore) // Delete a store by ID
}

func validateUUID(c *fiber.Ctx) error {
	id := c.Params("id")
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		if _, err := uuid.Parse(id); err != nil {
			errChan <- err
		}
	}()

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID format"})
	}
	return c.Next()
}
