package orderService

type Order struct {
	ID       string  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Category int     `gorm:"foreingKey"`
	Active   bool    `gorm:"column:active"`
	Price    float64 `gorm:"column:price"`
	Status   int     `gorm:"foreingKey"`
	UserId   string  `gorm:"column:user_id;type:uuid"`
	SalesMan string  `gorm:"type:uuid"`
}

type Status struct {
	ID    int    `gorm:"primaryKey"`
	Value string `gorm:"column:value"`
}

type Category struct {
	ID    int    `gorm:"column:categoryid;primaryKey;auto_increment"`
	Value string `gorm:"column:value"`
}
