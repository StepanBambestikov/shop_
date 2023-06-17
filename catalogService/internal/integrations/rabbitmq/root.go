package rabbitmq

import (
	"catalogServiceGit/internal/core"
	"catalogServiceGit/internal/integrations"
	"catalogServiceGit/internal/log"
	"encoding/json"
	"fmt"
	rmq "github.com/wagslane/go-rabbitmq"
)

type rabbitmqEventBus struct {
	conn      *rmq.Conn
	publisher *rmq.Publisher
	consumers map[string]*rmq.Consumer
	config    *core.RabbitConfig
}

func NewRabbitMQEventBus(config *core.RabbitConfig) (integrations.EventBus, error) {
	conn, err := rmq.NewConn(
		config.DSN,
	)
	if err != nil {
		return nil, err
	}
	publisher, err := rmq.NewPublisher(
		conn,
		rmq.WithPublisherOptionsExchangeName(config.Exchange),
		rmq.WithPublisherOptionsExchangeDurable,
		rmq.WithPublisherOptionsExchangeDeclare,
		rmq.WithPublisherOptionsExchangeKind("x-delayed-message"),
		rmq.WithPublisherOptionsExchangeArgs(map[string]interface{}{
			"x-delayed-type": "topic",
		}),
	)
	if err != nil {
		return nil, err
	}

	return &rabbitmqEventBus{
		conn:      conn,
		publisher: publisher,
		consumers: make(map[string]*rmq.Consumer),
		config:    config,
	}, nil
}

func (reb *rabbitmqEventBus) Send(topic string, message any) error {
	return reb.SendDelayed(topic, message, 0)
}

func (reb *rabbitmqEventBus) SendDelayed(topic string, message any, delay int32) error {
	encodedMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return reb.publisher.Publish(
		encodedMessage,
		[]string{topic},
		rmq.WithPublishOptionsExchange(reb.config.Exchange),
		rmq.WithPublishOptionsContentType("application/json"),
		rmq.WithPublishOptionsHeaders(map[string]interface{}{
			"x-delay": delay,
		}),
	)
}

func (reb *rabbitmqEventBus) Close() error {
	reb.publisher.Close()
	for consumer := range reb.consumers {
		reb.consumers[consumer].Close()
	}
	err := reb.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func getNackAction(discardOnError bool) rmq.Action {
	switch discardOnError {
	case true:
		return rmq.NackDiscard
	default:
		return rmq.NackRequeue
	}
}

func processDelivery(handler func([]byte) error, maxRetries uint64, discardOnError bool) func(rmq.Delivery) (act rmq.Action) {
	return func(d rmq.Delivery) (act rmq.Action) {
		dc, ok := d.Headers["x-delivery-count"].(uint64)
		if !ok {
			dc = 0
		}
		if dc >= maxRetries {
			log.Warn("Too many retries!")
			return rmq.NackDiscard
		}
		defer func() {
			if err := recover(); err != nil {
				log.Warn("Caught panic in rabbitmq consumer: ", err)
				act = getNackAction(discardOnError)
			}
		}()

		err := handler(d.Body)

		if err != nil {
			if discardOnError {
				log.Warn("[rabbitMQ] discarding message ", d, ": ", err.Error())
				return rmq.NackDiscard
			} else {
				log.Warn("[rabbitMQ] requeue message ", d, ": ", err.Error())
				return rmq.NackRequeue
			}
		}
		return rmq.Ack
	}
}

func (reb *rabbitmqEventBus) AddConsumer(name string, topic string, handler func([]byte) error, discardOnError bool) error {
	if reb.consumers[name] != nil {
		return fmt.Errorf("Consumer %s already exists!", name)
	}

	consumer, err := rmq.NewConsumer(
		reb.conn,
		processDelivery(handler, reb.config.MaxRetries, discardOnError),
		fmt.Sprintf("%s.%s", name, topic),
		rmq.WithConsumerOptionsRoutingKey(topic),
		rmq.WithConsumerOptionsConcurrency(1),
		rmq.WithConsumerOptionsExchangeName(reb.config.Exchange),
		rmq.WithConsumerOptionsExchangeKind("x-delayed-message"),
		rmq.WithConsumerOptionsQOSPrefetch(1),
		rmq.WithConsumerOptionsQueueQuorum,
		rmq.WithConsumerOptionsExchangeDurable,
		rmq.WithConsumerOptionsQueueDurable,
		rmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	reb.consumers[name] = consumer
	return nil
}
