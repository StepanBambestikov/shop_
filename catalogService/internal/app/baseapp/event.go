package baseapp

import (
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/integrations"
	"catalogServiceGit/internal/integrations/rabbitmq"
	"catalogServiceGit/internal/log"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type BaseEventApp struct {
	EventBus integrations.EventBus
	exchange string
	done     chan error
}

func (a *BaseEventApp) Start(ctx context.Context) error {
	var err error
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		a.done <- nil
	}()

	select {
	case err = <-a.done:
		log.Info("Exiting due to signal!")
		break
	case <-ctx.Done():
		log.Info("Exiting due to context closed!")
		break
	}

	log.Info("Closing event app!")
	err1 := a.EventBus.Close()
	if err1 != nil {
		return err1
	}
	return err
}

func (a *BaseEventApp) AddHandler(name string, topic string, handler func([]byte) error, discardOnError bool) error {
	err := a.EventBus.AddConsumer(fmt.Sprintf("%s.%s", a.exchange, name), topic, handler, discardOnError)
	if err != nil {
		return err
	}
	return nil
}

func NewBaseEventApp(cfg *core.RabbitConfig) (*BaseEventApp, error) {
	mq, err := rabbitmq.NewRabbitMQEventBus(cfg)
	if err != nil {
		return nil, err
	}

	return &BaseEventApp{
		EventBus: mq,
		exchange: cfg.Exchange,
		done:     make(chan error, 1),
	}, nil
}
