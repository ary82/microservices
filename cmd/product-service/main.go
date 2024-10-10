package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ary82/microservices/internal/products"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	mode := os.Getenv("MODE")
	if mode != "dev" && mode != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env:", err)
		}
	}

	port := os.Getenv("PRODUCTS_SERVICE_GRPC_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := products.NewSqlDb(os.Getenv("PRODUCTS_SERVICE_DB_URL"))
	if err != nil {
		log.Fatalf("failed to mount db: %v", err)
	}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_CONN_STR"))
	if err != nil {
		log.Fatal("can't initialize queue:", err)
	}

	mqChannel, err := conn.Channel()
	if err != nil {
		log.Fatal("can't initialize channel:", err)
	}

	err = products.InitializeConsumerQueue(mqChannel)
	if err != nil {
		log.Fatal("mq initialization failed:", err)
	}

	pr := products.NewProductsRepository(db)
	ps := products.NewProductsService(pr, mqChannel)
	s := products.NewGrpcServer(port, ps)

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
	// mqChannel.Close()
	// conn.Close()

	// grpc
	s.Stop()
}
