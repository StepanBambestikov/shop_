package selection

type FilterType int

const (
	NoFilter FilterType = iota
	Category_
	Active
	SellerID
)

type SortingType int

const (
	NoSorting SortingType = iota
	Price
)

type SelectionDTO struct {
	Filter   FilterType
	Sorting  SortingType
	UserID   string
	Category string
	SellerID string
}
