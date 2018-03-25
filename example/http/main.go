package main

import (
	"fmt"
	"os"

	"github.com/mantzas/patron/http"
)

func main() {

	s, err := http.New(http.Ports(50000, 50001))
	if err != nil {
		fmt.Printf("failed to create service %v", err)
		os.Exit(1)
	}

	err = s.ListenAndServe()
	if err != nil {
		fmt.Printf("failed to create service %v", err)
		os.Exit(1)
	}
}
