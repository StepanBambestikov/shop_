package api

import (
	"orderServiceGit/internal/core/selection"
	"orderServiceGit/internal/core/services/DTO"
)

type CreateOrderRequest struct {
	OrderDTO DTO.OrderDTO
}

type SetOrderStatusRequest struct {
	OrderID   string
	NewStatus string
}

type DeleteOrderRequest struct {
	OrderID string
}

type GetUserOrdersRequest struct {
	SelectionDTO selection.SelectionDTO
}

type GetOrderInfoRequest struct {
	OrderID string
}
