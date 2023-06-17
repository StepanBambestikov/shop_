package events

import (
	"orderServiceGit/internal/core/services/DTO"
)

type CreateOrderRequested struct {
	OrderDTO  DTO.OrderDTO
	ResendCnt uint64 `json:"resend_cnt" validate:"omitempty"`
}

type SetOrderStatusRequested struct {
	OrderID   string
	NewStatus string
	ResendCnt uint64 `json:"resend_cnt" validate:"omitempty"`
}

type DeleteOrderRequested struct {
	OrderID   string
	ResendCnt uint64 `json:"resend_cnt" validate:"omitempty"`
}
