package orderService

import (
	"gorm.io/gorm"
	"orderServiceGit/internal/core/services/DTO"
)

func (dbManager *PostgresOrderService) getCategoryID(category string) (id int, err error) {
	var gettedCategory Category
	err = dbManager.innerDB.First(&gettedCategory, &Category{Value: category}).Error
	if err != nil {
		return 0, err
	}
	return gettedCategory.ID, nil
}

func (dbManager *PostgresOrderService) getCategoryByID(id int) (category string, err error) {
	var gettedCategory Category
	err = dbManager.innerDB.First(&gettedCategory, &Category{ID: id}).Error
	if err != nil {
		return "", err
	}
	return gettedCategory.Value, nil
}

func copyCommonFieldsDTO2Order(order *Order, dto *DTO.OrderDTO) {
	order.ID = dto.ID
	order.Active = dto.Active
	order.Price = dto.Price
	order.UserId = dto.UserId
	order.SalesMan = dto.SalesMan
	return
}

func copyCommonFieldsOrder2DTO(order *Order, dto *DTO.OrderDTO) {
	dto.ID = order.ID
	dto.Active = order.Active
	dto.Price = order.Price
	dto.UserId = order.UserId
	dto.SalesMan = order.SalesMan
	return
}

func (dbManager *PostgresOrderService) makeOrderFromDTO(dto *DTO.OrderDTO) (order Order, err error) {
	copyCommonFieldsDTO2Order(&order, dto)
	categoryID, err := dbManager.getCategoryID(dto.Category)
	if err != nil {
		return Order{}, err
	}
	order.Category = categoryID
	statusID, err := dbManager.getStatusID(dto.Status)
	order.Status = statusID
	if err != nil {
		return Order{}, err
	}
	return
}

func (dbManager *PostgresOrderService) makeDTOFromOrder(order *Order) (dto DTO.OrderDTO, err error) {
	copyCommonFieldsOrder2DTO(order, &dto)
	category, err := dbManager.getCategoryByID(order.Category)
	if err != nil {
		return dto, err
	}
	dto.Category = category
	status, err := dbManager.getStatusByID(order.Status)
	if err != nil {
		return dto, err
	}
	dto.Status = status
	return
}

func addCryptoExtention(db *PostgresOrderService) (err error) {
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

func (dbManager *PostgresOrderService) getOrdersOfSellerQuery(sellerID string) *gorm.DB {
	return dbManager.innerDB.Model(&Order{}).Where("salesMan = ?", sellerID)
}

func (dbManager *PostgresOrderService) getActiveOrders() *gorm.DB {
	return dbManager.innerDB.Model(&Order{}).Where("active is true")
}

func (dbManager *PostgresOrderService) getOrdersOfCategory(category string) *gorm.DB {
	categoryID, err := dbManager.getCategoryID(category)
	if err != nil {
		panic("no such category")
	}
	return dbManager.innerDB.Model(&Order{}).Where("category = ?", categoryID)
}

func (dbManager *PostgresOrderService) getStatusID(status string) (id int, err error) {
	var gettedStatus Status
	err = dbManager.innerDB.First(&gettedStatus, &Status{Value: status}).Error
	if err != nil {
		return 0, err
	}
	return gettedStatus.ID, nil
}

func (dbManager *PostgresOrderService) getStatusByID(id int) (status string, err error) {
	var gettedStatus Status
	err = dbManager.innerDB.First(&gettedStatus, &Status{ID: id}).Error
	if err != nil {
		return "", err
	}
	return gettedStatus.Value, nil
}

func (dbManager *PostgresOrderService) fillCategoryTable() (err error) {
	for _, currentCategory := range Categories {
		err = dbManager.innerDB.Create(&Category{Value: string(currentCategory)}).Error
		if err != nil {
			return
		}
	}
	return
}

func (dbManager *PostgresOrderService) fillStatusTable() (err error) {
	for _, currentStatus := range Statuses {
		err = dbManager.innerDB.Create(&Status{Value: string(currentStatus)}).Error
		if err != nil {
			return
		}
	}
	return
}
