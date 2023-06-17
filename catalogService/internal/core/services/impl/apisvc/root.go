package apisvc

import (
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/core/entities/events"
	"catalogServiceGit/internal/core/selection"
	"catalogServiceGit/internal/core/services"
	"catalogServiceGit/internal/core/services/DTO"
	"catalogServiceGit/internal/integrations"
)

type apiServiceImpl struct {
	counterService services.CatalogService
	eventBus       integrations.EventBus
}

func NewApiServiceImpl(cs services.CatalogService, eb integrations.EventBus) services.GateService {
	return &apiServiceImpl{
		counterService: cs,
		eventBus:       eb,
	}
}

func (api *apiServiceImpl) GetSeveralProducts(dto selection.SelectionDTO) ([]string, error) {
	productDTOs, err := api.counterService.GetSeveralProducts(dto)
	if err != nil {
		return nil, err
	}
	return productDTOs, nil
}

func (api *apiServiceImpl) CreateProduct(info DTO.ProductDTO) (err error) {
	event := events.CreateProductRequested{
		ProductDTO: info,
	}
	err = api.eventBus.Send(core.TopicCreateProduct, event)
	return err
}

func (api *apiServiceImpl) ChangeProduct(info DTO.ProductDTO) (err error) {
	event := events.ChangeProductRequested{
		ProductDTO: info,
	}
	err = api.eventBus.Send(core.TopicChangeProduct, event)
	return err
}

func (api *apiServiceImpl) OrderProduct(dto DTO.ProductDTO) (err error) {
	event := events.OrderProductRequested{
		ProductID: dto.ID,
	}
	err = api.eventBus.Send(core.TopicOrderProduct, event)
	return err
}

func (api *apiServiceImpl) RateProduct(dto DTO.ProductDTO) (err error) {
	event := events.RateProductRequested{
		ProductID:  dto.ID,
		NewRatting: dto.Ratting,
	}
	err = api.eventBus.Send(core.TopicRateProduct, event)
	return err
}

func (api *apiServiceImpl) DeleteProduct(dto DTO.ProductDTO) (err error) {
	event := events.DeleteProductRequested{
		ProductID: dto.ID,
	}
	err = api.eventBus.Send(core.TopicDeleteProduct, event)
	return err
}
