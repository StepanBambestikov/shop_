package api

import (
	"catalogServiceGit/internal/core/selection"
	"catalogServiceGit/internal/core/services/DTO"
)

type GetAllProductsRequest struct {
	SelectionDTO selection.SelectionDTO
}

type CreateProductRequest struct {
	ProductDTO DTO.ProductDTO
}

type ChangeProductRequest struct {
	ProductDTO DTO.ProductDTO
}

type OrderProductRequest struct {
	ProductID string
}

type RateProductRequest struct {
	ProductID  string
	NewRatting float64
}

type DeleteProductRequest struct {
	ProductID string
}
