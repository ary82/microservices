package products

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventProducer interface {
	Produce(key string, data []byte)
}

type rabbitMqProducer struct {
	ch           *amqp.Channel
	exchangeName string
}

func NewEventProducer(exchange string, ch *amqp.Channel) EventProducer {
	return &rabbitMqProducer{
		ch:           ch,
		exchangeName: exchange,
	}
}

func (r *rabbitMqProducer) Produce(key string, data []byte) {
	err := r.ch.Publish(
		r.exchangeName, // exchange
		key,            // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		log.Println("CANNOT PUBLISH:", err)
	}
}
