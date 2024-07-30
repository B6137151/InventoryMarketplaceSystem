package controllers

//
//import (
//	"github.com/B6137151/InventoryMarketplaceSystem/internal/dtos"
//	"github.com/B6137151/InventoryMarketplaceSystem/internal/models"
//	"github.com/B6137151/InventoryMarketplaceSystem/internal/repositories"
//	"github.com/gofiber/fiber/v2"
//	"sync"
//)
//
//type PurchaseController struct {
//	purchaseRepository         repositories.PurchaseRepository
//	productVariantRepository   repositories.ProductVariantRepository
//	salesRoundDetailRepository repositories.SalesRoundDetailRepository
//}
//
//func NewPurchaseController(purchaseRepo repositories.PurchaseRepository, productVariantRepo repositories.ProductVariantRepository, salesRoundDetailRepo repositories.SalesRoundDetailRepository) *PurchaseController {
//	return &PurchaseController{
//		purchaseRepository:         purchaseRepo,
//		productVariantRepository:   productVariantRepo,
//		salesRoundDetailRepository: salesRoundDetailRepo,
//	}
//}
//
//// MakePurchase handles the purchase request
//func (c *PurchaseController) MakePurchase(ctx *fiber.Ctx) error {
//	dto := new(dtos.PurchaseDTO)
//	if err := ctx.BodyParser(dto); err != nil {
//		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "request body is not valid"})
//	}
//
//	var wg sync.WaitGroup
//	errChan := make(chan error, 1)
//
//	var orderDetails []models.OrderDetail
//	var totalAmount float64
//
//	for _, itemDTO := range dto.Items {
//		wg.Add(1)
//		go func(item dtos.PurchaseItemDTO) {
//			defer wg.Done()
//
//			salesRoundDetail, err := c.salesRoundDetailRepository.GetSalesRoundDetailByRoundAndVariantID(dto.RoundID, item.ProductVariantID)
//			if err != nil {
//				errChan <- err
//				return
//			}
//
//			if salesRoundDetail.Remaining < item.Quantity {
//				errChan <- fiber.NewError(fiber.StatusBadRequest, "Quantity exceeds available stock for the sales round")
//				return
//			}
//
//			productVariant, err := c.productVariantRepository.GetProductVariantByID(item.ProductVariantID)
//			if err != nil {
//				errChan <- err
//				return
//			}
//
//			orderDetail := models.OrderDetail{
//				VariantID:  item.ProductVariantID,
//				Quantity:   item.Quantity,
//				Price:      productVariant.Price,
//				TotalPrice: productVariant.Price * float64(item.Quantity),
//			}
//
//			orderDetails = append(orderDetails, orderDetail)
//			totalAmount += orderDetail.TotalPrice
//
//			salesRoundDetail.Remaining -= item.Quantity
//			if err := c.salesRoundDetailRepository.UpdateSalesRoundDetail(salesRoundDetail); err != nil {
//				errChan <- err
//				return
//			}
//		}(itemDTO)
//	}
//
//	go func() {
//		wg.Wait()
//		close(errChan)
//	}()
//
//	if err := <-errChan; err != nil {
//		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	}
//
//	purchase := models.Purchase{
//		CustomerID:      dto.CustomerID,
//		RoundID:         dto.RoundID,
//		TotalPrice:      totalAmount,
//		DeliveryAddress: dto.DeliveryAddress,
//		PaymentSource:   dto.PaymentSource,
//		OrderDetails:    orderDetails,
//	}
//
//	if err := c.purchaseRepository.CreatePurchase(&purchase); err != nil {
//		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create purchase"})
//	}
//
//	response := dtos.OrderResponseDTO{
//		ID:              purchase.ID,
//		CustomerID:      purchase.CustomerID,
//		RoundID:         purchase.RoundID,
//		OrderDate:       purchase.CreatedAt,
//		Status:          "completed", // Adjust as needed
//		Code:            "",          // Generate a code if necessary
//		TotalPrice:      purchase.TotalPrice,
//		DeliveryAddress: purchase.DeliveryAddress,
//		PaymentSource:   purchase.PaymentSource,
//		CreatedAt:       purchase.CreatedAt.Format("2006-01-02 15:04:05"),
//		UpdatedAt:       purchase.UpdatedAt.Format("2006-01-02 15:04:05"),
//	}
//
//	return ctx.Status(fiber.StatusCreated).JSON(response)
//}
