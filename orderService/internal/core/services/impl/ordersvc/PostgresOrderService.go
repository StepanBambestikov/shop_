package orderService

import (
	_ "github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresOrderService struct {
	innerDB *gorm.DB
}

//var (
//	databaseInit = "user=postgres password=qweasdzxc12321 dbname=postgres sslmode=disable"
//)
