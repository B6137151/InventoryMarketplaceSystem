package controllers

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"log"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderController interface {
	CreateOrder(c *fiber.Ctx) error
	GetAllOrders(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
	DeleteOrder(c *fiber.Ctx) error
	CompleteOrder(c *fiber.Ctx) error
	GetOrderByID(c *fiber.Ctx) error // Add GetOrderByID to the interface
}

type orderController struct {
	purchaseService services.PurchaseService
	cartRepository  repositories.CartRepository
	orderRepository repositories.OrderRepository
}

func NewOrderController(purchaseService services.PurchaseService, cartRepository repositories.CartRepository, orderRepository repositories.OrderRepository) OrderController {
	return &orderController{
		purchaseService: purchaseService,
		cartRepository:  cartRepository,
		orderRepository: orderRepository,
	}
}

// GetOrderByID retrieves an order by its ID
func (ctrl *orderController) GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Parse UUID from the string
	orderID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	// Fetch the order from the database using the repository
	order, err := ctrl.orderRepository.GetOrderByID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving the order",
		})
	}

	// Return the order details as a response
	return c.JSON(order)
}

func (ctrl *orderController) CreateOrder(c *fiber.Ctx) error {
	dto := new(dtos.OrderCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	orderID := uuid.New()
	order := models.Order{
		ID:              orderID,
		CustomerID:      dto.CustomerID,
		RoundID:         dto.RoundID,
		Status:          models.OrderStatus("Processing"),
		OrderDate:       time.Now(),
		TotalPrice:      dto.TotalPrice,
		DeliveryAddress: dto.DeliveryAddress,
		PaymentSource:   dto.PaymentSource, // Now treated as a regular string
	}
	if err := ctrl.orderRepository.CreateOrder(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"order_id": orderID})
}

func (ctrl *orderController) GetAllOrders(c *fiber.Ctx) error {
	orders, err := ctrl.purchaseService.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	orderResponses := make([]dtos.OrderResponseDTO, len(orders))
	for i, order := range orders {
		orderResponses[i] = dtos.OrderResponseDTO{
			ID:              order.ID,
			OrderDate:       order.OrderDate,
			Status:          string(order.Status),
			TotalPrice:      order.TotalPrice,
			DeliveryAddress: order.DeliveryAddress,
			PaymentSource:   order.PaymentSource, // Now treated as a regular string
			Customer: dtos.CustomerDTO{
				ID:    order.Customer.ID,
				Name:  order.Customer.Name,
				Email: order.Customer.Email,
			},
			SalesRound: dtos.SalesRoundDTO{
				ID:        order.SalesRound.ID,
				Name:      order.SalesRound.Name,
				StartDate: order.SalesRound.StartDate,
				EndDate:   order.SalesRound.EndDate,
			},
			PaymentInfo: dtos.PaymentInfoDTO{
				PaymentMethod: order.PaymentInfo.PaymentMethod,
				PaymentStatus: order.PaymentInfo.PaymentStatus,
				TransactionID: order.PaymentInfo.TransactionID,
			},
			ShippingInfo: dtos.ShippingInfoDTO{
				ShippingMethod:        order.ShippingInfo.ShippingMethod,
				EstimatedDeliveryDate: order.ShippingInfo.EstimatedDeliveryDate,
			},
		}
	}
	return c.JSON(orderResponses)
}

func (ctrl *orderController) UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.OrderUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	order, err := ctrl.purchaseService.GetOrderByID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	order.Status = models.OrderStatus(dto.Status)
	order.TotalPrice = dto.TotalPrice
	order.DeliveryAddress = dto.DeliveryAddress
	order.PaymentSource = dto.PaymentSource // Now treated as a regular string

	if err := ctrl.purchaseService.UpdateOrder(order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "updated"})
}

