package rabbitmq_service

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"parcel-service/internal/core/domain"
	"parcel-service/pkg/rabbitmq"
)

type rabbitmqPublisher rabbitmq.RabbitMQ

func New(rabbitmq *rabbitmq.RabbitMQ) *rabbitmqPublisher {
	return &rabbitmqPublisher{
		Connection: rabbitmq.Connection,
		Channel:    rabbitmq.Channel,
	}
}

func (rmq *rabbitmqPublisher) CreateParcel(parcel domain.Parcel) error {
	return rmq.publishJson("parcel.create", parcel)
}

func (rmq *rabbitmqPublisher) UpdateParcelStatus(id string, status int) error {
	body := struct {
		id     string
		status int
	}{
		id:     id,
		status: status,
	}

	return rmq.publishJson("parcel.update.status", body)
}

func (rmq *rabbitmqPublisher) publishJson(topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		"topics",
		topic,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         js,
		},
	)

	return err
}
