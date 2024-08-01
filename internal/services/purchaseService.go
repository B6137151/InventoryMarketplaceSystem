package services

import (
	"fmt"
	"time"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseService interface {
	MakePurchase(request dtos.PurchaseCreateDTO) (dtos.OrderResponseDTO, error)
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id uuid.UUID) error
}

type purchaseService struct {
	orderRepo            repositories.OrderRepository
	orderDetailRepo      repositories.OrderDetailRepository
	productVariantRepo   repositories.ProductVariantRepository
	productRepo          repositories.ProductRepository
	salesRoundDetailRepo repositories.SalesRoundDetailRepository
}

// NewPurchaseService creates a new instance of PurchaseService
func NewPurchaseService(
	orderRepo repositories.OrderRepository,
	orderDetailRepo repositories.OrderDetailRepository,
	productVariantRepo repositories.ProductVariantRepository,
	productRepo repositories.ProductRepository,
	salesRoundDetailRepo repositories.SalesRoundDetailRepository,
) PurchaseService {
	return &purchaseService{
		orderRepo:            orderRepo,
		orderDetailRepo:      orderDetailRepo,
		productVariantRepo:   productVariantRepo,
		productRepo:          productRepo,
		salesRoundDetailRepo: salesRoundDetailRepo,
	}
}

func (s *purchaseService) MakePurchase(request dtos.PurchaseCreateDTO) (dtos.OrderResponseDTO, error) {
	// Auto-generate order code
	orderCode := fmt.Sprintf("ORDER-%s", uuid.New().String())

	// Initialize total price
	totalPrice := 0.0

	// Check and adjust stock and calculate total price
	for _, item := range request.Items {
		productVariant, err := s.productVariantRepo.GetProductVariantByID(item.VariantID)
		if err != nil {
			return dtos.OrderResponseDTO{}, err
		}

		product, err := s.productRepo.GetProductByID(productVariant.ProductID)
		if err != nil {
			return dtos.OrderResponseDTO{}, err
		}

		if product.Stock < item.Quantity {
			return dtos.OrderResponseDTO{}, fmt.Errorf("not enough stock")
		}

		// Check sales round detail for limits
		salesRoundDetail, err := s.salesRoundDetailRepo.GetSalesRoundDetailByRoundIDAndVariantID(request.RoundID, item.VariantID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return dtos.OrderResponseDTO{}, fmt.Errorf("sales round detail not found")
			}
			return dtos.OrderResponseDTO{}, err
		}

		if item.Quantity > salesRoundDetail.QuantityLimit {
			return dtos.OrderResponseDTO{}, fmt.Errorf("quantity exceeds sales round limit")
		}

		totalPrice += productVariant.Price * float64(item.Quantity)

		// Update sales round detail
		salesRoundDetail.Quantity -= item.Quantity
		if err := s.salesRoundDetailRepo.UpdateSalesRoundDetail(salesRoundDetail); err != nil {
			return dtos.OrderResponseDTO{}, err
		}
	}

	// Create order
	order := models.Order{
		CustomerID:      request.CustomerID,
		RoundID:         request.RoundID,
		OrderDate:       time.Now(),
		Status:          "ซื้อ สำเร็จ", // Automatically set status
		Code:            orderCode,     // Auto-generated order code
		TotalPrice:      totalPrice,    // Calculated total price
		DeliveryAddress: request.DeliveryAddress,
		PaymentSource:   request.PaymentSource,
	}

	if err := s.orderRepo.CreateOrder(&order); err != nil {
		return dtos.OrderResponseDTO{}, err
	}

	// Create order details and update stock
	for _, item := range request.Items {
		productVariant, err := s.productVariantRepo.GetProductVariantByID(item.VariantID)
		if err != nil {
			return dtos.OrderResponseDTO{}, err
		}

		product, err := s.productRepo.GetProductByID(productVariant.ProductID)
		if err != nil {
			return dtos.OrderResponseDTO{}, err
		}

		product.Stock -= item.Quantity
		if err := s.productRepo.UpdateProduct(product); err != nil {
			return dtos.OrderResponseDTO{}, err
		}

		orderDetail := models.OrderDetail{
			OrderID:    order.ID,
			VariantID:  item.VariantID,
			Quantity:   item.Quantity,
			Price:      productVariant.Price,
			TotalPrice: productVariant.Price * float64(item.Quantity),
		}

		if err := s.orderDetailRepo.CreateOrderDetail(&orderDetail); err != nil {
			return dtos.OrderResponseDTO{}, err
		}
	}

	response := dtos.OrderResponseDTO{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		RoundID:         order.RoundID,
		OrderDate:       order.OrderDate,
		Status:          order.Status,
		Code:            order.Code,
		TotalPrice:      order.TotalPrice,
		DeliveryAddress: order.DeliveryAddress,
		PaymentSource:   order.PaymentSource,
		CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

func (s *purchaseService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAllOrders()
}

func (s *purchaseService) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(id)
}

func (s *purchaseService) UpdateOrder(order *models.Order) error {
	return s.orderRepo.UpdateOrder(order)
}

func (s *purchaseService) DeleteOrder(id uuid.UUID) error {
	return s.orderRepo.DeleteOrder(id)
}
