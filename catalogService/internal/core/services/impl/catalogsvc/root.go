package catalogService

import (
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/core/selection"
	"catalogServiceGit/internal/core/services"
	"catalogServiceGit/internal/core/services/DTO"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresCatalogClient(cfg *core.PostgresConfig) (services.CatalogService, error) {
	internalDB, err := gorm.Open(postgres.Open(cfg.Initialize), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbManager := &PostgresCatalogService{innerDB: internalDB}
	addCryptoExtention(dbManager)
	err = dbManager.DeleteAllTables()
	if err != nil {
		panic(err)
	}
	dbManager.MakeAllTables()

	return dbManager, nil
}

func FOR_TESTING_NewPostgresCatalogClient(cfg *core.PostgresConfig) (*PostgresCatalogService, error) {
	internalDB, err := gorm.Open(postgres.Open(cfg.Initialize), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbManager := &PostgresCatalogService{innerDB: internalDB}
	addCryptoExtention(dbManager)
	err = dbManager.DeleteAllTables()
	if err != nil {
		panic(err)
	}
	dbManager.MakeAllTables()
	return dbManager, nil
}

func (dbManager *PostgresCatalogService) GetSeveralProducts(selectionParams selection.SelectionDTO) (productIDs []string, err error) {
	query := dbManager.innerDB.Model(&Product{})
	var productsID []string
	var product Product
	switch selectionParams.Filter {
	case selection.MinPrice:
		err = query.Order("price").Limit(1).Find(&product).Error
		if err != nil {
			return nil, err
		}
		productIDs = []string{product.ID}
		return
	case selection.MaxPrice:
		err = query.Order("price DESC").Limit(1).Find(&product).Error
		if err != nil {
			return nil, err
		}
		productIDs = []string{product.ID}
		return
	case selection.MinRatting:
		err = query.Order("ratting").Last(&product).Error
		if err != nil {
			return nil, err
		}
		productIDs = []string{product.ID}
		return
	case selection.Category_:
		query = dbManager.getProductsOfCategoryQuery(selectionParams.Category)
	case selection.SellerID:
		query = dbManager.getProductsOfSellerQuery(selectionParams.Category)
	}

	switch selectionParams.Sorting {
	case selection.Price:
		query = query.Order("price")
	case selection.Ratting:
		query = query.Order("ratting")
	case selection.AvailableQuantity:
		query = query.Order("quantity")
	case selection.ReviewsCount:
		query = query.Order("reviews_count")
	}
	query.Scan(&productsID)
	return productsID, nil
}

func (dbManager *PostgresCatalogService) CreateProduct(dto DTO.ProductDTO) (err error) {
	product, err := dbManager.makeProductFromProductDTO(&dto)
	if err != nil {
		return err
	}
	err = dbManager.innerDB.Create(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (dbManager *PostgresCatalogService) ChangeProduct(dto DTO.ProductDTO) (err error) {
	product, err := dbManager.makeProductFromProductDTO(&dto)
	if err != nil {
		return err
	}
	err = dbManager.innerDB.Model(&product).Select("*").Updates(product).Error
	return err
}

func (dbManager *PostgresCatalogService) RateProduct(dto DTO.ProductDTO) (err error) {
	err = dbManager.innerDB.Model(&Product{}).Select("ratting").Where("ID = ?", dto.ID).Updates(Product{Ratting: dto.Ratting}).Error
	return
}

func (dbManager *PostgresCatalogService) DeleteProduct(dto DTO.ProductDTO) (err error) {
	err = dbManager.innerDB.Delete(&Product{}, Product{ID: dto.ID}).Error
	return
}

func (dbManager *PostgresCatalogService) GetProductByID(dto DTO.ProductDTO) (*DTO.ProductDTO, error) {
	var product Product
	err := dbManager.innerDB.First(&product, Product{ID: dto.ID}).Error
	if err != nil {
		return nil, err
	}
	dtoOutput, err := dbManager.makeDTOFromProduct(&product)
	if err != nil {
		return nil, err
	}
	return &dtoOutput, nil
}

func (dbManager *PostgresCatalogService) DeleteAllTables() (err error) {
	err = dbManager.innerDB.Migrator().DropTable(&Product{}, &Category{})
	return err
}

func (dbManager *PostgresCatalogService) MakeAllTables() {
	err := dbManager.innerDB.AutoMigrate(&Product{}, &Category{})
	if err != nil {
		panic(err)
	}
	err = dbManager.fillCategoryTable()
	if err != nil {
		panic(err)
	}
}
