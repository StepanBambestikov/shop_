package events

import (
	"catalogServiceGit/internal/core/services/DTO"
)

type CreateProductRequested struct {
	ProductDTO DTO.ProductDTO
	ResendCnt  uint64 `json:"resend_cnt" validate:"omitempty"`
}

type ChangeProductRequested struct {
	ProductDTO DTO.ProductDTO
	ResendCnt  uint64 `json:"resend_cnt" validate:"omitempty"`
}

type OrderProductRequested struct {
	ProductID string
	ResendCnt uint64 `json:"resend_cnt" validate:"omitempty"`
}

type RateProductRequested struct {
	ProductID  string
	NewRatting float64
	ResendCnt  uint64 `json:"resend_cnt" validate:"omitempty"`
}

type DeleteProductRequested struct {
	ProductID string
	ResendCnt uint64 `json:"resend_cnt" validate:"omitempty"`
}
