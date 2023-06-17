package baseapp

import (
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/core/entities"
	"catalogServiceGit/internal/log"
	"catalogServiceGit/internal/metrics"
	"catalogServiceGit/internal/middleware"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

type BaseWebApp struct {
	Router        *gin.Engine
	serverConfig  *core.ServerConfig
	healthConfig  *core.HealthConfig
	metricsConfig *core.MetricsConfig
}

func NewBaseWebApp(serverCfg *core.ServerConfig, healthCfg *core.HealthConfig, metricsCfg *core.MetricsConfig) (*BaseWebApp, error) {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.GinError())

	if metricsCfg != nil && metricsCfg.Enabled {
		metrep, err := metrics.NewPrometheusRepository(metricsCfg)
		if err != nil {
			return nil, err
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
		Router:        router,
		serverConfig:  serverCfg,
		healthConfig:  healthCfg,
		metricsConfig: metricsCfg,
	}, nil
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
		done <- app.Router.Run(fmt.Sprintf("%s:%d", app.serverConfig.Host, app.serverConfig.Port))
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
	app.Router.NoRoute(handle404)
}

func (app *BaseWebApp) AddRoute(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	app.Router.Handle(httpMethod, relativePath, handlers...)
}

func (app *BaseWebApp) AddMiddleware(_middleware gin.HandlerFunc) {
	app.Router.Use(_middleware)
}
