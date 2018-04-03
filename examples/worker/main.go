package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/mantzas/patron/worker"
	"github.com/mantzas/patron/worker/amqp"
	zl "github.com/rs/zerolog"
)

type helloProcessor struct {
}

func (hp helloProcessor) Process(ctx context.Context, msg []byte) error {
	fmt.Printf("message: %s", string(msg))
	return nil
}

func main() {

	// Set up logging
	err := log.Setup(zerolog.DefaultFactory(zl.InfoLevel))
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
