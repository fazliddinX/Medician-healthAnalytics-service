package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"medical-service/queue/rabbitMQ"
)

type Queues struct {
	MedicalRecords Queue
}

type Queue interface {
	CreateConsume() (<-chan amqp.Delivery, error)
	DeleteConsume() (<-chan amqp.Delivery, error)
}

func NewMedicalRecordConsume(log *slog.Logger, rabbitConn *amqp.Connection) Queue {
	return &rabbitMQ.RabbitMedicalRecordsConsumer{Log: log, Rabbit: rabbitConn}
}
func NewLifestyleConsumer(log *slog.Logger, rabbitConn *amqp.Connection) Queue {
	return &rabbitMQ.LifestyleConsumer{Log: log, Rabbit: rabbitConn}
}
func NewWearableData(log *slog.Logger, rabbitConn *amqp.Connection) Queue {
	return &rabbitMQ.WearableDataConsumer{Log: log, Rabbit: rabbitConn}
}
