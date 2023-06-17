package apisvc

import (
	"orderServiceGit/internal/core"
	"orderServiceGit/internal/core/entities/events"
	"orderServiceGit/internal/core/selection"
	"orderServiceGit/internal/core/services"
	"orderServiceGit/internal/core/services/DTO"
	"orderServiceGit/internal/integrations"
)

type apiServiceImpl struct {
	orderService services.OrderService
	eventBus     integrations.EventBus
}

func NewApiServiceImpl(cs services.OrderService, eb integrations.EventBus) services.GateService {
	return &apiServiceImpl{
		orderService: cs,
		eventBus:     eb,
	}
}

func (api *apiServiceImpl) GetUserOrders(dto *selection.SelectionDTO) ([]DTO.OrderDTO, error) {
	productDTOs, err := api.orderService.GetUserOrders(dto)
	if err != nil {
		return nil, err
	}
	return productDTOs, nil
}

func (api *apiServiceImpl) GetOrderByID(orderDTO DTO.OrderDTO) (orderOutput DTO.OrderDTO, err error) {
	orderOutput, err = api.orderService.GetOrderByID(orderDTO)
	return
}

func (api *apiServiceImpl) CreateOrder(orderDTO DTO.OrderDTO) (err error) {
	event := events.CreateOrderRequested{
		OrderDTO: orderDTO,
	}
	err = api.eventBus.Send(core.TopicCreateOrder, event)
	return err
}
func (api *apiServiceImpl) SetOrderStatus(orderDTO DTO.OrderDTO) (err error) {
	event := events.SetOrderStatusRequested{
		OrderID:   orderDTO.ID,
		NewStatus: orderDTO.Status,
	}
	err = api.eventBus.Send(core.TopicSetOrderStatus, event)
	return err
}
func (api *apiServiceImpl) DeleteOrder(orderDTO DTO.OrderDTO) (err error) {
	event := events.DeleteOrderRequested{
		OrderID: orderDTO.ID,
	}
	err = api.eventBus.Send(core.TopicDeleteOrder, event)
	return err
}
