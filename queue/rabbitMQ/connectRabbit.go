package rabbitMQ

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"medical-service/pkg/config"
)

func ConnectToRabbit(cfg config.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.RABBIT_URL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
