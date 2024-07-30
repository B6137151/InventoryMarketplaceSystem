package repositories

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"gorm.io/gorm"
)

type OrderHistoryRepository interface {
	CreateOrderHistory(orderHistory *models.OrderHistory) error
	GetAllOrderHistories() ([]models.OrderHistory, error)
	GetOrderHistoryByID(id string) (*models.OrderHistory, error)
	UpdateOrderHistory(orderHistory *models.OrderHistory) error
	DeleteOrderHistory(id string) error
}

type orderHistoryRepository struct {
	db *gorm.DB
}

func NewOrderHistoryRepository(db *gorm.DB) OrderHistoryRepository {
	return &orderHistoryRepository{db: db}
}

func (r *orderHistoryRepository) CreateOrderHistory(orderHistory *models.OrderHistory) error {
	return r.db.Create(orderHistory).Error
}

func (r *orderHistoryRepository) GetAllOrderHistories() ([]models.OrderHistory, error) {
	var orderHistories []models.OrderHistory
	err := r.db.Find(&orderHistories).Error
	return orderHistories, err
}

func (r *orderHistoryRepository) GetOrderHistoryByID(id string) (*models.OrderHistory, error) {
	var orderHistory models.OrderHistory
	err := r.db.First(&orderHistory, id).Error
	return &orderHistory, err
}

func (r *orderHistoryRepository) UpdateOrderHistory(orderHistory *models.OrderHistory) error {
	return r.db.Save(orderHistory).Error
}

func (r *orderHistoryRepository) DeleteOrderHistory(id string) error {
	return r.db.Delete(&models.OrderHistory{}, id).Error
}
