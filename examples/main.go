package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/async/amqp"
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
	err = config.Set("rabbitmq_url", "amqp://localhost:8081")
	if err != nil {
		fmt.Printf("failed to set rabbitmq URL config %v", err)
		os.Exit(1)
	}

	exporter := &logExporter{}
	view.RegisterExporter(exporter)
	trace.RegisterExporter(exporter)

	// Always trace for this demo.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Report stats at every second.
	view.SetReportingPeriod(10 * time.Second)
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

	rabbitmqURL, err := config.GetString("rabbitmq_url")
	if err != nil {
		fmt.Printf("failed to get rabbitmq URL config %v", err)
		os.Exit(1)
	}

	// setting up a amqp processor
	_, err = amqp.New(rabbitmqURL, "test", &helloProcessor{})
	if err != nil {
		fmt.Print("failed to create AMQP service", err)
		os.Exit(1)
	}

	// Set up routes
	routes := make([]sync_http.Route, 0)
	routes = append(routes, sync_http.NewRoute("/", http.MethodGet, index))

	options := []sync_http.Option{
		sync_http.SetPorts(50000),
		sync_http.SetRoutes(routes),
	}

	httpSrv, err := sync_http.New(httprouter.CreateHandler, options...)
	if err != nil {
		fmt.Print("failed to create HTTP service", err)
		os.Exit(1)
	}

	patron.Metric(&logExporter{}, 5*time.Second)

	srv, err := patron.New("test", httpSrv /*, amqpSrv*/)
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
