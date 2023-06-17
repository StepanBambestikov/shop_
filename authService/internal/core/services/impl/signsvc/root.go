package signsvc

import (
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities/api"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/services"
)

type authorizer struct {
	context  *authContext
	strategy AuthStrategy
}

func NewAuthorizer(login, password, refreshToken *string, strategy AuthStrategy, userService services.UserService) services.AuthService {
	return authorizer{
		context: &authContext{
			login:        login,
			password:     password,
			refreshToken: refreshToken,
			userService:  userService,
		},
		strategy: strategy,
	}
}

func (a authorizer) Authorize() (*api.LoginResponse, *entities.Error, error) {
	resp, outError, originError := a.strategy.Authorize(a.context)
	if originError != nil {
		return resp, outError, originError
	}
	return resp, nil, nil
}
