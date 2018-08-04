package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/sync/http"
)

// Audit records name and time of a processing step.
type Audit struct {
	Name     string
	Started  time.Time
	Duration time.Duration
}

// Audits is a collection of audit entries.
type Audits []Audit

func (a *Audits) append(aud Audit) {
	dur := time.Duration(0)
	if len(*a) > 0 {
		dur = aud.Started.Sub((*a)[len(*a)-1].Started)
	}
	aud.Duration = dur
	*a = append(*a, aud)
}

const (
	amqpURL      = "amqp://admin:admin@localhost:5672/"
	amqpExchange = "patron"
	amqpQueue    = "patron"
	kafkaTopic   = "patron-topic"
	kafkaBroker  = "localhost:9092"
)

var logger log.Logger

func init() {
	err := os.Setenv("PATRON_LOG_LEVEL", "debug")
	if err != nil {
		fmt.Printf("failed to set log level env var: %v", err)
		os.Exit(1)
	}
	err = os.Setenv("PATRON_JAEGER_SAMPLER_PARAM", "1.0")
	if err != nil {
		fmt.Printf("failed to set sampler env vars:: %v", err)
		os.Exit(1)
	}
}

func main() {
	cfg, err := patron.Configure("patron", "1.0.0")
	if err != nil {
		fmt.Printf("failed to configure patron: %v", err)
		os.Exit(1)
	}

	amqpCmp, err := newAmqpComponent(cfg.Name, amqpURL, amqpQueue, amqpExchange)
	if err != nil {
		logger.Fatalf("failed to create processor %v", err)
	}

	kafkaCmp, err := newKafkaComponent(cfg.Name, kafkaBroker, kafkaTopic, amqpURL, amqpExchange)
	if err != nil {
		logger.Fatalf("failed to create processor %v", err)
	}

	httpCmp, err := newHTTPComponent(kafkaBroker, kafkaTopic)
	if err != nil {
		logger.Fatalf("failed to create processor %v", err)
	}

	// Set up routes
	routes := []http.Route{
		http.NewPostRoute("/", httpCmp.process, true),
	}

	srv, err := patron.New(cfg, patron.Routes(routes), patron.Components(kafkaCmp.cmp, amqpCmp.cmp))
	if err != nil {
		logger.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		logger.Fatalf("failed to create service %v", err)
	}

}
