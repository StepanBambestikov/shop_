package entities

type Error struct {
	Code    int    `json:"-" example:"500"`
	Message string `json:"message" example:"Unknown error"`
}

func (app *Error) Error() string {
	return app.Message
}

func (app *Error) ToReply() *ApiReply {
	return &ApiReply{
		Error: &ErrorMessage{
			Message: app.Message,
		},
		Message: "error",
	}
}

// NewError creates a new Error instance with an optional message
func NewError(code int, message string) *Error {
	e := &Error{
		Code:    code,
		Message: message,
	}
	return e
}
