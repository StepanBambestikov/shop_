package catalogService

type Product struct {
	ID           string  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Price        float64 `gorm:"column:price"`
	Ratting      float64 `gorm:"column:ratting"`
	Quantity     int     `gorm:"column:quantity"`
	CategoryID   int     `gorm:"column:categoryid;foreingKey"`
	SalesMan     string  `gorm:"type:uuid"`
	ReviewsCount int     `gorm:"column:reviews_count"`
}

type Category struct {
	CategoryID int `gorm:"column:categoryid;primaryKey;auto_increment"`
	Value      string
}
