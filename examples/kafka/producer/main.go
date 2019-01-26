package main

import (
	"context"
	"log"
	"time"

	"github.com/mantzas/patron/trace/kafka"
)

const (
	kafkaTopic  = "patron-topic"
	kafkaBroker = "localhost:9092"
)

func main() {
	var err error
	prd, err := kafka.NewProducer([]string{kafkaBroker})
	if err != nil {
		log.Fatal(err)
	}
	err = prd.Send(context.Background(), kafkaTopic, "test")
	if err != nil {
		log.Fatal(err)
	}

	count := 100000
	now := time.Now()
	for i := 0; i < count; i++ {
		err = prd.Send(context.Background(), kafkaTopic, "test")
		if err != nil {
			log.Fatal(err)
		}
	}
	dur := time.Since(now)
	rps := float64(count) / dur.Seconds()
	log.Printf("count: %d dur: %v rps: %f", count, dur, rps)
}
