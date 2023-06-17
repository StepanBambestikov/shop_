package integrations

import "orderServiceGit/internal/core/entities"

type EventBus interface {
	Send(topic string, message any) error
	SendDelayed(topic string, message any, delay int32) error
	AddConsumer(
		name string,
		topic string,
		handler func([]byte) error,
		discardOnError bool,
	) error
	Close() error
}

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
