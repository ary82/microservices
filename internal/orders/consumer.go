package orders

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerService interface {
	Consume(events []string) error
}

type rabbitMqConsumer struct {
	ch           *amqp.Channel
	exchangeName string
	ordersRepo   OrdersRepository
}

func NewEventConsumer(ch *amqp.Channel, exchangeName string, ordersRepo OrdersRepository) ConsumerService {
	return &rabbitMqConsumer{
		ch:           ch,
		exchangeName: exchangeName,
		ordersRepo:   ordersRepo,
	}
}

func (r *rabbitMqConsumer) Consume(events []string) error {
	err := r.ch.ExchangeDeclare(
		r.exchangeName, // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	queue, err := r.ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, key := range events {
		err = r.ch.QueueBind(
			queue.Name,     // queue name
			key,            // routing key
			r.exchangeName, // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := r.ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			switch d.RoutingKey {
			case "USER_REGISTERED":
				log.Printf("New user: %s", string(d.Body))
			case "PRODUCT_CREATED":
				log.Printf("New product: %s", string(d.Body))
			case "INVENTORY_UPDATE":
				log.Printf("Inventory Update: %s", string(d.Body))
			default:
				log.Printf(" [rabbitMQ unmapped: %s] %s", d.RoutingKey, d.Body)
			}
		}
	}()

	return nil
}
