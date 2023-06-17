package DB_tests

import (
	"catalogServiceGit/internal/core/services/DTO"
	"github.com/google/uuid"
)

type DBTestInfo struct {
	Products []DTO.ProductDTO
}

func NewDBTestInfo(products []DTO.ProductDTO) *DBTestInfo {
	return &DBTestInfo{Products: products}
}

func (TestInfo *DBTestInfo) GetMinPrice() float64 {
	product := TestInfo.Products[0]
	for _, currentValue := range TestInfo.Products {
		if currentValue.Price < product.Price {
			product = currentValue
		}
	}
	return product.Price
}

func (TestInfo *DBTestInfo) GetMaxPrice() (price float64) {
	product := TestInfo.Products[0]
	for _, currentValue := range TestInfo.Products {
		if currentValue.Price > product.Price {
			product = currentValue
		}
	}
	return product.Price
}

func (TestInfo *DBTestInfo) GetMinRatting() float64 {
	product := TestInfo.Products[0]
	for _, currentValue := range TestInfo.Products {
		if currentValue.Ratting < product.Ratting {
			product = currentValue
		}
	}
	return product.Ratting
}

func (TestInfo *DBTestInfo) GetPrevMinRatting() float64 {
	product := TestInfo.Products[0]
	minRatting := TestInfo.GetMinRatting()
	for _, currentValue := range TestInfo.Products {
		if currentValue.Ratting < product.Ratting && currentValue.Ratting != minRatting {
			product = currentValue
		}
	}
	return product.Ratting
}

func (TestInfo *DBTestInfo) GetProductCountBySeller(seller string) int {
	var productCount int
	for _, currentValue := range TestInfo.Products {
		if currentValue.SalesMan == seller {
			productCount++
		}
	}
	return productCount
}

func (TestInfo *DBTestInfo) MakeNewProduct(productID string) *DTO.ProductDTO {
	return &DTO.ProductDTO{ID: productID, Ratting: 0, Price: 100, SalesMan: productID, Quantity: 43, Category: "laptop"}
}

var (
	testInfo = NewDBTestInfo(
		[]DTO.ProductDTO{
			{
				ID:           uuid.New().String(),
				Price:        10,
				Ratting:      5,
				Category:     "laptop",
				Quantity:     100,
				SalesMan:     uuid.New().String(),
				ReviewsCount: 10,
			},
			{
				ID:           uuid.New().String(),
				Price:        30,
				Ratting:      4,
				Category:     "laptop",
				Quantity:     100,
				SalesMan:     uuid.New().String(),
				ReviewsCount: 10,
			},
			{
				ID:           uuid.New().String(),
				Price:        20,
				Ratting:      3,
				Category:     "notebook",
				Quantity:     100,
				SalesMan:     uuid.New().String(),
				ReviewsCount: 20,
			},
			{
				ID:           uuid.New().String(),
				Price:        40,
				Ratting:      2,
				Category:     "table",
				Quantity:     100,
				SalesMan:     uuid.New().String(),
				ReviewsCount: 30,
			},
			{
				ID:           uuid.New().String(),
				Price:        50,
				Ratting:      1,
				Category:     "table",
				Quantity:     100,
				SalesMan:     uuid.New().String(),
				ReviewsCount: 40,
			},
		})
)

//func (TestInfo *DBTestInfo) GetProductCountByCategory(category string) int {
//	neededId := catalogService.DB.GetProductsOfCategoryQuery()
//	var productCount int
//	for _, currentValue := range TestInfo.Products {
//		if currentValue.CategoryID == context {
//			productCount++
//		}
//	}
//	return productCount
//}
