package repositories

import (
	"errors"
	"log"
	"strings"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CartRepository defines the methods for cart data operations.
type CartRepository interface {
	CreateCart(cart *models.Cart) (string, string, error) // Updated to match the new method signature
	GetCart(customerID uuid.UUID, roundID uuid.UUID) (*models.Cart, error)
	UpdateCart(cart *models.Cart) error
	DeleteCart(customerID uuid.UUID, roundID uuid.UUID) error
	GetDetailedCartByCustomerID(customerID uuid.UUID, expandFields []string) (*models.Cart, error)
	GetCartByID(cartID uuid.UUID) (*models.Cart, error)
	GetAggregatedCartByCustomerID(customerID uuid.UUID) (*models.Cart, []models.OrderItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new CartRepository.
func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

// CreateCart creates a new cart and returns the associated store name and customer name
func (r *cartRepository) CreateCart(cart *models.Cart) (string, string, error) {
	// Start a new transaction
	tx := r.db.Begin()
	log.Printf("Starting transaction for creating cart: %+v", cart)

	// Ensure the transaction is rolled back if an error occurs
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to panic: %v", r)
		}
	}()

	// Retrieve the store information based on StoreID
	var store models.Store
	log.Printf("Retrieving store information for StoreID: %s", cart.StoreID)
	if err := tx.Where("id = ?", cart.StoreID).First(&store).Error; err != nil {
		tx.Rollback()
		log.Printf("Error retrieving store: %v", err)
		return "", "", err
	}
	log.Printf("Store retrieved: %+v", store)

	// Retrieve the customer information based on CustomerID
	var customer models.Customer
	log.Printf("Retrieving customer information for CustomerID: %s", cart.CustomerID)
	if err := tx.Where("id = ?", cart.CustomerID).First(&customer).Error; err != nil {
		tx.Rollback()
		log.Printf("Error retrieving customer: %v", err)
		return "", "", err
	}
	log.Printf("Customer retrieved: %+v", customer)

	// Attempt to create the cart
	log.Printf("Attempting to create cart: %+v", cart)
	if err := tx.Create(cart).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating cart: %v", err)
		return "", "", err
	}
	log.Printf("Cart created successfully: %+v", cart)

	// Iterate over the cart items and create or update them as needed
	for _, item := range cart.Items {
		var existingItem models.CartItem
		log.Printf("Checking for existing cart item with CartID: %s and VariantID: %s", cart.ID, item.VariantID)
		err := tx.Where("cart_id = ? AND variant_id = ?", cart.ID, item.VariantID).First(&existingItem).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Item doesn't exist, so create a new one
			item.CartID = cart.ID // Ensure the CartID is set on the item
			log.Printf("Creating new cart item: %+v", item)
			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				log.Printf("Error creating cart item: %v", err)
				return "", "", err
			}
			log.Printf("Cart item created successfully: %+v", item)
		} else if err != nil {
			// Some other error occurred
			tx.Rollback()
			log.Printf("Error checking existing cart item: %v", err)
			return "", "", err
		} else {
			// Item already exists, so update the quantity
			log.Printf("Updating quantity for existing cart item: %+v", existingItem)
			existingItem.Quantity += item.Quantity
			if err := tx.Save(&existingItem).Error; err != nil {
				tx.Rollback()
				log.Printf("Error updating cart item quantity: %v", err)
				return "", "", err
			}
			log.Printf("Cart item quantity updated successfully: %+v", existingItem)
		}
	}

	// Commit the transaction if no errors occurred
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return "", "", err
	}
	log.Printf("Transaction committed successfully")

	// Return the store name and customer name along with success
	return store.StoreName, customer.Name, nil
}

func (r *cartRepository) GetCart(customerID uuid.UUID, roundID uuid.UUID) (*models.Cart, error) {
	log.Printf("Fetching cart for CustomerID: %s and RoundID: %s", customerID, roundID)
	var cart models.Cart
	// Preload Items, Customer, and Store
	err := r.db.Preload("Items").Preload("Customer").Preload("Store").Where("customer_id = ? AND round_id = ?", customerID, roundID).First(&cart).Error
	if err != nil {
		log.Printf("Error fetching cart: %v", err)
		return nil, err
	}

	// Log the fetched cart
	log.Printf("Fetched cart: %+v", cart)

	// Log Customer details if it exists
	if cart.Customer.ID != uuid.Nil {
		log.Printf("Customer details - ID: %s, Name: %s, Email: %s", cart.Customer.ID, cart.Customer.Name, cart.Customer.Email)
	} else {
		log.Println("Customer is missing or not properly loaded.")
	}

	// Log Store details if it exists
	if cart.Store.ID != uuid.Nil {
		log.Printf("Store details - ID: %s, Name: %s", cart.Store.ID, cart.Store.StoreName)
	} else {
		log.Println("Store is missing or not properly loaded.")
	}

	// Log Items in the cart
	if len(cart.Items) > 0 {
		for i, item := range cart.Items {
			log.Printf("Item %d - VariantID: %s, ProductName: %s, Quantity: %d, Price: %f", i+1, item.VariantID, item.ProductName, item.Quantity, item.Price)
		}
	} else {
		log.Println("No items found in the cart.")
	}

	return &cart, nil
}

func (r *cartRepository) GetCartByID(cartID uuid.UUID) (*models.Cart, error) {
	log.Printf("Fetching cart by CartID: %s", cartID)
	var cart models.Cart
	err := r.db.Preload("Items").Where("id = ?", cartID).First(&cart).Error
	if err != nil {
		log.Printf("Error fetching cart by ID: %v", err)
		return nil, err
	}
	log.Printf("Fetched cart by ID: %+v", cart)
	return &cart, nil
}

