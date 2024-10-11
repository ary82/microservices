package products

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/ary82/microservices/internal/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerService interface {
	Consume(events []string) error
}

type rabbitMqConsumer struct {
	ch           *amqp.Channel
	exchangeName string
	productsRepo ProductsRepository
}

func NewEventConsumer(ch *amqp.Channel, exchangeName string, productsRepo ProductsRepository) ConsumerService {
	return &rabbitMqConsumer{
		ch:           ch,
		exchangeName: exchangeName,
		productsRepo: productsRepo,
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
			case "ORDER_PLACED":
				r.consumeOrderPlaced(d.Body)
			default:
				log.Printf(" [rabbitMQ unmapped: %s] %s", d.RoutingKey, d.Body)
			}
		}
	}()

	return nil
}

func (r *rabbitMqConsumer) consumeOrderPlaced(body []byte) {
	decodedBody := mq.OrderPlacedPayload{}

	err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&decodedBody)
	if err != nil {
		log.Println(err)
	}

	for _, v := range decodedBody.Products {
		err = r.productsRepo.UpdateInventorySubtract(v.ProductId, int64(v.Quantity))
		if err != nil {
			log.Println(err)
		}
	}
}
