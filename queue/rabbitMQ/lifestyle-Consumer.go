package rabbitMQ

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type LifestyleConsumer struct {
	Log    *slog.Logger
	Rabbit *amqp.Connection
}

func (r *LifestyleConsumer) CreateConsume() (<-chan amqp.Delivery, error) {
	channel, err := r.Rabbit.Channel()
	if err != nil {
		r.Log.Error("Failed to open a channel", "error", err)
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"CreateLifestyle",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.Log.Error("Failed to declare a queue", "error", err)
		return nil, err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.Log.Error("Failed to register a consumer", "error", err)
		return nil, err
	}

	return msgs, nil
}

func (r *LifestyleConsumer) DeleteConsume() (<-chan amqp.Delivery, error) {
	channel, err := r.Rabbit.Channel()
	if err != nil {
		r.Log.Error("Failed to open a channel", "error", err)
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"DeleteLifestyle",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.Log.Error("Failed to declare a queue", "error", err)
		return nil, err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.Log.Error("Failed to register a consumer", "error", err)
		return nil, err
	}

	return msgs, nil
}
