package coreapp

import (
	"context"
	"log"
	"orderServiceGit/internal/app/baseapp"
	"orderServiceGit/internal/core"
	"orderServiceGit/internal/core/services"
	"orderServiceGit/internal/core/services/impl/coresvc"
	catalogService "orderServiceGit/internal/core/services/impl/ordersvc"
	"orderServiceGit/internal/util"

	"golang.org/x/sync/errgroup"
)

type coreApp struct {
	webApp      baseapp.BaseWebApp
	eventApp    baseapp.BaseEventApp
	coreService services.CoreService
	config      *core.Config
}

func (f *coreApp) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return f.webApp.Start(ctx)
	})

	g.Go(func() error {
		return f.eventApp.Start(ctx)
	})

	return g.Wait()
}

func (c *coreApp) initConsumers() error {
	err := c.eventApp.AddHandler("CreateOrder", core.TopicCreateOrder, util.WrapEvent(c.coreService.CreateOrder), false)
	if err != nil {
		return err
	}
	err = c.eventApp.AddHandler("ChangeProduct", core.TopicDeleteOrder, util.WrapEvent(c.coreService.DeleteOrder), false)
	if err != nil {
		return err
	}
	err = c.eventApp.AddHandler("DeleteProduct", core.TopicSetOrderStatus, util.WrapEvent(c.coreService.SetOrderStatus), false)
	if err != nil {
		return err
	}
	return nil
}

func NewCoreApp(config *core.Config) (baseapp.IApp, error) {
	eventApp, err := baseapp.NewBaseEventApp(&config.Integrations.Rabbit)
	if err != nil {
		log.Fatal("Cannot initialize event app: ", err)
	}
	webApp, err := baseapp.NewBaseWebApp(&config.Server, &config.Health, &config.Metrics)
	if err != nil {
		log.Fatal("Cannot initialize web app: ", err)
	}
	catalogservice, err := catalogService.NewPostgresOrderClient(&config.Integrations.Postgres)
	if err != nil {
		log.Fatal("Cannot initialize counter service: ", err)
	}
	coreSvc := coresvc.NewCoreServiceImpl(catalogservice)

	app := &coreApp{
		webApp:      *webApp,
		eventApp:    *eventApp,
		coreService: coreSvc,
		config:      config,
	}
	err = app.initConsumers()
	if err != nil {
		log.Fatal("Cannot initialize event consumers: ", err)
	}

	return app, nil
}
