package dtos

type ProductUpdateStockDTO struct {
	Stock int `json:"stock" validate:"required"`
}
