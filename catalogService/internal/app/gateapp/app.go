package gateapp

import (
	"fmt"
	"log"
	"net/http"

	"catalogServiceGit/internal/app/baseapp"
	"catalogServiceGit/internal/app/gateapp/handlers"
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/core/services"
	"catalogServiceGit/internal/core/services/impl/apisvc"
	"catalogServiceGit/internal/integrations"
	"catalogServiceGit/internal/integrations/rabbitmq"

	"catalogServiceGit/docs/gen"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type gateApp struct {
	*baseapp.BaseWebApp
	apiService services.GateService
	eventBus   integrations.EventBus
	config     *core.Config
}

func NewGateApp(config *core.Config) (baseapp.IApp, error) {
	var err error

	webApp, err := baseapp.NewBaseWebApp(&config.Server, &config.Health, &config.Metrics)
	if err != nil {
		log.Fatal("Cannot initialize base web app: ", err)
	}

	gate := &gateApp{
		BaseWebApp: webApp,
		config:     config,
	}

	eb, err := rabbitmq.NewRabbitMQEventBus(&config.Integrations.Rabbit)
	if err != nil {
		log.Fatal("Cannot initialize rabbit: ", err)
	}
	catalog, err := catalogsvc.NewPostgresCatalogClient(&config.Integrations.Postgres)
	if err != nil {
		log.Fatal("Cannot ping redis: ", err)
	}
	gate.apiService = apisvc.NewApiServiceImpl(catalog, eb)

	gate.initRoutes()
	return gate, nil
}

func (a *gateApp) initRoutes() {
	a.AddUnprotectedRoutes()
}

func (a *gateApp) AddUnprotectedRoutes() {
	if a.config.Swagger.Enabled {
		a.AddSwaggerHandler()
	}
	a.AddRoute(
		http.MethodPost,
		"/api/v1/products",
		a.CreateProductHandler)
	a.AddRoute(
		http.MethodPost,
		"/api/v1/products/{id}",
		a.ChangeProductHandler)
	a.AddRoute(
		http.MethodPost,
		"/api/v1/products/{id}/order",
		a.OrderProductHandler)
	a.AddRoute(
		http.MethodPost,
		"/api/v1/products/{id}/rate",
		a.RateProductHandler)
	a.AddRoute(
		http.MethodPost,
		"/api/v1/products/{id}",
		a.DeleteProductHandler)
	a.AddRoute(
		http.MethodGet,
		"/api/v1/products",
		a.GetSeveralProductsHandler)
}

func (a *gateApp) AddSwaggerHandler() {
	gen.SwaggerInfo.BasePath = "/"
	a.AddRoute(http.MethodGet, fmt.Sprintf("%s/*any", a.config.Swagger.Endpoint), ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (a *gateApp) CreateProductHandler(ctx *gin.Context) {
	handlers.CreateProductHandler(ctx, a.apiService)
}

func (a *gateApp) ChangeProductHandler(ctx *gin.Context) {
	handlers.ChangeProductHandler(ctx, a.apiService)
}

func (a *gateApp) DeleteProductHandler(ctx *gin.Context) {
	handlers.DeleteProductHandler(ctx, a.apiService)
}

func (a *gateApp) OrderProductHandler(ctx *gin.Context) {
	handlers.DeleteProductHandler(ctx, a.apiService)
}

func (a *gateApp) RateProductHandler(ctx *gin.Context) {
	handlers.DeleteProductHandler(ctx, a.apiService)
}

func (a *gateApp) GetSeveralProductsHandler(ctx *gin.Context) {
	handlers.GetSeveralProductsHandler(ctx, a.apiService)
}
