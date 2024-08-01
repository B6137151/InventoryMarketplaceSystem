package repositories

import (
	"fmt"
	"log"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id uuid.UUID) error
	GetRecognizedRevenue(roundID uuid.UUID) (float64, error)
	GetTotalOrders(roundID uuid.UUID) (int, error)
	GetTotalItemsOrdered(roundID uuid.UUID) (int, error)
	GetTotalItemsSold(roundID uuid.UUID) (int, error)
	GetOrdersByRoundID(roundID uuid.UUID) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *models.Order) error {
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in CreateOrder: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Create(order).Error
	}()

	return <-errChan
}

func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetAllOrders: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Find(&orders).Error
	}()

	err := <-errChan
	return orders, err
}

func (r *orderRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetOrderByID: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.First(&order, "id = ?", id).Error
	}()

	err := <-errChan
	return &order, err
}

func (r *orderRepository) UpdateOrder(order *models.Order) error {
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in UpdateOrder: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Save(order).Error
	}()

	return <-errChan
}

func (r *orderRepository) DeleteOrder(id uuid.UUID) error {
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in DeleteOrder: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Delete(&models.Order{}, "id = ?", id).Error
	}()

	return <-errChan
}

// New metrics methods
func (r *orderRepository) GetRecognizedRevenue(roundID uuid.UUID) (float64, error) {
	var totalRevenue float64
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetRecognizedRevenue: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Model(&models.Order{}).Where("round_id = ? AND status = ?", roundID, "paid").Select("SUM(total_price)").Scan(&totalRevenue).Error
	}()
	err := <-errChan
	return totalRevenue, err
}

func (r *orderRepository) GetTotalOrders(roundID uuid.UUID) (int, error) {
	var count int64
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetTotalOrders: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Model(&models.Order{}).Where("round_id = ?", roundID).Count(&count).Error
	}()
	err := <-errChan
	return int(count), err
}

func (r *orderRepository) GetTotalItemsOrdered(roundID uuid.UUID) (int, error) {
	var totalItems int64
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetTotalItemsOrdered: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Model(&models.OrderDetail{}).
			Joins("JOIN orders ON orders.id = order_details.order_id").
			Where("orders.round_id = ?", roundID).
			Select("SUM(order_details.quantity)").
			Scan(&totalItems).Error
	}()
	err := <-errChan
	return int(totalItems), err
}

func (r *orderRepository) GetTotalItemsSold(roundID uuid.UUID) (int, error) {
	var totalItems int64
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetTotalItemsSold: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Model(&models.OrderDetail{}).
			Joins("JOIN \"order\" ON \"order\".id = \"order-detail\".order_id").
			Where("\"order\".round_id = ? AND \"order\".status = ?", roundID, "paid").
			Select("SUM(\"order-detail\".quantity)").
			Scan(&totalItems).Error

	}()
	err := <-errChan
	return int(totalItems), err
}

func (r *orderRepository) GetOrdersByRoundID(roundID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetOrdersByRoundID: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Where("round_id = ?", roundID).Preload("OrderDetails").Find(&orders).Error
	}()
	err := <-errChan
	return orders, err
}
