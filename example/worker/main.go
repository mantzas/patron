package main

import (
	"fmt"
	"os"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zero"
	"github.com/mantzas/patron/worker"
	"github.com/rs/zerolog"
)

func main() {

	// Set up logging
	err := log.Setup(zero.DefaultFactory(zerolog.InfoLevel))
	if err != nil {
		fmt.Printf("failed to setup logging %v", err)
		os.Exit(1)
	}

	// Set up worker
	w, err := worker.New("test")
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
