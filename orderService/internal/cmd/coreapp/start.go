package coreapp

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	capp "orderServiceGit/internal/app/coreapp"
	"orderServiceGit/internal/core"
	"orderServiceGit/internal/log"

	"github.com/spf13/cobra"
)

func StartCommand(cmd *cobra.Command, args []string) error {
	log.InitLogger()

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewOptionalBackend("config.yaml"),
		file.NewOptionalBackend("config/config.yaml"),
	)

	cfg := core.MultipleConfig{}
	err := loader.Load(context.Background(), &cfg)

	if err != nil {
		log.Errorf("Error loading config: ", err.Error())
		return err
	}

	log.Debug(cfg)

	app, err := capp.NewCoreApp(&cfg.CoreConfig)
	if err != nil {
		return err
	}

	err = app.Start(context.Background())
	if err != nil {
		return err
	}
	return nil
}
