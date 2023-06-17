package DB_tests

import (
	"github.com/google/uuid"
	"orderServiceGit/internal/core/services/DTO"
)

type DBTestInfo struct {
	Orders []DTO.OrderDTO
}

func NewDBTestInfo(orders []DTO.OrderDTO) *DBTestInfo {
	return &DBTestInfo{Orders: orders}
}

func (TestInfo *DBTestInfo) GetOrdersOfUser(userID string) []DTO.OrderDTO {
	var userOrders []DTO.OrderDTO
	for _, currentOrder := range TestInfo.Orders {
		if currentOrder.UserId == userID {
			userOrders = append(userOrders, currentOrder)
		}
	}
	return userOrders
}

func (TestInfo *DBTestInfo) GetActiveOrdersOfUser(userID string) []DTO.OrderDTO {
	var userOrders []DTO.OrderDTO
	for _, currentOrder := range TestInfo.Orders {
		if currentOrder.UserId == userID && currentOrder.Active {
			userOrders = append(userOrders, currentOrder)
		}
	}
	return userOrders
}

func (TestInfo *DBTestInfo) GetSellerOrdersOfUser(userID string, sellerID string) []DTO.OrderDTO {
	var userOrders []DTO.OrderDTO
	for _, currentOrder := range TestInfo.Orders {
		if currentOrder.UserId == userID && currentOrder.SalesMan == sellerID {
			userOrders = append(userOrders, currentOrder)
		}
	}
	return userOrders
}

var (
	userID1, userID2 = uuid.New().String(), uuid.New().String()
	testInfo         = NewDBTestInfo(
		[]DTO.OrderDTO{
			{
				Price:    10,
				Active:   true,
				UserId:   userID1,
				SalesMan: userID2,
				Category: "laptop",
				Status:   "Delivering",
			},
			{
				Price:    20,
				Active:   false,
				UserId:   userID1,
				SalesMan: userID2,
				Category: "laptop",
				Status:   "Delivering",
			},
			{
				Price:    40,
				Active:   true,
				UserId:   userID1,
				SalesMan: userID2,
				Category: "notebook",
				Status:   "Delivered",
			},
			{
				Price:    30,
				Active:   false,
				UserId:   userID2,
				SalesMan: userID1,
				Category: "notebook",
				Status:   "Sending",
			},
			{
				Price:    50,
				Active:   true,
				UserId:   userID2,
				SalesMan: userID1,
				Category: "notebook",
				Status:   "Sending",
			},
		})
)
