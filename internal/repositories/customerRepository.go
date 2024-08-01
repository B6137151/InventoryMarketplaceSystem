package repositories

import (
	"fmt"
	"log"

	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(customer *models.Customer) error
	GetAllCustomers() ([]models.Customer, error)
	GetCustomerByID(id uuid.UUID) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	DeleteCustomer(id uuid.UUID) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) CreateCustomer(customer *models.Customer) error {
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in CreateCustomer: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Create(customer).Error
	}()

	return <-errChan
}

func (r *customerRepository) GetAllCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetAllCustomers: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Find(&customers).Error
	}()

	err := <-errChan
	return customers, err
}

func (r *customerRepository) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in GetCustomerByID: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.First(&customer, "id = ?", id).Error
	}()

	err := <-errChan
	return &customer, err
}

func (r *customerRepository) UpdateCustomer(customer *models.Customer) error {
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in UpdateCustomer: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Save(customer).Error
	}()

	return <-errChan
}

func (r *customerRepository) DeleteCustomer(id uuid.UUID) error {
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in DeleteCustomer: %v", r)
				errChan <- fmt.Errorf("internal server error")
			}
			close(errChan)
		}()
		errChan <- r.db.Delete(&models.Customer{}, "id = ?", id).Error
	}()

	return <-errChan
}
