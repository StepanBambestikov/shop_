package signsvc

import (
	"errors"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities/api"
	"gitea.teneshag.ru/gigabit/goauth/internal/integrations"
	"gitea.teneshag.ru/gigabit/goauth/internal/util"
	"net/http"
)

type AuthStrategy interface {
	Authorize(context *authContext) (*api.LoginResponse, *entities.Error, error)
}

type LoginAuthStrategy struct{}

func (s LoginAuthStrategy) Authorize(context *authContext) (*api.LoginResponse, *entities.Error, error) {
	response, err := context.userService.Login(util.FromPointer(context.login), util.FromPointer(context.password))
	if err != nil && errors.Is(err, &integrations.ErrUnauthorized{}) {
		return nil, entities.NewError(http.StatusUnauthorized, "invalid credentials"), err
	} else if err != nil {
		return nil, entities.NewError(http.StatusInternalServerError, "intrnal error"), err
	}
	return &response, nil, nil
}

type TokenAuthStrategy struct{}

func (s TokenAuthStrategy) Authorize(context *authContext) (*api.LoginResponse, *entities.Error, error) {
	response, err := context.userService.RefreshToken(util.FromPointer(context.refreshToken))
	if err != nil && errors.Is(err, &integrations.ErrUnauthorized{}) {
		return nil, entities.NewError(http.StatusUnauthorized, "token is expired or revoked"), err
	} else if err != nil {
		return nil, entities.NewError(http.StatusInternalServerError, "intrnal error"), err
	}
	return &response, nil, nil

}
