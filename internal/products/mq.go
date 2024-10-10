package products

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitializeConsumerQueue(messageCh *amqp.Channel) error {
	err := messageCh.ExchangeDeclare(
		"inter_microservice", // name
		"direct",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return err
	}

	queue, err := messageCh.QueueDeclare(
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

	keys := []string{
		"USER_REGISTERED",
		"ORDER_PLACED",
	}

	for _, key := range keys {
		err = messageCh.QueueBind(
			queue.Name,           // queue name
			key,                  // routing key
			"inter_microservice", // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := messageCh.Consume(
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
			log.Printf(" [x] %s", d.Body)
		}
	}()

	return nil
}
