package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("product")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	fmt.Println("Ctrl+C to stop")
	<-done
	fmt.Println("Stopping...")
}
