package DTO

type ProductDTO struct {
	ID           string
	Price        float64
	Ratting      float64
	Quantity     int
	Category     string
	SalesMan     string
	ReviewsCount int
}

func (u ProductDTO) Equal(other interface{}) bool {
	o, ok := other.(ProductDTO)
	if !ok {
		return false
	}
	return u.ID == o.ID &&
		u.Price == o.Price &&
		u.Ratting == o.Ratting &&
		u.Quantity == o.Quantity &&
		u.Category == o.Category &&
		u.SalesMan == o.SalesMan &&
		u.ReviewsCount == o.ReviewsCount
}
