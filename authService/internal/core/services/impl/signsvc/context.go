package signsvc

import "gitea.teneshag.ru/gigabit/goauth/internal/core/services"

type authContext struct {
	login        *string
	password     *string
	refreshToken *string
	userService  services.UserService
}
