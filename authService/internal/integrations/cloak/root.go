package cloak

import (
	"context"
	"fmt"
	"gitea.teneshag.ru/gigabit/goauth/internal/core"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/integrations"
	"gitea.teneshag.ru/gigabit/goauth/internal/log"
	"gitea.teneshag.ru/gigabit/goauth/internal/util"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"

	gocloak "github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
)

type keycloakUserStorage struct {
	config     core.KeycloakConfig
	cloak      *gocloak.GoCloak
	cloakMutex sync.Mutex
	adminJWT   *gocloak.JWT
}

var ctx = context.Background()

func refreshAdminTokenByTicker(us *keycloakUserStorage, ticker *time.Ticker, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			log.Info("Refreshing token in keycloak")
			err := us.refreshAdminToken(ctx)
			if err != nil {
				log.Error("Can't refresh token! ", err.Error())
			}
		}
	}

}

func NewKeycloakUserStorage(config core.KeycloakConfig) (integrations.UserStorage, error) {
	var err error
	store := keycloakUserStorage{
		cloak:    gocloak.NewClient(config.Uri),
		config:   config,
		adminJWT: &gocloak.JWT{},
	}

	store.adminJWT, err = store.cloak.LoginAdmin(
		ctx,
		config.Admin.Username,
		config.Admin.Password,
		config.Admin.Realm,
	)

	if err != nil {
		return nil, errors.Wrap(err, "initializing")
	}

	ticker := time.NewTicker(time.Duration(config.TokenRefreshInterval) * time.Second)
	go refreshAdminTokenByTicker(&store, ticker, ctx)

	return &store, nil
}

func (k *keycloakUserStorage) refreshAdminToken(ctx context.Context) (err error) {
	adminJWTpt, err := k.cloak.RefreshToken(
		ctx,
		k.adminJWT.RefreshToken,
		"admin-cli",
		"",
		k.config.Admin.Realm,
	)
	if err != nil && (strings.Contains(err.Error(), "400 Bad Request: invalid_grant: Session not active") ||
		strings.Contains(err.Error(), "400 Bad Request: invalid_grant: Token is not active")) {
		adminJWTpt, err = k.cloak.LoginAdmin(
			ctx,
			k.config.Admin.Username,
			k.config.Admin.Password,
			k.config.Admin.Realm,
		)
	}
	if err != nil {
		return errors.Wrap(err, "refreshing admin jwt")
	}

	k.cloakMutex.Lock()
	k.adminJWT = adminJWTpt
	k.cloakMutex.Unlock()

	return nil
}
func (k *keycloakUserStorage) getFromTokenOffline(accessToken string) (*entities.User, error) {
	claims := jwt.MapClaims{}
	result, err := k.cloak.DecodeAccessTokenCustomClaims(context.Background(), accessToken, k.config.Realm, claims)
	if err != nil {
		return nil, err
	}

	if !result.Valid {
		return nil, &integrations.ErrInvalidCredentials{Message: "token is invalid"}
	}

	log.Debug(claims)
	emailVerified, err := strconv.ParseBool(fmt.Sprint(claims["email_verified"]))
	if err != nil {
		return nil, errors.Wrap(err, "can't get email verfication status")
	}

	rawMarketRole := fmt.Sprint(claims["market_role"])
	marketRole, ok := entities.UserRoleFromString(rawMarketRole)
	if !ok {
		return nil, errors.Wrap(err, fmt.Sprintf("unknown market role: %s", rawMarketRole))
	}

	user := entities.User{
		ID:            util.ToPointer(fmt.Sprint(claims["sub"])),
		FirstName:     util.ToPointer(fmt.Sprint(claims["first_name"])),
		LastName:      util.ToPointer(fmt.Sprint(claims["last_name"])),
		Username:      fmt.Sprint(claims["preferred_username"]),
		Email:         fmt.Sprint(claims["email"]),
		MarketRole:    marketRole,
		EmailVerified: emailVerified,
	}

	return &user, nil
}

func (k *keycloakUserStorage) GetFromToken(accessToken string) (*entities.User, error) {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()
	introspectResult, err := k.cloak.RetrospectToken(ctx,
		accessToken,
		k.config.Client.ID,
		k.config.Client.Secret,
		k.config.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "get from token")
	}

	if !(*introspectResult.Active) {
		return nil, errors.Wrap(integrations.ErrUnauthorized{}, "get from token")
	}

	user, err := k.getFromTokenOffline(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "get from token")
	}
	return user, nil
}

func (k *keycloakUserStorage) SignupUser(user *entities.User) error {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()

	userId, err := k.cloak.CreateUser(context.Background(),
		k.adminJWT.AccessToken,
		k.config.Realm,
		gocloak.User{
			Username:      util.ToPointer(user.Username),
			Enabled:       util.ToPointer(true),
			EmailVerified: util.ToPointer(false),
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         util.ToPointer(user.Email),
			Attributes: &map[string][]string{
				"market_role": {"notverified"},
			},
			RequiredActions: nil,
		},
	)

	// todo: add already exists handler
	if err != nil && strings.Contains(err.Error(), "409 Conflict: User exists with same username") {
		err = &integrations.ErrAlreadyExists{}
	}
	if err != nil {
		return errors.Wrap(err, "creating user")
	}
	err = k.cloak.SetPassword(ctx,
		k.adminJWT.AccessToken,
		userId,
		k.config.Realm,
		util.FromPointer(user.Password),
		false,
	)
	if err != nil {
		return errors.Wrap(err, "creating user")
	}

	log.Info("Successfully registered user with id ", userId)

	return nil
}

