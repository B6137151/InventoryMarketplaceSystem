package database

import (
	"fmt"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func SetupDatabase() *gorm.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)

	log.Println("Starting database setup...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	log.Println("Database connection established.")

	// Auto-migrate all tables
	if err := db.AutoMigrate(
		&models.Store{},
		&models.Category{},
		&models.Product{},
		&models.OrderDetail{},
		&models.Order{},
		//&models.SalesRoundDetail{},
		&models.SalesRound{},
		&models.ProductVariant{},
		&models.OrderHistory{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate tables: %v", err)
		os.Exit(1)
	}
	log.Println("Tables migrated successfully.")

	return db
}
