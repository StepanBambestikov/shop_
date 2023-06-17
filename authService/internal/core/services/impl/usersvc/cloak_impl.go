package usersvc

import (
	"gitea.teneshag.ru/gigabit/goauth/internal/integrations"
	"gitea.teneshag.ru/gigabit/goauth/internal/log"
	"gitea.teneshag.ru/gigabit/goauth/internal/util"

	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities/api"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/services"

	"github.com/pkg/errors"
)

type cloakUserService struct {
	userStorage integrations.UserStorage
}

func NewCloakUserService(userStorage integrations.UserStorage) services.UserService {
	svc := &cloakUserService{
		userStorage: userStorage,
	}
	return svc
}

func (u *cloakUserService) Login(login string, password string) (api.LoginResponse, error) {
	user, access, refresh, err := u.userStorage.LoginUser(login, password)
	if err != nil {
		return api.LoginResponse{}, errors.Wrap(err, "can't login")
	}
	log.Infof("User %s %s logged in", util.FromPointer(user.ID), user.Username)

	return api.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *cloakUserService) RefreshToken(refreshToken string) (api.LoginResponse, error) {
	access, refresh, err := u.userStorage.RefreshToken(refreshToken)
	if err != nil {
		return api.LoginResponse{}, err
	}

	return api.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *cloakUserService) Logout(refreshToken string) error {
	return u.userStorage.LogoutUser(refreshToken)
}

func (u *cloakUserService) Signup(request api.SignupRequest) (api.SignupResponse, error) {
	user := &entities.User{
		Username:      request.Username,
		FirstName:     util.ToPointer(request.FirstName),
		LastName:      util.ToPointer(request.LastName),
		Email:         request.Email,
		Password:      util.ToPointer(request.Password),
		EmailVerified: false,
	}

	err := u.userStorage.SignupUser(user)
	if err != nil && errors.Is(err, &integrations.ErrAlreadyExists{}) {
		return api.SignupResponse{}, err
	} else if err != nil {
		return api.SignupResponse{}, errors.Wrap(err, "can't register user")
	}

	return api.SignupResponse{Message: "OK"}, nil
}

func (u *cloakUserService) SetPassword(user *entities.User, password string) error {
	if user.ID == nil {
		return errors.New("nil user id to delete")
	}
	err := u.userStorage.SetPassword(user, password)
	if err != nil {
		return err
	}
	return nil
}

func (u *cloakUserService) DeleteUser(user *entities.User) error {
	if user.ID == nil {
		return errors.New("nil user id to delete")
	}
	err := u.userStorage.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *cloakUserService) DeleteUserById(userId string) error {
	user, err := u.userStorage.GetUserById(userId)
	if err != nil {
		return errors.Wrap(err, "delete by id")
	}
	err = u.userStorage.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *cloakUserService) GetUserById(userId string) (*entities.User, error) {
	user, err := u.userStorage.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *cloakUserService) SetRoleById(userId string, role entities.UserRole) error {
	user, err := u.userStorage.GetUserById(userId)
	if err != nil {
		return errors.Wrap(err, "set role")
	}
	err = u.userStorage.SetRole(user, role)
	if err != nil {
		return errors.Wrap(err, "set role")
	}
	return nil
}