func (k *keycloakUserStorage) LoginUser(login string, password string) (*entities.User, string, string, error) {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()

	_jwt, err := k.cloak.Login(ctx,
		k.config.Client.ID,
		k.config.Client.Secret,
		k.config.Realm,
		login,
		password,
	)

	if err != nil && strings.HasPrefix(err.Error(), "401 Unauthorized") {
		err = &integrations.ErrUnauthorized{}
	}
	if err != nil {
		return nil, "", "", errors.Wrap(err, "login")
	}

	user, err := k.getFromTokenOffline(_jwt.AccessToken)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "login")
	}

	return user, _jwt.AccessToken, _jwt.RefreshToken, nil
}

func (k *keycloakUserStorage) RefreshToken(refreshToken string) (string, string, error) {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()

	token, err := k.cloak.RefreshToken(ctx,
		refreshToken,
		k.config.Client.ID,
		k.config.Client.Secret,
		k.config.Realm,
	)

	if err != nil && strings.Contains(err.Error(), "Session not active") {
		err = &integrations.ErrUnauthorized{}
	}
	if err != nil {
		return "", "", errors.Wrap(err, "refresh token")
	}

	return token.AccessToken, token.RefreshToken, nil
}

func (k *keycloakUserStorage) LogoutByID(userId string) error {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()

	err := k.cloak.LogoutAllSessions(ctx,
		k.adminJWT.AccessToken,
		k.config.Realm,
		userId,
	)
	if err != nil {
		return errors.Wrap(err, "logout by id")
	}

	return nil
}

func (k *keycloakUserStorage) LogoutUser(refreshToken string) error {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()

	err := k.cloak.Logout(ctx,
		k.config.Client.ID,
		k.config.Client.Secret,
		k.config.Realm,
		refreshToken,
	)
	if err != nil {
		return errors.Wrap(err, "logout user")
	}

	return nil
}

func (k *keycloakUserStorage) DeleteUser(user *entities.User) error {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()

	err := k.cloak.DeleteUser(ctx,
		k.adminJWT.AccessToken,
		k.config.Realm,
		util.FromPointer(user.ID),
	)

	if err != nil {
		return errors.Wrap(err, "delete user")
	}
	return nil
}

func (k *keycloakUserStorage) SetPassword(user *entities.User, password string) error {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()
	err := k.cloak.SetPassword(ctx,
		k.adminJWT.AccessToken,
		util.FromPointer(user.ID),
		k.config.Realm,
		password,
		false,
	)
	if err != nil {
		return errors.Wrap(err, "set password")
	}
	err = k.LogoutByID(util.FromPointer(user.ID))
	if err != nil {
		return errors.Wrap(err, "logout")
	}
	return nil
}

func (k *keycloakUserStorage) GetUserById(userId string) (*entities.User, error) {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()

	cloakUser, err := k.cloak.GetUserByID(
		context.Background(),
		k.adminJWT.AccessToken,
		k.config.Realm,
		userId,
	)
	if err != nil {
		return nil, errors.Wrap(err, "get by id")
	}
	if cloakUser.ID == nil {
		return nil, errors.Wrap(&integrations.ErrNotFound{}, "get by id")
	}
	rawMarketRole := (*cloakUser.Attributes)["market_role"][0]
	marketRole, ok := entities.UserRoleFromString(rawMarketRole)
	if !ok {
		return nil, errors.Wrap(err, fmt.Sprintf("unknown market role: %s", rawMarketRole))
	}
	user := &entities.User{
		ID:            cloakUser.ID,
		FirstName:     cloakUser.FirstName,
		LastName:      cloakUser.LastName,
		Username:      *cloakUser.Username,
		Email:         *cloakUser.Email,
		EmailVerified: *cloakUser.EmailVerified,
		MarketRole:    marketRole,
	}
	return user, nil
}

func (k *keycloakUserStorage) SetRole(user *entities.User, role entities.UserRole) error {
	k.cloakMutex.Lock()
	defer k.cloakMutex.Unlock()
	err := k.cloak.UpdateUser(ctx,
		k.adminJWT.AccessToken,
		k.config.Realm,
		gocloak.User{
			ID:            user.ID,
			EmailVerified: &user.EmailVerified,
			Attributes: &map[string][]string{
				"market_role": {role.String()},
			},
			RequiredActions: nil,
		},
	)
	if err != nil {
		return errors.Wrap(err, "set role")
	}
	err = k.LogoutByID(util.FromPointer(user.ID))
	if err != nil {
		return errors.Wrap(err, "logout")
	}

	return nil
}
