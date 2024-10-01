package controllers

import (
	"log"

	"strings"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartController struct {
	cartRepository           repositories.CartRepository
	salesRoundRepository     repositories.SalesRoundRepository
	productVariantRepository repositories.ProductVariantRepository
	storeRepository          repositories.StoreRepository
	productRepository        repositories.ProductRepository
	customerRepository       repositories.CustomerRepository
}

func NewCartController(cartRepo repositories.CartRepository, salesRoundRepo repositories.SalesRoundRepository, productVariantRepo repositories.ProductVariantRepository, storeRepo repositories.StoreRepository, productRepo repositories.ProductRepository) *CartController {
	return &CartController{
		cartRepository:           cartRepo,
		salesRoundRepository:     salesRoundRepo,
		productVariantRepository: productVariantRepo,
		storeRepository:          storeRepo,
		productRepository:        productRepo,
	}
}

// CreateCart godoc
// @Summary Create a cart
// @Description Create a new cart or update an existing one
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart body dtos.CartRequestDTO true "Cart"
// @Success 201 {object} dtos.CartInteractionResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /cart [post]
func (h *CartController) CreateCart(c *fiber.Ctx) error {
	dto := new(dtos.CartRequestDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	customerID := dto.CustomerID
	roundID := dto.RoundID
	storeID := dto.StoreID

	// ตรวจสอบว่ามีตะกร้าสำหรับลูกค้านี้อยู่แล้วหรือไม่
	existingCart, err := h.cartRepository.GetCart(customerID, roundID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve cart"})
	}

	var items []models.CartItem
	var totalAmount float64

	for _, item := range dto.Items {
		variant, err := h.productVariantRepository.GetProductVariantByID(item.VariantID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid variant ID"})
		}

		product, err := h.productRepository.GetProductByID(variant.ProductID)
		if err != nil {
			log.Printf("Error fetching product details for ProductID %s: %v", variant.ProductID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching product details"})
		}

		totalPrice := float64(item.Quantity) * variant.Price
		totalAmount += totalPrice

		items = append(items, models.CartItem{
			VariantID:   item.VariantID,
			Quantity:    item.Quantity,
			ProductName: product.ProductName,
			SKUCode:     variant.SKUCode,
			Price:       variant.Price,
		})
	}

	if existingCart != nil {
		// If the cart already exists, add new items to it and update the total amount
		existingCart.Items = append(existingCart.Items, items...)
		existingCart.TotalAmount += totalAmount

		if err := h.cartRepository.UpdateCart(existingCart); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update cart"})
		}

		// Retrieve the store information after updating the cart
		store, err := h.storeRepository.GetStoreByID(existingCart.StoreID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch store details"})
		}

		// Prepare response items
		responseItems := make([]dtos.CartItemDTO, len(existingCart.Items))
		for i, item := range existingCart.Items {
			responseItems[i] = dtos.CartItemDTO{
				VariantID: item.VariantID,
				Quantity:  item.Quantity,
			}
		}

		return c.Status(fiber.StatusOK).JSON(dtos.CartInteractionResponseDTO{
			Message:     "Cart updated successfully",
			CartID:      existingCart.ID,
			CustomerID:  existingCart.CustomerID,
			RoundID:     existingCart.RoundID,
			StoreID:     existingCart.StoreID,
			StoreName:   store.StoreName, // Use the fetched store name
			Items:       responseItems,
			TotalAmount: existingCart.TotalAmount,
		})
	}

	// ถ้าไม่มีตะกร้าสินค้า ให้สร้างใหม่
	cart := models.Cart{
		CustomerID:  customerID,
		RoundID:     roundID,
		StoreID:     storeID,
		Items:       items,
		TotalAmount: totalAmount, // Ensure you have this field in your Cart model
	}

	// Adjusted to handle three return values
	storeName, customerName, err := h.cartRepository.CreateCart(&cart)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create cart"})
	}

	responseItems := make([]dtos.CartItemDTO, len(items))
	for i, item := range items {
		responseItems[i] = dtos.CartItemDTO{
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		}
	}

	return c.Status(fiber.StatusCreated).JSON(dtos.CartInteractionResponseDTO{
		Message:      "Cart created successfully",
		CartID:       cart.ID,
		CustomerID:   customerID,
		RoundID:      roundID,
		StoreID:      storeID,
		StoreName:    storeName,    // Use the storeName returned by CreateCart
		CustomerName: customerName, // Include the customerName in the response
		Items:        responseItems,
		TotalAmount:  totalAmount,
	})
}

// Checkout godoc
// @Summary Checkout a cart
// @Description Complete the checkout process for a cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param checkout body dtos.CheckoutRequestDTO true "Checkout"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /cart/checkout [post]
func (h *CartController) Checkout(c *fiber.Ctx) error {
	// Parse the request body into the CheckoutRequestDTO
	dto := new(dtos.CheckoutRequestDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Log the DTO to check the CartID and other details
	log.Printf("Received Checkout request - CartID: %s, PaymentSource: %s, DeliveryAddress: %s, Currency: %s",
		dto.CartID, dto.PaymentSource, dto.DeliveryAddress, dto.Currency)

	// Fetch the cart by CartID
	cart, err := h.cartRepository.GetCartByID(dto.CartID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Log detailed info if the cart was not found
			log.Printf("Error fetching cart by ID: %s - record not found", dto.CartID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
		}
		// Log any other error
		log.Printf("Error fetching cart by ID: %s - %v", dto.CartID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve cart"})
	}

	// Log the cart details to ensure correct retrieval
	log.Printf("Fetched cart - CartID: %s, CustomerID: %s, RoundID: %s, StoreID: %s, TotalAmount: %f",
		cart.ID, cart.CustomerID, cart.RoundID, cart.StoreID, cart.TotalAmount)

	// Update payment source, delivery address, and currency
	cart.PaymentSource = dto.PaymentSource
	cart.DeliveryAddress = dto.DeliveryAddress
	cart.Currency = dto.Currency

	// Log the updated cart details before saving
	log.Printf("Updating cart - CartID: %s, PaymentSource: %s, DeliveryAddress: %s, Currency: %s",
		cart.ID, cart.PaymentSource, cart.DeliveryAddress, cart.Currency)

	// Update the cart in the database
	if err := h.cartRepository.UpdateCart(cart); err != nil {
		log.Printf("Error updating cart: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not complete checkout"})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Checkout completed successfully",
		"cart_id":        cart.ID,
		"payment_source": cart.PaymentSource,
	})
}

// GetCartDetails godoc
// @Summary Get details of a cart
// @Description Get details of a cart by cart ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart_id path string true "Cart ID"
// @Param expand query string false "Comma separated list of related entities to expand (e.g., customer, store)"
// @Success 200 {object} dtos.CartDetailResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /cart/{cart_id}/details [get]
func (h *CartController) GetCartDetails(c *fiber.Ctx) error {
	// Log the request method, path, and parameters
	log.Printf("Received request: %s %s", c.Method(), c.Path())

	// Extract expand query parameter
	expand := c.Query("expand")
	log.Printf("Expand query parameter: '%s'", expand)
	expandFields := strings.Split(expand, ",")

	// Log all route parameters for debugging
	log.Printf("All route parameters: %+v", c.AllParams())

	// Extract the cart_id parameter
	log.Printf("Extracting cart_id parameter...")
	cartIDParam := c.Params("cart_id")
	log.Printf("Extracted cart_id: '%s'", cartIDParam)

	// Check if the cart_id parameter is missing or empty
	if cartIDParam == "" {
		log.Println("Cart ID is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cart ID is required"})
	}

	// Try parsing the cart_id into UUID
	cartID, err := uuid.Parse(cartIDParam)
	if err != nil {
		log.Printf("Invalid cart ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID"})
	}

	log.Printf("Fetching detailed cart for cart ID: %s", cartIDParam)

	// Fetch cart details with optional expansion based on expandFields
	cart, err := h.cartRepository.GetDetailedCartByCustomerID(cartID, expandFields)
	if err != nil {
		log.Printf("Cart not found for cart ID %s: %v", cartIDParam, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
	}

	log.Printf("Fetched detailed cart: %+v", cart)

	var items []dtos.CartItemDetailDTO

	for _, item := range cart.Items {
		log.Printf("Processing cart item: %+v", item)
		variant, err := h.productVariantRepository.GetProductVariantByID(item.VariantID)
		if err != nil {
			log.Printf("Error fetching product variant details for VariantID %s: %v", item.VariantID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching product variant details"})
		}

		product, err := h.productRepository.GetProductByID(variant.ProductID)
		if err != nil {
			log.Printf("Error fetching product details for ProductID %s: %v", variant.ProductID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching product details"})
		}

		log.Printf("Fetched product details: %+v", product)
		totalPrice := float64(item.Quantity) * variant.Price

		items = append(items, dtos.CartItemDetailDTO{
			VariantID:   item.VariantID,
			Quantity:    item.Quantity,
			ProductName: product.ProductName,
			SKUCode:     variant.SKUCode,
			Price:       variant.Price,
			TotalPrice:  totalPrice,
			ImageURL:    variant.ImageURL,
		})
	}

	// Create MetaData object (Line 69)
	metaData := dtos.MetaData{
		Total: len(cart.Items), // จำนวนรายการทั้งหมด
		Count: len(items),      // จำนวนรายการที่ถูกประมวลผล
	}

	// Response with cart details and MetaData (Line 74)
	response := dtos.CartDetailResponseDTO{
		Message:       "Cart details fetched successfully",
		CartID:        cart.ID,
		CustomerID:    cart.CustomerID,
		CustomerName:  cart.Customer.Name,
		RoundID:       cart.RoundID,
		StoreID:       cart.StoreID,
		StoreName:     cart.Store.StoreName,
		PaymentSource: cart.PaymentSource,
		Items:         items,
		TotalAmount:   cart.TotalAmount, // ใช้ totalAmount จากฐานข้อมูล
		MetaData:      metaData,         // ใส่ MetaData ใน response (บรรทัดที่ 81)
	}

	log.Printf("Final response: %+v", response)
	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdateCart godoc
// @Summary Update a cart
// @Description Update a cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart body dtos.CartRequestDTO true "Cart"
// @Success 200 {object} dtos.CartInteractionResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /cart [put]
func (h *CartController) UpdateCart(c *fiber.Ctx) error {
	dto := new(dtos.CartRequestDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	customerID, err := uuid.Parse(dto.CustomerID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid customer ID format"})
	}

	roundID, err := uuid.Parse(dto.RoundID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid round ID format"})
	}

	// Validate sales round exists
	if _, err := h.salesRoundRepository.GetSalesRoundByID(roundID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid round ID"})
	}

	var items []models.CartItem
	for _, item := range dto.Items {
		variantID, err := uuid.Parse(item.VariantID.String())
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid variant ID format"})
		}

		// Validate product variant exists
		if _, err := h.productVariantRepository.GetProductVariantByID(variantID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid variant ID"})
		}

		items = append(items, models.CartItem{
			VariantID: variantID,
			Quantity:  item.Quantity,
		})
	}

	cart := models.Cart{
		CustomerID: customerID,
		RoundID:    roundID,
		Items:      items,
	}

	if err := h.cartRepository.UpdateCart(&cart); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update cart"})
	}

	var responseItems []dtos.CartItemDTO
	for _, item := range cart.Items {
		responseItems = append(responseItems, dtos.CartItemDTO{
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		})
	}

	return c.JSON(dtos.CartInteractionResponseDTO{
		Message:    "Cart updated successfully",
		CustomerID: cart.CustomerID,
		RoundID:    cart.RoundID,
		Items:      responseItems,
	})
}

// DeleteCart godoc
// @Summary Delete a cart
// @Description Delete a cart by customer ID and round ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param customer_id query string true "Customer ID"
// @Param round_id query string true "Round ID"
// @Success 204
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /cart [delete]
func (h *CartController) DeleteCart(c *fiber.Ctx) error {
	customerID := c.Query("customer_id")
	roundID := c.Query("round_id")

	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid customer ID"})
	}

	roundUUID, err := uuid.Parse(roundID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid round ID"})
	}

	if err := h.cartRepository.DeleteCart(customerUUID, roundUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete cart"})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Cart deleted successfully"})
}
