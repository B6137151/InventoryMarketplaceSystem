package repositories

import (
	"strings"
	"sync"

	"log"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreRepository interface {
	CreateStore(store *models.Store) error
	GetAllStores(expand string) ([]models.Store, error) // Updated to include expand parameter
	GetStoreByID(id uuid.UUID) (*models.Store, error)
	UpdateStore(store *models.Store) error
	DeleteStore(id uuid.UUID) error
	GetStoresBySalesRoundID(salesRoundID uuid.UUID) ([]models.Store, error) // Added method

}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) CreateStore(store *models.Store) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.db.Create(store).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

func (r *storeRepository) GetAllStores(expand string) ([]models.Store, error) { // Updated signature
	var stores []models.Store
	query := r.db.Model(&models.Store{})

	// Check if expand contains "products" and preload if necessary
	if expand != "" && strings.Contains(expand, "products") {
		query = query.Preload("Products")
	}

	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := query.Find(&stores).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) GetStoreByID(id uuid.UUID) (*models.Store, error) {
	var store models.Store
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.db.Preload("Products").First(&store, "id = ?", id).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) UpdateStore(store *models.Store) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.db.Save(store).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

func (r *storeRepository) DeleteStore(id uuid.UUID) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.db.Delete(&models.Store{}, id).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

func (r *storeRepository) GetStoresBySalesRoundID(salesRoundID uuid.UUID) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Joins(`JOIN product ON product.store_id = store.id`).
		Joins(`JOIN "sales-round-detail" ON "sales-round-detail".product_id = product.id`).
		Where(`"sales-round-detail".round_id = ?`, salesRoundID).
		Find(&stores).Error
	if err != nil {
		log.Println("Error fetching stores by sales round ID:", err)
	}
	log.Println("Fetched stores by sales round ID:", stores)
	return stores, err
}
