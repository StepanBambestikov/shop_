package coresvc

import (
	"catalogServiceGit/internal/core/entities/events"
	"catalogServiceGit/internal/core/services"
	"catalogServiceGit/internal/core/services/DTO"
)

type coreServiceImpl struct {
	catalogService services.CatalogService
}

func NewCoreServiceImpl(cs services.CatalogService) services.CoreService {
	return &coreServiceImpl{
		catalogService: cs,
	}
}

func (core *coreServiceImpl) CreateProduct(ev *events.CreateProductRequested) (err error) {
	err = core.catalogService.CreateProduct(ev.ProductDTO)
	return err
}

func (core *coreServiceImpl) ChangeProduct(ev *events.ChangeProductRequested) (err error) {
	err = core.catalogService.ChangeProduct(ev.ProductDTO)
	return err
}

func (core *coreServiceImpl) RateProduct(ev *events.RateProductRequested) (err error) {
	err = core.catalogService.RateProduct(DTO.ProductDTO{ID: ev.ProductID, Ratting: ev.NewRatting})
	return err
}

func (core *coreServiceImpl) DeleteProduct(ev *events.DeleteProductRequested) (err error) {
	err = core.catalogService.DeleteProduct(DTO.ProductDTO{ID: ev.ProductID})
	return err
}