func (ctrl *orderController) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}
	if err := ctrl.purchaseService.DeleteOrder(orderID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *orderController) CompleteOrder(c *fiber.Ctx) error {
	// Log the start of the request
	log.Println("Starting CompleteOrder request")

	// Parse the cart ID from the request
	cartID := c.Params("cartID")
	log.Printf("Received cartID: %s", cartID)

	// Convert cartID string to uuid.UUID
	cartUUID, err := uuid.Parse(cartID)
	if err != nil {
		log.Printf("Invalid cart ID format: %s, error: %v", cartID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid cart ID format",
			"error":   err.Error(),
		})
	}

	log.Printf("Parsed cart UUID: %s", cartUUID)

	var cart *models.Cart
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// Log before fetching cart details
	log.Println("Fetching cart details")
	wg.Add(1)
	go func() {
		defer wg.Done()
		cart, err = ctrl.cartRepository.GetCartByID(cartUUID)
		if err != nil {
			log.Printf("Error fetching cart with ID %s: %v", cartUUID, err)
			errChan <- fmt.Errorf("Cart not found: %v", err)
			return
		}
		log.Printf("Fetched cart: %+v", cart)
	}()
	wg.Wait()

	// Check if any errors occurred during fetching the cart
	select {
	case err := <-errChan:
		log.Printf("Failed to fetch cart: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to fetch cart",
			"error":   err.Error(),
		})
	default:
		log.Println("Successfully fetched cart")
	}

	// Ensure the cart has a valid PaymentSource, else return an error
	log.Println("Checking if PaymentSource is set in the cart")
	if cart.PaymentSource == "" {
		log.Println("PaymentSource is not set in the cart")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "PaymentSource is not set in the cart",
		})
	}

	// Log before checking existing order
	log.Println("Checking if an order already exists for this cart")
	existingOrder, err := ctrl.orderRepository.GetOrderByCartID(cartUUID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Error checking existing order for cart ID %s: %v", cartUUID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check existing order",
			"error":   err.Error(),
		})
	}

	var order models.Order
	if existingOrder != nil {
		// If an order already exists, use the existing order ID
		order = *existingOrder
		log.Printf("Using existing order ID: %s", order.ID)
	} else {
		// If no existing order, create a new order
		log.Println("Creating a new order")
		orderID := uuid.New()
		order = models.Order{
			ID:              orderID,
			CustomerID:      cart.CustomerID,
			RoundID:         cart.RoundID,
			Status:          models.OrderStatus("Processing"),
			OrderDate:       time.Now(),
			PaymentSource:   cart.PaymentSource, // Now treated as a regular string
			DeliveryAddress: cart.DeliveryAddress,
			CartID:          cart.ID,
			TotalPrice:      cart.TotalAmount,

			PaymentInfo: models.PaymentInfo{
				PaymentMethod: cart.PaymentSource, // Now treated as a regular string
				TransactionID: "ABC123",
				PaymentStatus: "Completed",
				PaidAmount:    cart.TotalAmount,
				PaidCurrency:  cart.Currency,
			},
			ShippingInfo: models.ShippingInfo{
				ShippingMethod:   "Express",
				TrackingNumber:   "TRACK123456",
				CarrierName:      "DHL",
				ShippingStatus:   "Shipped",
				EstimatedArrival: time.Now().AddDate(0, 0, 5),
			},
		}

		// Log before order creation
		log.Println("Inserting new order into the database")
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := ctrl.orderRepository.CreateOrder(&order); err != nil {
				log.Printf("Error creating order for cart ID %s: %v", cartUUID, err)
				errChan <- fmt.Errorf("Failed to create order: %v", err)
			} else {
				log.Printf("Order created successfully: %+v", order)
			}
		}()
		wg.Wait()
	}

	// Check if any errors occurred during the order creation process
	select {
	case err := <-errChan:
		log.Printf("Failed to create order: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create order",
			"error":   err.Error(),
		})
	default:
		log.Printf("Order completed successfully with ID: %s", order.ID)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":        "Order completed successfully",
			"order_id":       order.ID,
			"payment_source": order.PaymentSource,
		})
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
