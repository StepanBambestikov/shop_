package integrations

import "gitea.teneshag.ru/gigabit/goauth/internal/core/entities"

type UserStorage interface {
	GetFromToken(accessToken string) (*entities.User, error)

	SignupUser(user *entities.User) error
	LoginUser(login string, password string) (*entities.User, string, string, error)
	LogoutUser(refreshToken string) error
	LogoutByID(userID string) error
	RefreshToken(refreshToken string) (string, string, error)

	SetPassword(user *entities.User, password string) error

	DeleteUser(user *entities.User) error

	GetUserById(userId string) (*entities.User, error)
	SetRole(userId *entities.User, role entities.UserRole) error
}
