package coreapp

import (
	"catalogServiceGit/internal/app/baseapp"
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/core/services"
	catalogService "catalogServiceGit/internal/core/services/impl/catalogsvc"
	"catalogServiceGit/internal/core/services/impl/coresvc"
	"catalogServiceGit/internal/util"
	"context"
	"fmt"
	"log"

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
	err := c.eventApp.AddHandler("CreateProduct", core.TopicCreateProduct, util.WrapEvent(c.coreService.CreateProduct), false)
	if err != nil {
		return err
	}
	err = c.eventApp.AddHandler("ChangeProduct", core.TopicChangeProduct, util.WrapEvent(c.coreService.ChangeProduct), false)
	if err != nil {
		return err
	}
	err = c.eventApp.AddHandler("DeleteProduct", core.TopicDeleteProduct, util.WrapEvent(c.coreService.DeleteProduct), false)
	if err != nil {
		return err
	}
	err = c.eventApp.AddHandler("RateProduct", core.TopicRateProduct, util.WrapEvent(c.coreService.RateProduct), false)
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
	fmt.Println("=====================================================================")
	fmt.Println(config.Integrations.Postgres)
	fmt.Println("=====================================================================")
	catalogservice, err := catalogService.NewPostgresCatalogClient(&config.Integrations.Postgres)
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
