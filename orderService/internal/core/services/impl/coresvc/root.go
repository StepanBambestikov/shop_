package coresvc

import (
	"orderServiceGit/internal/core/entities/events"
	"orderServiceGit/internal/core/services"
	"orderServiceGit/internal/core/services/DTO"
)

type coreServiceImpl struct {
	catalogService services.OrderService
}

func NewCoreServiceImpl(cs services.OrderService) services.CoreService {
	return &coreServiceImpl{
		catalogService: cs,
	}
}

func (core *coreServiceImpl) CreateOrder(ev *events.CreateOrderRequested) (err error) {
	err = core.catalogService.CreateOrder(ev.OrderDTO)
	return err
}

func (core *coreServiceImpl) SetOrderStatus(ev *events.SetOrderStatusRequested) (err error) {
	err = core.catalogService.SetOrderStatus(DTO.OrderDTO{Status: ev.NewStatus, ID: ev.OrderID})
	return err
}
func (core *coreServiceImpl) DeleteOrder(ev *events.DeleteOrderRequested) (err error) {
	err = core.catalogService.DeleteOrder(DTO.OrderDTO{ID: ev.OrderID})
	return err
}
