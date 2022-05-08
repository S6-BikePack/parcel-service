package rabbitmq_service

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"parcel-service/config"
	"parcel-service/internal/core/domain"
	"parcel-service/pkg/rabbitmq"
)

type rabbitmqPublisher struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Config     *config.Config
}

func New(rabbitmq *rabbitmq.RabbitMQ, cfg *config.Config) *rabbitmqPublisher {
	return &rabbitmqPublisher{
		Connection: rabbitmq.Connection,
		Channel:    rabbitmq.Channel,
		Config:     cfg,
	}
}

func (rmq *rabbitmqPublisher) CreateParcel(parcel domain.Parcel) error {
	return rmq.publishJson("create", parcel)
}

func (rmq *rabbitmqPublisher) UpdateParcelStatus(id string, status int) error {
	body := struct {
		id     string
		status int
	}{
		id:     id,
		status: status,
	}

	return rmq.publishJson("update.status", body)
}

func (rmq *rabbitmqPublisher) publishJson(topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		rmq.Config.RabbitMQ.Exchange,
		fmt.Sprintf("parcel.%s.%s", rmq.Config.ServiceArea.Identifier, topic),
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
