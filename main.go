package main

import (
	"log"
	"runtime"
	"sync"

	_ "github.com/B6137151/InventoryMarketplaceSystem/docs" // Swagger docs
	"github.com/B6137151/InventoryMarketplaceSystem/internal/controllers"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/route"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/services"
	"github.com/B6137151/InventoryMarketplaceSystem/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger" // swagger middleware for Fiber
	"gorm.io/gorm"
)

// @title Inventory Marketplace System API
// @version 1.0
// @description API documentation for Inventory Marketplace System
// @host localhost:3000
// @BasePath /
func main() {
	runtime.GOMAXPROCS(4)
	app := fiber.New()

	// Swagger endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Logger middleware
	app.Use(logger.New())

	// Initialize the WaitGroup and DB channel for asynchronous setup
	var wg sync.WaitGroup
	dbChan := make(chan *gorm.DB, 1)

	// Setup the database asynchronously using a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		db := database.SetupDatabase()
		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
			log.Fatalf("Failed to create extension: %v", err)
		}

		if err := db.AutoMigrate(
			&models.Store{},
			&models.Category{},
			&models.Product{},
			&models.OrderDetail{},
			&models.Order{},
			&models.SalesRoundDetail{},
			&models.SalesRound{},
			&models.ProductVariant{},
			&models.OrderHistory{},
			&models.Cart{},     // Add Cart model
			&models.CartItem{}, // Add CartItem model
			//&models.Purchase{},
		); err != nil {
			log.Fatalf("Failed to migrate the tables: %v", err)
		}

		dbChan <- db
		close(dbChan)
	}()

	wg.Wait()      // Wait for all goroutines in the WaitGroup to complete
	db := <-dbChan // Receive db instance from the channel
	if db == nil {
		log.Fatal("Failed to connect to the database")
		return
	}

	// Initialize repositories
	storeRepository := repositories.NewStoreRepository(db)
	categoryRepository := repositories.NewCategoryRepository(db)
	customerRepository := repositories.NewCustomerRepository(db)
	productRepository := repositories.NewProductRepository(db)
	productVariantRepository := repositories.NewProductVariantRepository(db)
	salesRoundRepository := repositories.NewSalesRoundRepository(db)
	orderRepository := repositories.NewOrderRepository(db)
	salesRoundDetailRepository := repositories.NewSalesRoundDetailRepository(db)
	orderDetailRepository := repositories.NewOrderDetailRepository(db)
	orderHistoryRepository := repositories.NewOrderHistoryRepository(db)
	cartRepository := repositories.NewCartRepository(db) // New cart repository

	// Initialize services
	purchaseService := services.NewPurchaseService(orderRepository, orderDetailRepository, productVariantRepository, productRepository, salesRoundDetailRepository)

	// Initialize controllers
	storeController := controllers.NewStoreController(storeRepository)
	categoryController := controllers.NewCategoryController(categoryRepository)
	customerController := controllers.NewCustomerController(customerRepository)
	productController := controllers.NewProductController(productRepository)
	productVariantController := controllers.NewProductVariantController(productVariantRepository)

	// Corrected to pass all required arguments
	salesRoundController := controllers.NewSalesRoundController(
		salesRoundRepository,
		productRepository,
		productVariantRepository,
		salesRoundDetailRepository,
		categoryRepository,
		storeRepository,
	)

	salesRoundDetailController := controllers.NewSalesRoundDetailController(salesRoundDetailRepository)
	orderController := controllers.NewOrderController(purchaseService, cartRepository, orderRepository) // Corrected to use cartRepository and orderRepository
	orderDetailController := controllers.NewOrderDetailController(orderDetailRepository)
	orderHistoryController := controllers.NewOrderHistoryController(orderHistoryRepository)
	purchaseController := controllers.NewPurchaseController(purchaseService, salesRoundDetailRepository)
	cartController := controllers.NewCartController(cartRepository, salesRoundRepository, productVariantRepository, storeRepository, productRepository) // New cart controller with all five arguments

	// Register routes
	route.RegisterStoreRoutes(app, storeController)
	route.RegisterCategoryRoutes(app, categoryController)
	route.RegisterCustomerRoutes(app, customerController)
	route.RegisterProductRoutes(app, productController)
	route.RegisterProductVariantRoutes(app, productVariantController)
	route.RegisterSalesRoundRoutes(app, salesRoundController)
	route.RegisterSalesRoundDetailRoutes(app, salesRoundDetailController)
	route.RegisterOrderRoutes(app, orderController)
	route.RegisterOrderDetailRoutes(app, orderDetailController)
	route.RegisterOrderHistoryRoutes(app, orderHistoryController)
	route.RegisterPurchaseRoutes(app, purchaseController)
	route.CartRoute(app, cartController) // New cart route

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running!")
	})

	// Start the server on port 3000 using a goroutine to not block the main goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Use `select{}` to keep the main goroutine running, avoiding the program from exiting
	select {}
}
