package api

type Reply struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

func NewReply(status string, data interface{}) *Reply {
	return &Reply{
		Status: status,
		Result: data,
	}
}
