package integrations

type ErrConnectionFailed struct {
	Message string
}

func (e ErrConnectionFailed) Error() string {
	return e.Message
}

type ErrAlreadyExists struct{}

func (e ErrAlreadyExists) Error() string {
	return "already exists"
}

type ErrInvalidCredentials struct {
	Message string
}

func (e ErrInvalidCredentials) Error() string {
	return e.Message
}

type ErrUnauthorized struct{}

func (e ErrUnauthorized) Error() string {
	return "unauthorized"
}

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "not found"
}
