package dtos

// SalesRoundDetailUpdateQuantityDTO is used for updating the quantity of a sales round detail
type SalesRoundDetailUpdateQuantityDTO struct {
	Quantity  int `json:"quantity" validate:"required"`
	Remaining int `json:"remaining" validate:"required"`
}
