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
	amqpURL      = "amqp://guest:guest@localhost:5672/"
	amqpExchange = "patron"
	amqpQueue    = "patron"
	kafkaTopic   = "patron-topic"
	kafkaBroker  = "localhost:9092"
)

func init() {
	err := os.Setenv("PATRON_LOG_LEVEL", "debug")
	if err != nil {
		fmt.Printf("failed to set log level env var: %v", err)
		os.Exit(1)
	}
	err = os.Setenv("PATRON_JAEGER_SAMPLER_PARAM", "1.0")
	if err != nil {
		fmt.Printf("failed to set sampler env vars: %v", err)
		os.Exit(1)
	}
}

func main() {
	name := "patron"
	version := "1.0.0"

	err := patron.SetupLogging(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	amqpCmp, err := newAmqpComponent(amqpURL, amqpQueue, amqpExchange)
	if err != nil {
<<<<<<< HEAD
		log.Fatalf("failed to create processor %v", err)
=======
		log.Create().Fatalf("failed to create processor: %v", err)
>>>>>>> Fixed kafka producer header tracing
	}

	kafkaCmp, err := newKafkaComponent(name, kafkaBroker, kafkaTopic, amqpURL, amqpExchange)
	if err != nil {
<<<<<<< HEAD
		log.Fatalf("failed to create processor %v", err)
=======
		log.Create().Fatalf("failed to create processor: %v", err)
>>>>>>> Fixed kafka producer header tracing
	}

	httpCmp, err := newHTTPComponent(kafkaBroker, kafkaTopic, "http://localhost:50000/second")
	if err != nil {
<<<<<<< HEAD
		log.Fatalf("failed to create processor %v", err)
=======
		log.Create().Fatalf("failed to create processor: %v", err)
>>>>>>> Fixed kafka producer header tracing
	}

	// Set up routes
	routes := []http.Route{
		http.NewPostRoute("/", httpCmp.first, true),
		http.NewGetRoute("/second", httpCmp.second, true),
	}

	srv, err := patron.New(name, version, patron.Routes(routes), patron.Components(kafkaCmp.cmp, amqpCmp.cmp))
	if err != nil {
<<<<<<< HEAD
		log.Fatalf("failed to create service %v", err)
=======
		log.Create().Fatalf("failed to create service: %v", err)
>>>>>>> Fixed kafka producer header tracing
	}

	err = srv.Run()
	if err != nil {
<<<<<<< HEAD
		log.Fatalf("failed to run service %v", err)
=======
		log.Create().Fatalf("service failure: %v", err)
>>>>>>> Fixed kafka producer header tracing
	}
}
