package api

import "orderServiceGit/internal/core/services/DTO"

type OrdersResponse struct {
	Orders []DTO.OrderDTO `json:"values"`
}
