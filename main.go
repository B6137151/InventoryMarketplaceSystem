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
			//&models.Purchase{}, // Added Purchase model
		); err != nil {
			log.Fatalf("Failed to migrate the tables: %v", err)
		}

		// if err := db.Exec(`
		// 	CREATE TABLE IF NOT EXISTS "sales-round-detail" (
		// 		id             uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
		// 		created_at     timestamp with time zone DEFAULT now(),
		// 		updated_at     timestamp with time zone DEFAULT now(),
		// 		deleted_at     timestamp with time zone,
		// 		round_id       uuid NOT NULL,
		// 		variant_id     uuid NOT NULL,
		// 		quantity       int NOT NULL,
		// 		remaining      int NOT NULL,
		// 		product_stock  int NOT NULL,
		// 		quantity_limit int NOT NULL
		// 	);

		// 	ALTER TABLE "sales-round-detail" OWNER TO postgres;

		// 	CREATE INDEX IF NOT EXISTS "idx_sales-round-detail_variant_id" ON "sales-round-detail" (variant_id);
		// 	CREATE INDEX IF NOT EXISTS "idx_sales-round-detail_round_id" ON "sales-round-detail" (round_id);
		// 	CREATE INDEX IF NOT EXISTS "idx_sales-round-detail_deleted_at" ON "sales-round-detail" (deleted_at);
		// `).Error; err != nil {
		// 	log.Fatalf("Failed to perform manual migration: %v", err)
		// }

		dbChan <- db
		close(dbChan)
	}()

	wg.Wait()      // Wait for all goroutines in the WaitGroup to complete
	db := <-dbChan // Receive db instance from the channel
	if db == nil {
		log.Fatal("Failed to connect to the database")
		return
	}

	// Initialize repositories and controllers, and register routes

	// Store
	storeRepository := repositories.NewStoreRepository(db)
	storeController := controllers.NewStoreController(storeRepository)
	route.RegisterStoreRoutes(app, storeController)

	// Category
	categoryRepository := repositories.NewCategoryRepository(db)
	categoryController := controllers.NewCategoryController(categoryRepository)
	route.RegisterCategoryRoutes(app, categoryController)

	// Customer
	customerRepository := repositories.NewCustomerRepository(db)
	customerController := controllers.NewCustomerController(customerRepository)
	route.RegisterCustomerRoutes(app, customerController)

	// Product
	productRepository := repositories.NewProductRepository(db)
	productController := controllers.NewProductController(productRepository)
	route.RegisterProductRoutes(app, productController)

	// Product Variant
	productVariantRepository := repositories.NewProductVariantRepository(db)
	productVariantController := controllers.NewProductVariantController(productVariantRepository)
	route.RegisterProductVariantRoutes(app, productVariantController)

	// Sales Round
	salesRoundRepository := repositories.NewSalesRoundRepository(db)
	orderRepository := repositories.NewOrderRepository(db) // Initialize order repository
	salesRoundDetailRepository := repositories.NewSalesRoundDetailRepository(db)

	salesRoundController := controllers.NewSalesRoundController(salesRoundRepository, orderRepository, salesRoundDetailRepository)
	route.RegisterSalesRoundRoutes(app, salesRoundController)

	// Sales Round Detail
	salesRoundDetailController := controllers.NewSalesRoundDetailController(salesRoundDetailRepository)
	route.RegisterSalesRoundDetailRoutes(app, salesRoundDetailController)

	// Order
	orderController := controllers.NewOrderController(orderRepository)
	route.RegisterOrderRoutes(app, orderController)

	// Order Detail
	orderDetailRepository := repositories.NewOrderDetailRepository(db)
	orderDetailController := controllers.NewOrderDetailController(orderDetailRepository)
	route.RegisterOrderDetailRoutes(app, orderDetailController)

	// Order History
	orderHistoryRepository := repositories.NewOrderHistoryRepository(db)
	orderHistoryController := controllers.NewOrderHistoryController(orderHistoryRepository)
	route.RegisterOrderHistoryRoutes(app, orderHistoryController)

	// Purchase
	//purchaseRepository := repositories.NewPurchaseRepository(db)
	//purchaseController := controllers.NewPurchaseController(purchaseRepository, productVariantRepository, salesRoundDetailRepository)
	//route.RegisterPurchaseRoutes(app, purchaseController)

	// Serve a simple message at the root URL
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
