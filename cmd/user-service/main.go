package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ary82/microservices/internal/users"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("user")

	mode := os.Getenv("MODE")
	if mode != "dev" && mode != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env:", err)
		}
	}

	port := "8002"
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := users.NewGrpcServer(port)

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

	// grpc
	s.Stop()
}
