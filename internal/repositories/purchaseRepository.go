package repositories

//import (
//	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
//	"github.com/google/uuid"
//	"gorm.io/gorm"
//)
//
//type PurchaseRepository interface {
//	CreatePurchase(purchase *models.Purchase) error
//	GetAllPurchases() ([]models.Purchase, error)
//	GetPurchaseByID(id uuid.UUID) (*models.Purchase, error)
//	UpdatePurchase(purchase *models.Purchase) error
//	DeletePurchase(id uuid.UUID) error
//}
//
//type purchaseRepository struct {
//	db *gorm.DB
//}
//
//func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
//	return &purchaseRepository{db: db}
//}
//
//func (r *purchaseRepository) CreatePurchase(purchase *models.Purchase) error {
//	return r.db.Create(purchase).Error
//}
//
//func (r *purchaseRepository) GetAllPurchases() ([]models.Purchase, error) {
//	var purchases []models.Purchase
//	err := r.db.Preload("OrderDetails").Find(&purchases).Error
//	return purchases, err
//}
//
//func (r *purchaseRepository) GetPurchaseByID(id uuid.UUID) (*models.Purchase, error) {
//	var purchase models.Purchase
//	err := r.db.Preload("OrderDetails").First(&purchase, "id = ?", id).Error
//	return &purchase, err
//}
//
//func (r *purchaseRepository) UpdatePurchase(purchase *models.Purchase) error {
//	return r.db.Save(purchase).Error
//}
//
//func (r *purchaseRepository) DeletePurchase(id uuid.UUID) error {
//	return r.db.Delete(&models.Purchase{}, "id = ?", id).Error
//}
