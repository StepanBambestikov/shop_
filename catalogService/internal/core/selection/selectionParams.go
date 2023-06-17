package selection

type FilterType int

const (
	NoFilter FilterType = iota
	MinPrice
	MaxPrice
	Category_
	MinRatting
	MinSold
	SellerID
)

type SortingType int

const (
	NoSorting SortingType = iota
	Price
	Ratting
	AvailableQuantity
	ReviewsCount
)

type SelectionDTO struct {
	Filter   FilterType
	Sorting  SortingType
	Category string
	SellerID string
}
