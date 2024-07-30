package repositories

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetailRepository interface {
	CreateOrderDetail(orderDetail *models.OrderDetail) error
	GetAllOrderDetails() ([]models.OrderDetail, error)
	GetOrderDetailByID(id uuid.UUID) (*models.OrderDetail, error)
	UpdateOrderDetail(orderDetail *models.OrderDetail) error
	DeleteOrderDetail(id uuid.UUID) error
}

type orderDetailRepository struct {
	db *gorm.DB
}

func NewOrderDetailRepository(db *gorm.DB) OrderDetailRepository {
	return &orderDetailRepository{db: db}
}

func (r *orderDetailRepository) CreateOrderDetail(orderDetail *models.OrderDetail) error {
	return r.db.Create(orderDetail).Error
}

func (r *orderDetailRepository) GetAllOrderDetails() ([]models.OrderDetail, error) {
	var orderDetails []models.OrderDetail
	err := r.db.Find(&orderDetails).Error
	return orderDetails, err
}

func (r *orderDetailRepository) GetOrderDetailByID(id uuid.UUID) (*models.OrderDetail, error) {
	var orderDetail models.OrderDetail
	err := r.db.First(&orderDetail, "detail_id = ?", id).Error
	return &orderDetail, err
}

func (r *orderDetailRepository) UpdateOrderDetail(orderDetail *models.OrderDetail) error {
	return r.db.Save(orderDetail).Error
}

func (r *orderDetailRepository) DeleteOrderDetail(id uuid.UUID) error {
	return r.db.Delete(&models.OrderDetail{}, "detail_id = ?", id).Error
}
