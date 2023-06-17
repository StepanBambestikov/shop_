package DB_tests

import (
	"github.com/stretchr/testify/assert"
	"orderServiceGit/internal/core"
	"orderServiceGit/internal/core/selection"
	"orderServiceGit/internal/core/services/DTO"
	orderService "orderServiceGit/internal/core/services/impl/ordersvc"
	"testing"
)

var (
	databaseInit = "user=postgres password=Qqqwwweee12321 dbname=postgres sslmode=disable"
)

func MakeAndPrepareDB(testInfo *DBTestInfo) (DB *orderService.PostgresOrderService, err error) {
	DB, err = orderService.FOR_TESTING_NewPostgresOrderClient(&core.PostgresConfig{Initialize: databaseInit})
	for _, order := range testInfo.Orders {
		err := DB.CreateOrder(order)
		if err != nil {
			return nil, err
		}
	}
	return DB, nil
}

func TestDatabase_SuccessGettingUserOrders(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetUserOrders(&selection.SelectionDTO{UserID: userID1})
	assert.NoError(t, err)
	assert.NotEmpty(t, products)
	for _, val := range products {
		assert.Equal(t, val.UserId, userID1)
	}
}

func TestDatabase_SuccessGettingActiveUserOrders(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetUserOrders(&selection.SelectionDTO{UserID: userID1, Filter: selection.Active})
	assert.NoError(t, err)
	assert.NotEmpty(t, products)
	for _, val := range products {
		assert.Equal(t, val.UserId, userID1)
		assert.True(t, val.Active)
	}
}

func TestDatabase_SuccessDeletingOrder(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetUserOrders(&selection.SelectionDTO{UserID: userID1, Filter: selection.Active})
	id := products[0].ID
	err = DB.DeleteOrder(DTO.OrderDTO{ID: id})
	assert.NoError(t, err)
	newProducts, err := DB.GetUserOrders(&selection.SelectionDTO{UserID: userID1, Filter: selection.Active})
	assert.NoError(t, err)
	assert.NotEmpty(t, newProducts)
	assert.Equal(t, len(products)-1, len(newProducts))
}

func TestDatabase_SuccessSetStatus(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetUserOrders(&selection.SelectionDTO{UserID: userID1, Filter: selection.Active})
	id := products[0].ID
	err = DB.SetOrderStatus(DTO.OrderDTO{ID: id, Status: "Delivering"})
	assert.NoError(t, err)
	product, err := DB.GetOrderByID(DTO.OrderDTO{ID: id})
	assert.NotNil(t, product)
	assert.Equal(t, "Delivering", product.Status)
}

//TODO добавить проверку на невозможность добавить несуществующий статус и категорию!!!!
