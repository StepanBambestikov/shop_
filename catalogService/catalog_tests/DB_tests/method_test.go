package DB_tests

import (
	"catalogServiceGit/internal/core/selection"
	"catalogServiceGit/internal/core/services/DTO"
	catalogService "catalogServiceGit/internal/core/services/impl/catalogsvc"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	databaseInit = "user=postgres password=Qqqwwweee12321 dbname=postgres sslmode=disable"
)

func MakeAndPrepareDB(testInfo *DBTestInfo) (DB *catalogService.PostgresCatalogService, err error) {
	DB, err = catalogService.FOR_TESTING_NewPostgresCatalogClient(databaseInit)
	for _, product := range testInfo.Products {
		err = DB.CreateProduct(product)
		if err != nil {
			return nil, err
		}
	}
	return DB, nil
}

func TestDatabase_SuccessMinPriceGetting(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetSeveralProducts(selection.SelectionDTO{
		Filter: selection.MinPrice,
	})
	assert.NoError(t, err)
	product, err := DB.GetProductByID(DTO.ProductDTO{ID: products[0]})
	assert.NoError(t, err)
	assert.Equal(t, testInfo.GetMinPrice(), product.Price)
	return

}

func TestDatabase_SuccessMaxPriceGetting(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	//MaxPrice check
	products, err := DB.GetSeveralProducts(selection.SelectionDTO{
		Filter: selection.MaxPrice,
	})
	assert.NoError(t, err)
	product, err := DB.GetProductByID(DTO.ProductDTO{ID: products[0]})
	assert.NoError(t, err)
	assert.Equal(t, testInfo.GetMaxPrice(), product.Price)
	return
}

func TestDatabase_SuccessMinRattingGetting(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetSeveralProducts(selection.SelectionDTO{
		Filter: selection.MinRatting,
	})
	assert.NoError(t, err)
	product, err := DB.GetProductByID(DTO.ProductDTO{ID: products[0]})
	assert.NoError(t, err)
	assert.Equal(t, testInfo.GetMinRatting(), product.Ratting)
	return
}

func TestDatabase_SuccessDeleting(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetSeveralProducts(selection.SelectionDTO{
		Filter: selection.MinRatting,
	})
	assert.NoError(t, err)
	product, err := DB.GetProductByID(DTO.ProductDTO{ID: products[0]})
	assert.NoError(t, err)
	assert.Equal(t, testInfo.GetMinRatting(), product.Ratting)

	err = DB.DeleteProduct(DTO.ProductDTO{ID: products[0]})
	assert.NoError(t, err)

	//New MinRatting check
	products, err = DB.GetSeveralProducts(selection.SelectionDTO{
		Filter: selection.MinRatting,
	})
	assert.NoError(t, err)
	product, err = DB.GetProductByID(DTO.ProductDTO{ID: products[0]})
	assert.NoError(t, err)
	assert.Equal(t, testInfo.GetPrevMinRatting(), product.Ratting)
	return
}

func TestDatabase_SuccessChangingProduct(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetSeveralProducts(selection.SelectionDTO{Filter: selection.MinPrice})
	assert.NoError(t, err)

	newProduct := *testInfo.MakeNewProduct(products[0])
	err = DB.ChangeProduct(newProduct)
	assert.NoError(t, err)
	//product, err := DB.GetProductByID(DTO.ProductDTO{ID: products[0]})
	//assert.NoError(t, err)
	//assert.Equal(t, newProduct, product) TODO Equaler for DTO.ProductDTO doesnt work
	return
}

func TestDatabase_SuccessRateProduct(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)

	products, err := DB.GetSeveralProducts(selection.SelectionDTO{Filter: selection.MinPrice})
	assert.NoError(t, err)

	newRatting := 43.0
	err = DB.RateProduct(DTO.ProductDTO{ID: products[0], Ratting: newRatting})
	assert.NoError(t, err)
	products, err = DB.GetSeveralProducts(selection.SelectionDTO{Filter: selection.MinPrice})
	assert.NoError(t, err)
	product, err := DB.GetProductByID(DTO.ProductDTO{ID: products[0]})
	assert.NoError(t, err)
	assert.Equal(t, newRatting, product.Ratting)
	return
}

func TestDatabase_SuccessCategory(t *testing.T) {
	DB, err := MakeAndPrepareDB(testInfo)
	assert.NoError(t, err)
	productsIDs, err := DB.GetSeveralProducts(selection.SelectionDTO{Filter: selection.Category_, Category: "laptop"})
	assert.NoError(t, err)
	assert.NotEmpty(t, productsIDs)
	for _, currentProductID := range productsIDs {
		currentProduct, err := DB.GetProductByID(DTO.ProductDTO{ID: currentProductID})
		assert.NoError(t, err)
		assert.Equal(t, currentProduct.Category, "laptop")
	}
	return
}

//TODO add test on selection seller!

//TODO добавить проверку на невозможность добавить несуществующий статус и категорию!!!!
