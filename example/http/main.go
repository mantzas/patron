package main

import (
	"fmt"
	"net/http"
	"os"

	patron_http "github.com/mantzas/patron/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
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
