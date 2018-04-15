package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/http/httprouter"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/mantzas/patron/worker/amqp"
)

type helloProcessor struct {
}

func (hp helloProcessor) Process(ctx context.Context, msg []byte) error {
	fmt.Printf("message: %s", string(msg))
	return nil
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from patron!"))
}

func main() {

	// Set up logging
	err := log.Setup(zerolog.DefaultFactory(log.InfoLevel))
	if err != nil {
		fmt.Printf("failed to setup logging %v", err)
		os.Exit(1)
	}

	// Set up routes
	routes := make([]patron.Route, 0)
	routes = append(routes, patron.NewRoute("/", http.MethodGet, index))

	// setting up a amqp processor
	p, err := amqp.New("http://localhost:8081", "test", &helloProcessor{})
	if err != nil {
		fmt.Print("failed to setup amqp processor", err)
		os.Exit(1)
	}

	options := []patron.Option{
		patron.SetPorts(50000),
		patron.SetRoutes(routes),
		patron.SetProcessor(p),
	}

	s, err := patron.New("test", httprouter.CreateHandler, options...)
	if err != nil {
		fmt.Printf("failed to create service %v", err)
		os.Exit(1)
	}

	err = s.Run()
	if err != nil {
		fmt.Printf("failed to create service %v", err)
		os.Exit(1)
	}
}
