package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	patron_http "github.com/mantzas/patron/http"
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

	// Set up HTTP router
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	s, err := patron_http.New("test", mux, patron_http.Ports(50000, 50001))
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
