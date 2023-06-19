package services

import (
	"orderServiceGit/internal/core/entities/events"
	"orderServiceGit/internal/core/selection"
	"orderServiceGit/internal/core/services/DTO"
)

type OrderService interface {
	// return result
	GetUserOrders(dto *selection.SelectionDTO) ([]DTO.OrderDTO, error)
	GetOrderByID(orderDTO DTO.OrderDTO) (DTO.OrderDTO, error)
	//GetHealthz() int
	//GetMetrics()

	//create event
	CreateOrder(orderDTO DTO.OrderDTO) (err error)
	SetOrderStatus(orderDTO DTO.OrderDTO) (err error)
	DeleteOrder(orderDTO DTO.OrderDTO) (err error)
}

type GateService interface {
	// return result
	GetUserOrders(dto *selection.SelectionDTO) ([]DTO.OrderDTO, error)
	GetOrderByID(orderDTO DTO.OrderDTO) (DTO.OrderDTO, error)
	CreateOrder(orderDTO DTO.OrderDTO) (err error)
	SetOrderStatus(orderDTO DTO.OrderDTO) (err error)
	DeleteOrder(orderDTO DTO.OrderDTO) (err error)
}

type CoreService interface {
	CreateOrder(*events.CreateOrderRequested) (err error)
	SetOrderStatus(*events.SetOrderStatusRequested) (err error)
	DeleteOrder(*events.DeleteOrderRequested) (err error)
}
