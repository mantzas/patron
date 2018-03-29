package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	patron_http "github.com/mantzas/patron/http"
	"github.com/mantzas/patron/http/httprouter"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zero"
)

func index(w http.ResponseWriter, r *http.Request) {
}

func main() {

	// Set up logging
	err := log.Setup(zero.DefaultFactory(zerolog.InfoLevel))
	if err != nil {
		fmt.Printf("failed to setup logging %v", err)
		os.Exit(1)
	}

	// Set up routes
	routes := make([]patron_http.Route, 0)
	routes = append(routes, patron_http.NewRoute("/", http.MethodGet, index))

	s, err := patron_http.New("test", routes, patron_http.Ports(50000, 50001), httprouter.Handler())
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
