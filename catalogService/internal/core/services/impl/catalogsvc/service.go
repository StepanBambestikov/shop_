package catalogService

import (
	"catalogServiceGit/internal/core/services/DTO"
	"gorm.io/gorm"
)

func copyCommonFieldsDTO2Product(product *Product, dto *DTO.ProductDTO) {
	product.ID = dto.ID
	product.Price = dto.Price
	product.Ratting = dto.Ratting
	product.Quantity = dto.Quantity
	product.ReviewsCount = dto.ReviewsCount
	product.SalesMan = dto.SalesMan
	return
}

func copyCommonFieldsProduct2DTO(product *Product, dto *DTO.ProductDTO) {
	dto.ID = product.ID
	dto.Price = product.Price
	dto.Ratting = product.Ratting
	dto.Quantity = product.Quantity
	dto.ReviewsCount = product.ReviewsCount
	dto.SalesMan = product.SalesMan
	return
}

func (dbManager *PostgresCatalogService) makeProductFromProductDTO(dto *DTO.ProductDTO) (product Product, err error) {
	copyCommonFieldsDTO2Product(&product, dto)
	product.CategoryID, err = dbManager.getCategoryID(dto.Category)
	return
}

func (dbManager *PostgresCatalogService) makeDTOFromProduct(product *Product) (dto DTO.ProductDTO, err error) {
	copyCommonFieldsProduct2DTO(product, &dto)
	dto.Category, err = dbManager.getCategoryByID(product.CategoryID)
	return
}

func (dbManager *PostgresCatalogService) getProductsOfCategoryQuery(category string) *gorm.DB {
	var categoryID int
	dbManager.innerDB.Model(&Category{}).Where("value = ?", category).Pluck("categoryid", &categoryID)
	return dbManager.innerDB.Table("products").Where("categoryid = ?", categoryID)
}

func (dbManager *PostgresCatalogService) getProductsOfSellerQuery(sellerID string) *gorm.DB {
	return dbManager.innerDB.Model(&Product{}).Where("salesMan = ?", sellerID)
}

func (dbManager *PostgresCatalogService) getCategoryID(category string) (id int, err error) {
	var gettedCategory Category
	err = dbManager.innerDB.Model(&Category{}).First(&Category{Value: category}).Scan(&gettedCategory).Error
	if err != nil {
		return 0, err
	}
	return gettedCategory.CategoryID, nil
}

func (dbManager *PostgresCatalogService) getCategoryByID(categoryID int) (category string, err error) {
	var gettedCategory Category
	err = dbManager.innerDB.First(&gettedCategory, Category{CategoryID: categoryID}).Error
	if err != nil {
		return "", err
	}
	return gettedCategory.Value, nil
}

func addCryptoExtention(db *PostgresCatalogService) (err error) {
	var count int64
	err = db.innerDB.Raw("SELECT COUNT(*) FROM pg_extension WHERE extname = 'uuid-ossp'").Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		err = db.innerDB.Exec("CREATE EXTENSION \"uuid-ossp\";").Error
	}
	return
}

func (dbManager *PostgresCatalogService) fillCategoryTable() (err error) {
	for _, currentCategory := range Categories {
		err = dbManager.innerDB.Create(&Category{Value: string(currentCategory)}).Error
		if err != nil {
			return
		}
	}
	return
}