func (r *cartRepository) UpdateCart(cart *models.Cart) error {
	log.Printf("Updating cart: %+v", cart)
	err := r.db.Save(cart).Error
	if err != nil {
		log.Printf("Error updating cart: %v", err)
	}
	return err
}

func (r *cartRepository) DeleteCart(customerID uuid.UUID, roundID uuid.UUID) error {
	log.Printf("Deleting cart for CustomerID: %s and RoundID: %s", customerID, roundID)
	err := r.db.Where("customer_id = ? AND round_id = ?", customerID, roundID).Delete(&models.Cart{}).Error
	if err != nil {
		log.Printf("Error deleting cart: %v", err)
	}
	return err
}

func (r *cartRepository) GetDetailedCartByCustomerID(customerID uuid.UUID, expandFields []string) (*models.Cart, error) {
	log.Printf("Fetching detailed cart for CustomerID: %s", customerID)
	var cart models.Cart
	db := r.db // Initialize the db object to chain conditional preload queries

	// Always preload ProductVariant and Product
	db = db.Preload("Items.ProductVariant").Preload("Items.ProductVariant.Product")

	// Conditionally preload related entities based on the expandFields parameter
	for _, field := range expandFields {
		switch strings.TrimSpace(field) {
		case "store":
			// Preload the Store entity if 'store' is specified in the expand query
			db = db.Preload("Store")
		case "sales_round":
			// Preload the SalesRound entity if 'sales_round' is specified in the expand query
			db = db.Preload("SalesRound")
		case "customer":
			// Preload the Customer entity if 'customer' is specified in the expand query
			db = db.Preload("Customer")
		case "items":
			// Preload additional entities for items if 'items' is expanded
			db = db.Preload("Items.ProductVariant.Product.Store").
				Preload("Items.ProductVariant.Product.Category")
		}
	}

	// Execute the query, filtering by the customer ID
	err := db.Where("customer_id = ?", customerID).First(&cart).Error
	if err != nil {
		log.Printf("Error fetching detailed cart: %v", err)
		return nil, err
	}

	log.Printf("Fetched detailed cart: %+v", cart)
	return &cart, nil
}

func (r *cartRepository) GetAggregatedCartByCustomerID(customerID uuid.UUID) (*models.Cart, []models.OrderItem, error) {
	log.Printf("Fetching aggregated cart for CustomerID: %s", customerID)

	// Pass an empty expandFields slice or define what fields to expand
	expandFields := []string{} // You can specify what to expand, e.g., []string{"store", "items"}

	// Fix: Pass expandFields as the second argument
	cart, err := r.GetDetailedCartByCustomerID(customerID, expandFields)
	if err != nil {
		log.Printf("Error fetching detailed cart for aggregation: %v", err)
		return nil, nil, err
	}

	// Aggregate cart items into order items, considering product variants
	aggregatedItems, err := r.aggregateItems(cart.Items)
	if err != nil {
		log.Printf("Error aggregating cart items: %v", err)
		return nil, nil, err
	}

	log.Printf("Aggregated cart items: %+v", aggregatedItems)
	return cart, aggregatedItems, nil
}

// Helper function to aggregate cart items
func (r *cartRepository) aggregateItems(cartItems []models.CartItem) ([]models.OrderItem, error) {
	log.Printf("Aggregating cart items: %+v", cartItems)
	itemMap := make(map[string]models.OrderItem)

	for _, item := range cartItems {
		// Fetch ProductVariant for the current CartItem
		log.Printf("Fetching ProductVariant for CartItem with VariantID: %s", item.VariantID)
		variant, err := r.getProductVariantByID(item.VariantID)
		if err != nil {
			log.Printf("Error fetching ProductVariant: %v", err)
			return nil, err
		}
		log.Printf("Fetched ProductVariant: %+v", variant)

		if existingItem, ok := itemMap[item.SKUCode]; ok {
			existingItem.Quantity += item.Quantity
			existingItem.TotalPrice += float64(item.Quantity) * variant.Price
			itemMap[item.SKUCode] = existingItem
		} else {
			itemMap[item.SKUCode] = models.OrderItem{
				ProductName:  variant.Product.ProductName, // Ensure the correct product name is used
				SKUCode:      item.SKUCode,
				Quantity:     item.Quantity,
				PricePerUnit: variant.Price,
				TotalPrice:   float64(item.Quantity) * variant.Price,
				VariantID:    item.VariantID, // Store the VariantID for future reference
			}
		}
		log.Printf("Current aggregation map: %+v", itemMap)
	}

	orderItems := make([]models.OrderItem, 0, len(itemMap))
	for _, v := range itemMap {
		orderItems = append(orderItems, v)
	}
	log.Printf("Final aggregated order items: %+v", orderItems)
	return orderItems, nil
}

// Helper function to fetch a ProductVariant by ID
func (r *cartRepository) getProductVariantByID(variantID uuid.UUID) (*models.ProductVariant, error) {
	log.Printf("Fetching ProductVariant by ID: %s", variantID)
	var variant models.ProductVariant
	if err := r.db.Preload("Product").Where("id = ?", variantID).First(&variant).Error; err != nil {
		log.Printf("Error fetching ProductVariant by ID: %v", err)
		return nil, err
	}
	log.Printf("Fetched ProductVariant: %+v", variant)
	return &variant, nil
}
