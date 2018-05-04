package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/config"
	"github.com/mantzas/patron/config/viper"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	sync_http "github.com/mantzas/patron/sync/http"
	"github.com/mantzas/patron/sync/http/httprouter"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type logExporter struct{}

// ExportView logs the view data.
func (e *logExporter) ExportView(vd *view.Data) {
	log.Infof("view export: %v", *vd)
}

// ExportSpan logs the trace span.
func (e *logExporter) ExportSpan(vd *trace.SpanData) {
	log.Infof("span export: %v", *vd)
}

type helloProcessor struct {
}

func (hp helloProcessor) Process(ctx context.Context, msg []byte) error {
	log.Infof("message: %s", string(msg))
	return nil
}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from patron!"))
	if err != nil {
		log.Errorf("failed to write response %v", err)
	}
}

func init() {
	// Set up config (should come from flag, env, file etc)
	err := config.Setup(viper.New())
	if err != nil {
		fmt.Printf("failed to setup config %v", err)
		os.Exit(1)
	}

	err = config.Set("log_level", log.InfoLevel)
	if err != nil {
		fmt.Printf("failed to set log level config %v", err)
		os.Exit(1)
	}
}

func main() {

	// Set up logging
	lvl, err := config.Get("log_level")
	if err != nil {
		fmt.Printf("failed to get log level config %v", err)
		os.Exit(1)
	}

	err = log.Setup(zerolog.DefaultFactory(lvl.(log.Level)))
	if err != nil {
		fmt.Printf("failed to setup logging %v", err)
		os.Exit(1)
	}

	// Set up routes
	routes := make([]sync_http.Route, 0)
	routes = append(routes, sync_http.NewRoute("/", http.MethodGet, index))

	options := []sync_http.Option{
		sync_http.Port(50000),
		sync_http.Routes(routes),
	}

	httpSrv, err := sync_http.New(httprouter.CreateHandler, options...)
	if err != nil {
		fmt.Print("failed to create HTTP service", err)
		os.Exit(1)
	}

	le := logExporter{}
	opts := []patron.Option{
		patron.Metric(&le, 5*time.Second),
		patron.Trace(&le, trace.Config{DefaultSampler: trace.AlwaysSample()}),
	}

	srv, err := patron.New("test", []patron.Service{httpSrv}, opts...)
	if err != nil {
		fmt.Printf("failed to create service %v", err)
		os.Exit(1)
	}

	err = srv.Run()
	if err != nil {
		fmt.Printf("failed to create service %v", err)
		os.Exit(1)
	}
}
