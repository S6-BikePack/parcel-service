package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"parcel-service/config"
	"parcel-service/internal/core/domain"
	"parcel-service/internal/core/interfaces"
	"parcel-service/pkg/rabbitmq"
)

type rabbitmqHandler struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  interfaces.ParcelService
	handlers map[string]func(topic string, body []byte, handler *rabbitmqHandler) error
	config   *config.Config
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, service interfaces.ParcelService, cfg *config.Config) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq: rabbitmq,
		service:  service,
		handlers: map[string]func(topic string, body []byte, handler *rabbitmqHandler) error{
			"customer.create": CustomerCreateOrUpdate,
			"customer.update": CustomerCreateOrUpdate,
			"delivery." + cfg.ServiceArea.Identifier + ".create": DeliveryCreated,
		},
		config: cfg,
	}
}

func DeliveryCreated(topic string, body []byte, handler *rabbitmqHandler) error {
	var delivery struct {
		Parcel domain.Parcel
	}

	if err := json.Unmarshal(body, &delivery); err != nil {
		return err
	}

	_, err := handler.service.UpdateParcelStatus(delivery.Parcel.ID, 1)

	if err != nil {
		return err
	}

	return nil
}

func CustomerCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var customer domain.Customer

	if err := json.Unmarshal(body, &customer); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateCustomer(customer); err != nil {
		return err
	}

	return nil
}

func (handler *rabbitmqHandler) Listen() {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		handler.config.Server.Service+"-"+handler.config.ServiceArea.Identifier,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	for _, s := range maps.Keys(handler.handlers) {
		err = handler.rabbitmq.Channel.QueueBind(
			q.Name,
			s,
			handler.config.RabbitMQ.Exchange,
			false,
			nil)
		if err != nil {
			return
		}
	}

	msgs, err := handler.rabbitmq.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fun, exist := handler.handlers[msg.RoutingKey]

			if exist {
				err = fun(msg.RoutingKey, msg.Body, handler)
				if err == nil {
					msg.Ack(false)
					continue
				}
			}

			fmt.Println(err)
			msg.Nack(false, true)
		}
	}()

	<-forever
}

type MessageHandler struct {
	topic   string
	handler func(topic string, body []byte, handler *rabbitmqHandler) error
}
