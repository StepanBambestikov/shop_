package app

import (
	"context"
	"fmt"
	"gitea.teneshag.ru/gigabit/goauth/internal/core"
	"gitea.teneshag.ru/gigabit/goauth/internal/core/entities"
	"gitea.teneshag.ru/gigabit/goauth/internal/log"
	"gitea.teneshag.ru/gigabit/goauth/internal/metrics"
	"gitea.teneshag.ru/gigabit/goauth/internal/middleware"

	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type IApp interface {
	Start(ctx context.Context) error
}

type BaseWebApp struct {
	router        *gin.Engine
	serverConfig  *core.ServerConfig
	healthConfig  *core.HealthConfig
	metricsConfig *core.MetricsConfig
}

func NewApp(srv *core.ServerConfig, healthCfg *core.HealthConfig, metricsCfg *core.MetricsConfig) *BaseWebApp {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.GinError())

	if metricsCfg != nil && metricsCfg.Enabled {
		metrep, err := metrics.NewPrometheusRepository(metricsCfg)
		if err != nil {
			log.Fatal("cannot initialize metrics repository")
		}
		router.GET(metricsCfg.Endpoint, metrep.GetHandlerFunc())
		router.Handlers = append([]gin.HandlerFunc{middleware.Metrics(metrep)}, router.Handlers...)
	}

	if healthCfg != nil && healthCfg.Enabled {
		router.GET(healthCfg.Endpoint, func(context *gin.Context) {
			context.JSON(200, entities.ApiReply{
				Message: "OK",
			})
		})
	}

	return &BaseWebApp{
		router:        router,
		serverConfig:  srv,
		healthConfig:  healthCfg,
		metricsConfig: metricsCfg,
	}
}

func (app *BaseWebApp) Start(ctx context.Context) (err error) {
	sigs := make(chan os.Signal, 1)
	done := make(chan error)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Warn("Stopping due to signal ", sig)
		done <- nil
	}()

	go func() {
		done <- app.router.Run(fmt.Sprintf("%s:%d", app.serverConfig.Host, app.serverConfig.Port))
	}()

	select {
	case err = <-done:
		log.Info("Exiting due to error!")
		break
	case <-ctx.Done():
		log.Info("Exiting due to context closed!")
		break
	}
	return err
}

func (app *BaseWebApp) SetHandlers(handle404 gin.HandlerFunc) {
	app.router.NoRoute(handle404)
}

func (app *BaseWebApp) AddRoute(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	app.router.Handle(httpMethod, relativePath, handlers...)
}

func (app *BaseWebApp) AddMiddleware(_middleware gin.HandlerFunc) {
	app.router.Use(_middleware)
}
