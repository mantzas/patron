package main

import (
	"fmt"
	"net/http"
	"os"

	patron_http "github.com/mantzas/patron/http"
	"github.com/mantzas/patron/http/httprouter"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from patron!"))
}

func main() {
	// Set up routes
	routes := make([]patron_http.Route, 0)
	routes = append(routes, patron_http.NewRoute("/", http.MethodGet, index))

	s, err := patron_http.New("test", routes, zerolog.Log(log.InfoLevel),
		patron_http.Ports(50000), httprouter.Handler())
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
