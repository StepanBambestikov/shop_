package orderService

type StatusTable string

var (
	Statuses = []StatusTable{
		"Sending", "Delivering", "Delivered",
	}
)
