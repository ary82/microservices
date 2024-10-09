package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ary82/microservices/internal/products"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("can't initialize queue:", err)
	}

	mqChannel, err := conn.Channel()
	if err != nil {
		log.Fatal("can't initialize channel:", err)
	}
	queue, err := mqChannel.QueueDeclare(
		"q1",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("can't initialize queue:", err)
	}

	// body := "Hello World!"
	// err = mqChannel.PublishWithContext(context.Background(),
	// 	"",         // exchange
	// 	queue.Name, // routing key
	// 	false,      // mandatory
	// 	false,      // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        []byte(body),
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal("can't send message:", err)
	// }
	msgs, err := mqChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatal("can't consume:", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	port := "8001"
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := products.NewGrpcServer(port)

	log.Println("starting grpc server on:", port)
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for an interrupt
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	fmt.Println("Ctrl+C to stop")
	<-done
	fmt.Println("Stopping...")

	// RabbitMQ
	mqChannel.Close()
	conn.Close()

	// grpc
	s.Stop()
}
