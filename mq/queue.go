package mq

import (
	"github.com/streadway/amqp"
)

const exchangeName = "cita"

type Queue interface {
	Sub() (<-chan amqp.Delivery, error)
	Pub(key RoutingKey, data []byte) error
}

func New(url, name string, keys []RoutingKey) (Queue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		return nil, err
	}

	newQueue, err := ch.QueueDeclare(name, true, false, false, false, nil)
	for _, key := range keys {
		if err := ch.QueueBind(newQueue.Name, string(key), exchangeName, false, nil); err != nil {
			return nil, err
		}
	}

	return &queue{ch: ch, name: name}, nil
}

type queue struct {
	ch   *amqp.Channel
	name string
}

func (q *queue) Sub() (<-chan amqp.Delivery, error) {
	return q.ch.Consume(q.name, "", false, false, true, false, nil)
}

func (q *queue) Pub(key RoutingKey, data []byte) error {
	msg := amqp.Publishing{
		Body: data,
	}
	return q.ch.Publish(exchangeName, string(key), false, false, msg)
}
