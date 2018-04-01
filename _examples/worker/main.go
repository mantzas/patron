package main

import (
	"fmt"
	"os"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zero"
	"github.com/mantzas/patron/worker"
	"github.com/mantzas/patron/worker/amqp"
	"github.com/rs/zerolog"
)

type helloProcessor struct {
}

func (hp helloProcessor) Process(msg []byte) error {
	fmt.Printf("message: %s", string(msg))
	return nil
}

func main() {

	// Set up logging
	err := log.Setup(zero.DefaultFactory(zerolog.InfoLevel))
	if err != nil {
		fmt.Printf("failed to setup logging %v", err)
		os.Exit(1)
	}

	// setting up a amqp processor
	p, err := amqp.New("http://localhost:8081", "test", &helloProcessor{})

	// Set up worker
	w, err := worker.New("test", p)
	if err != nil {
		fmt.Printf("failed to create worker %v", err)
		os.Exit(1)
	}

	err = w.Run()
	if err != nil {
		fmt.Printf("failed to run worker %v", err)
		os.Exit(1)
	}
}
