package api

import "catalogServiceGit/internal/core/services/DTO"

type ProductsResponse struct {
	Products map[string]DTO.ProductDTO `json:"values"`
}
