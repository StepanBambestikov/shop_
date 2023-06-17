package entities

type ApiReply struct {
	Data    interface{}   `json:"data,omitempty"`
	Error   *ErrorMessage `json:"error"`
	Message string        `json:"message,omitempty"`
} // @name Response

type ErrorMessage struct {
	Message string `json:"message"`
} // @name ErrorMessage

func ReplyError(message string, code int) ApiReply {
	return ApiReply{
		Error: &ErrorMessage{
			Message: message,
		},
		Message: "Error",
	}
}
