package orderService

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"orderServiceGit/internal/core"
	"orderServiceGit/internal/core/selection"
	"orderServiceGit/internal/core/services"
	"orderServiceGit/internal/core/services/DTO"
)

func NewPostgresOrderClient(cfg *core.PostgresConfig) (services.OrderService, error) {
	internalDB, err := gorm.Open(postgres.Open("user="+cfg.User+" password="+cfg.Password+" dbname="+cfg.Dbname+" sslmode=disable "+"host="+cfg.Host+" port="+cfg.Port), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbManager := &PostgresOrderService{innerDB: internalDB}
	addCryptoExtention(dbManager)
	err = dbManager.DeleteAllTables()
	if err != nil {
		panic(err)
	}
	dbManager.MakeAllTables()

	return dbManager, nil
}

func FOR_TESTING_NewPostgresOrderClient(cfg *core.PostgresConfig) (*PostgresOrderService, error) {
	internalDB, err := gorm.Open(postgres.Open("" /*TODO мне лень делать сецчас это блин*/), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbManager := &PostgresOrderService{innerDB: internalDB}
	addCryptoExtention(dbManager)
	err = dbManager.DeleteAllTables()
	if err != nil {
		panic(err)
	}
	dbManager.MakeAllTables()
	return dbManager, nil
}

func (dbManager *PostgresOrderService) GetOrderByID(orderDTO DTO.OrderDTO) (dto DTO.OrderDTO, err error) {
	var order Order
	err = dbManager.innerDB.First(&order, Order{ID: orderDTO.ID}).Error
	if err != nil {
		return dto, err
	}
	dto, err = dbManager.makeDTOFromOrder(&order)
	return
}

func (dbManager *PostgresOrderService) DeleteOrder(orderDTO DTO.OrderDTO) (err error) {
	dbManager.innerDB.Delete(&Order{}, Order{ID: orderDTO.ID})
	return
}

func (dbManager *PostgresOrderService) GetUserOrders(selectionParams *selection.SelectionDTO) (orderDTO []DTO.OrderDTO, err error) {
	query := dbManager.innerDB.Model(&Order{})
	switch selectionParams.Filter {
	case selection.SellerID:
		query = dbManager.getOrdersOfSellerQuery(selectionParams.SellerID)
	case selection.Category_:
		query = dbManager.getOrdersOfCategory(selectionParams.Category)
	case selection.Active:
		query = dbManager.getActiveOrders()
	}
	switch selectionParams.Sorting {
	case selection.Price:
		query = query.Order("price")
	}
	query = query.Where("user_id = ?", selectionParams.UserID)
	var orders []Order
	err = query.Scan(&orders).Error
	if err != nil {
		return nil, err
	}
	for _, val := range orders {
		currentDTO, err := dbManager.makeDTOFromOrder(&val)
		if err != nil {
			return nil, err
		}
		orderDTO = append(orderDTO, currentDTO)
	}
	return
}

func (dbManager *PostgresOrderService) CreateOrder(orderDTO DTO.OrderDTO) (err error) {
	order, err := dbManager.makeOrderFromDTO(&orderDTO)
	if err != nil {
		return err
	}
	err = dbManager.innerDB.Create(&order).Error
	return err
}

func (dbManager *PostgresOrderService) SetOrderStatus(dto DTO.OrderDTO) (err error) {
	statusID, err := dbManager.getStatusID(dto.Status)
	if err != nil {
		return err
	}
	err = dbManager.innerDB.Model(&Order{}).Where("ID = ?", dto.ID).Update("status", statusID).Error
	return err
}

func (dbManager *PostgresOrderService) DeleteAllTables() (err error) {
	err = dbManager.innerDB.Migrator().DropTable(&Order{}, &Category{}, &Status{})
	return err
}

func (dbManager *PostgresOrderService) MakeAllTables() {
	err := dbManager.innerDB.AutoMigrate(&Order{}, &Status{}, &Category{})
	if err != nil {
		panic(err)
	}
	err = dbManager.fillCategoryTable()
	if err != nil {
		panic(err)
	}
	err = dbManager.fillStatusTable()
	if err != nil {
		panic(err)
	}
}
