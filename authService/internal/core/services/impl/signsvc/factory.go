package signsvc

import (
	"errors"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities/api"
	"gitea.teneshag.ru/gigabit/goauth/internal/util"
)

var ErrInvalidBody = errors.New("invalid body: login-password or refresh token required")

func NewStrategy(request *api.LoginRequest) (AuthStrategy, error) {
	if request.Login != nil && request.Password != nil &&
		util.FromPointer(request.Login) != "" && util.FromPointer(request.Password) != "" {
		return &LoginAuthStrategy{}, nil
	} else if request.RefreshToken != nil && util.FromPointer(request.RefreshToken) != "" {
		return &TokenAuthStrategy{}, nil
	}
	return nil, ErrInvalidBody
}
