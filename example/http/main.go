package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mantzas/patron/http"
)

func main() {

	s, err := http.New(http.Ports(80, 81))
	if err != nil {
		fmt.Printf("failed to create service %v", err)
		os.Exit(1)
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			fmt.Printf("failed to create service %v", err)
			os.Exit(1)
		}
	}()

	err = s.WaitSignalAndShutdown(5 * time.Second)
	if err != nil {
		fmt.Printf("failed to shutdown service %v", err)
		os.Exit(1)
	}
}
