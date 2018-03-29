package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	patron_http "github.com/mantzas/patron/http"
	"github.com/mantzas/patron/http/httprouter"
	"github.com/mantzas/patron/http/route"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zero"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {

	// Set up logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	f := zero.NewFactory(&zl)
	log.Setup(f)

	// Set up routes
	routes := make([]route.Route, 0)
	routes = append(routes, route.New("/", http.MethodGet, index))

	// Set up HTTP router
	h, err := httprouter.CreateHandler(routes)
	if err != nil {
		fmt.Printf("failed to create handler %v", err)
		os.Exit(1)
	}

	s, err := patron_http.New("test", h, patron_http.Ports(50000, 50001))
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
