package gateapp

import (
	gapp "catalogServiceGit/internal/app/gateapp"
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/log"
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
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

	app, err := gapp.NewGateApp(&cfg.GateConfig)
	if err != nil {
		return err
	}

	err = app.Start(context.Background())
	if err != nil {
		return err
	}
	return nil
}
