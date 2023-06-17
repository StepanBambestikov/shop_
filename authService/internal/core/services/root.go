package services

import (
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities/api"
)

type AuthService interface {
	Authorize() (*api.LoginResponse, *entities.Error, error)
}

type UserService interface {
	Login(login string, password string) (api.LoginResponse, error)
	RefreshToken(refreshToken string) (api.LoginResponse, error)
	Logout(refreshToken string) error
	Signup(request api.SignupRequest) (api.SignupResponse, error)
	SetPassword(user *entities.User, password string) error
	DeleteUser(user *entities.User) error
	DeleteUserById(userId string) error
	GetUserById(userId string) (*entities.User, error)
	SetRoleById(userId string, role entities.UserRole) error
}
