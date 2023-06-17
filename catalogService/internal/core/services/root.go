package services

import (
	"catalogServiceGit/internal/core/entities/events"
	"catalogServiceGit/internal/core/selection"
	"catalogServiceGit/internal/core/services/DTO"
)

type CatalogService interface {
	// return result
	GetSeveralProducts(dto selection.SelectionDTO) ([]string, error)
	//GetHealthz() int
	//GetMetrics()

	//create event
	CreateProduct(productDTO DTO.ProductDTO) (err error)
	ChangeProduct(productDTO DTO.ProductDTO) (err error)
	RateProduct(productDTO DTO.ProductDTO) (err error)
	DeleteProduct(productDTO DTO.ProductDTO) (err error)
}

type GateService interface {
	// return result
	GetSeveralProducts(selectionParams selection.SelectionDTO) ([]string, error)

	//create event
	CreateProduct(productDTO DTO.ProductDTO) (err error)
	ChangeProduct(productDTO DTO.ProductDTO) (err error)
	OrderProduct(productDTO DTO.ProductDTO) (err error)
	RateProduct(productDTO DTO.ProductDTO) (err error)
	DeleteProduct(productDTO DTO.ProductDTO) (err error)
}

type CoreService interface {
	CreateProduct(*events.CreateProductRequested) (err error)
	ChangeProduct(*events.ChangeProductRequested) (err error)
	RateProduct(*events.RateProductRequested) (err error)
	DeleteProduct(*events.DeleteProductRequested) (err error)
}
