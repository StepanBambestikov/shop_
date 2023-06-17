package goauth

import (
	"context"
	"gitea.teneshag.ru/gigabit/goauth/internal/app/goauth"
	"gitea.teneshag.ru/gigabit/goauth/internal/core"
	"gitea.teneshag.ru/gigabit/goauth/internal/log"

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

	cfg := core.Config{}
	err := loader.Load(context.Background(), &cfg)

	if err != nil {
		log.Errorf("Error loading config: ", err.Error())
		return err
	}

	log.Debug(cfg)

	app, err := goauth.NewAuthApp(&cfg)
	if err != nil {
		return err
	}

	err = app.Start(context.Background())
	if err != nil {
		return err
	}
	return nil
}
