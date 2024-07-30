package controllers

import (
	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"runtime"
	"sync"
)

type CustomerController interface {
	CreateCustomer(c *fiber.Ctx) error
	GetAllCustomers(c *fiber.Ctx) error
	UpdateCustomer(c *fiber.Ctx) error
	DeleteCustomer(c *fiber.Ctx) error
}

type customerController struct {
	customerRepository repositories.CustomerRepository
}

func NewCustomerController(customerRepository repositories.CustomerRepository) CustomerController {
	return &customerController{customerRepository: customerRepository}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Create a new customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param customer body dtos.CustomerCreateDTO true "Customer"
// @Success 201 {object} dtos.CustomerResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /customers [post]
func (h *customerController) CreateCustomer(c *fiber.Ctx) error {
	dto := new(dtos.CustomerCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	customer := models.Customer{Name: dto.Name, Email: dto.Email}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.customerRepository.CreateCustomer(&customer)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create customer"})
	}

	response := dtos.CustomerResponseDTO{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllCustomers godoc
// @Summary Get all customers
// @Description Get all customers
// @Tags Customers
// @Accept json
// @Produce json
// @Success 200 {array} dtos.CustomerResponseDTO
// @Failure 500 {object} fiber.Map
// @Router /customers [get]
func (h *customerController) GetAllCustomers(c *fiber.Ctx) error {
	var customers []models.Customer

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		customers, err = h.customerRepository.GetAllCustomers()
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve customers"})
	}

	var customerResponses []dtos.CustomerResponseDTO
	for _, customer := range customers {
		customerResponses = append(customerResponses, dtos.CustomerResponseDTO{
			ID:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(customerResponses)
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update a customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body dtos.CustomerUpdateDTO true "Customer"
// @Success 200 {object} dtos.CustomerResponseDTO
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /customers/{id} [put]
func (h *customerController) UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	dto := new(dtos.CustomerUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
	}

	var customer *models.Customer

	// Use a goroutine to fetch customer with synchronization
	errChan := make(chan error, 1)
	done := make(chan bool)
	go func() {
		defer close(errChan)
		customer, err = h.customerRepository.GetCustomerByID(uuid)
		if err != nil {
			errChan <- err
			return
		}
		done <- true // signal completion of fetch
	}()

	// Wait for the fetch to complete or an error
	select {
	case <-done:
		// Continue with update if fetch was successful
		customer.Name = dto.Name
		customer.Email = dto.Email

		if updateErr := h.customerRepository.UpdateCustomer(customer); updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update customer"})
		}
	case err := <-errChan:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found", "details": err.Error()})
	}

	response := dtos.CustomerResponseDTO{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.JSON(response)
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer
// @Tags Customers
// @Param id path string true "Customer ID"
// @Success 204
// @Failure 500 {object} fiber.Map
// @Router /customers/{id} [delete]
func (h *customerController) DeleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID format"})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		errChan <- h.customerRepository.DeleteCustomer(uuid)
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete customer"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func init() {
	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
