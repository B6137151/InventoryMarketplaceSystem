package repositories

import (
	"sync"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreRepository interface {
	CreateStore(store *models.Store) error
	GetAllStores() ([]models.Store, error)
	GetStoreByID(id uuid.UUID) (*models.Store, error)
	UpdateStore(store *models.Store) error
	DeleteStore(id uuid.UUID) error
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

func (r *storeRepository) GetAllStores() ([]models.Store, error) {
	var stores []models.Store
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.db.Preload("Products").Preload("Products.Category").Find(&stores).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) GetStoreByID(id uuid.UUID) (*models.Store, error) {
	var store models.Store
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.db.Preload("Products.Category").First(&store, "id = ?", id).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

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
