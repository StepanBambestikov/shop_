package gateapp

import (
	"fmt"
	"log"
	"net/http"

	"orderServiceGit/internal/app/baseapp"
	"orderServiceGit/internal/app/gateapp/handlers"
	"orderServiceGit/internal/core"
	"orderServiceGit/internal/core/services"
	"orderServiceGit/internal/core/services/impl/apisvc"
	ordersvc "orderServiceGit/internal/core/services/impl/ordersvc"
	"orderServiceGit/internal/integrations"
	"orderServiceGit/internal/integrations/rabbitmq"

	"orderServiceGit/docs/gen"

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
	catalog, err := ordersvc.NewPostgresOrderClient(&config.Integrations.Postgres)
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
		"/internal/orders/setStatus/{id}",
		a.SetOrderStatusHandler)
	a.AddRoute(
		http.MethodDelete,
		"/orders/{id}",
		a.DeleteOrderHandler)
	a.AddRoute(
		http.MethodGet,
		"/orders/{id}",
		a.GetOrderInfoHandler)
	a.AddRoute(
		http.MethodGet,
		"/orders",
		a.GetUserOrdersHandler)
	a.AddRoute(
		http.MethodPost,
		"/orders/{id}",
		a.CreateOrderHandler)
}

func (a *gateApp) AddSwaggerHandler() {
	gen.SwaggerInfo.BasePath = "/"
	a.AddRoute(http.MethodGet, fmt.Sprintf("%s/*any", a.config.Swagger.Endpoint), ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (a *gateApp) DeleteOrderHandler(ctx *gin.Context) {
	handlers.DeleteOrderHandler(ctx, a.apiService)
}

func (a *gateApp) GetOrderInfoHandler(ctx *gin.Context) {
	handlers.GetOrderInfoHandler(ctx, a.apiService)
}

func (a *gateApp) GetUserOrdersHandler(ctx *gin.Context) {
	handlers.GetUserOrdersHandler(ctx, a.apiService)
}

func (a *gateApp) SetOrderStatusHandler(ctx *gin.Context) {
	handlers.SetOrderStatusHandler(ctx, a.apiService)
}

func (a *gateApp) CreateOrderHandler(ctx *gin.Context) {
	handlers.CreateOrderHandler(ctx, a.apiService)
}
