package goauth

import (
	"fmt"
	"gitea.teneshag.ru/gigabit/goauth/internal/app"
	"gitea.teneshag.ru/gigabit/goauth/internal/app/goauth/handlers"
	"gitea.teneshag.ru/gigabit/goauth/internal/core"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/services"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/services/impl/usersvc"
	"gitea.teneshag.ru/gigabit/goauth/internal/integrations"
	"gitea.teneshag.ru/gigabit/goauth/internal/integrations/cloak"
	"gitea.teneshag.ru/gigabit/goauth/internal/middleware"
	"net/http"

	"gitea.teneshag.ru/gigabit/goauth/docs/gen"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type authApp struct {
	*app.BaseWebApp
	userStorage integrations.UserStorage
	userService services.UserService
	config      *core.Config
}

func NewAuthApp(config *core.Config) (app.IApp, error) {
	var err error

	auth := &authApp{
		BaseWebApp: app.NewApp(&config.Server, &config.Health, &config.Metrics),
		config:     config,
	}
	auth.userStorage, err = cloak.NewKeycloakUserStorage(config.Integrations.Keycloak)
	if err != nil {
		return nil, err
	}
	auth.userService = usersvc.NewCloakUserService(auth.userStorage)

	auth.initRoutes()
	return auth, nil
}

func (a *authApp) initRoutes() {
	a.AddUnprotectedRoutes()
	// gw.AddMiddleware(func(c *gin.Context) { authMiddleware(c, gw.userStorage) })
	a.AddAuthMiddleware()
	a.AddProtectedRoutes()
}

func (a *authApp) AddUnprotectedRoutes() {
	if a.config.Swagger.Enabled {
		a.AddSwaggerHandler()
	}
	a.AddRoute(
		http.MethodPost,
		"/api/v1/sign/in",
		a.LoginHandler)

	a.AddRoute(
		http.MethodPost,
		"/api/v1/sign/up",
		a.SignupHandler)

	a.AddRoute(
		http.MethodPost,
		"/api/v1/sign/out",
		a.LogoutHandler)
	a.AddRoute(
		http.MethodPost,
		"internal/api/v1/users/:userId/giveRole",
		a.GiveUserRoleHandler)
	a.AddRoute(
		http.MethodDelete,
		"internal/api/v1/users/:userId",
		a.DeleleUserByIdHandler)
}

func (a *authApp) AddAuthMiddleware() {
	a.AddMiddleware(func(c *gin.Context) { middleware.Authorization(c, a.userStorage) })
}

func (a *authApp) AddProtectedRoutes() {
	a.AddRoute(
		http.MethodDelete,
		"/api/v1/users/me",
		a.DeleleUserHandler)
	a.AddRoute(
		http.MethodGet,
		"/api/v1/users/me",
		a.MeHandler)
}

func (a *authApp) AddSwaggerHandler() {
	gen.SwaggerInfo.BasePath = "/"
	a.AddRoute(http.MethodGet, fmt.Sprintf("%s/*any", a.config.Swagger.Endpoint), ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (a *authApp) LoginHandler(ctx *gin.Context) {
	handlers.LoginHandler(ctx, a.userService)
}

func (a *authApp) LogoutHandler(ctx *gin.Context) {
	handlers.LogoutHandler(ctx, a.userService)
}

func (a *authApp) SignupHandler(ctx *gin.Context) {
	handlers.SignupHandler(ctx, a.userService)
}

func (a *authApp) MeHandler(ctx *gin.Context) {
	handlers.MeHandler(ctx)
}

func (a *authApp) DeleleUserHandler(ctx *gin.Context) {
	handlers.DeleteUserHandler(ctx, a.userService)
}

func (a *authApp) DeleleUserByIdHandler(ctx *gin.Context) {
	handlers.DeleteUserByIdHandler(ctx, a.userService)
}

func (a *authApp) GiveUserRoleHandler(ctx *gin.Context) {
	handlers.GiveUserRoleHandler(ctx, a.userService)
}
