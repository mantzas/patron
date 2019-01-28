package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mantzas/patron/async/kafka"
)

const (
	kafkaTopic  = "patron-topic"
	kafkaBroker = "localhost:9092"
)

func main() {
	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, os.Interrupt, syscall.SIGTERM)
	cf, err := kafka.New("test", []string{kafkaTopic}, []string{kafkaBroker})
	if err != nil {
		log.Fatal(err)
	}
	cns, err := cf.Create()
	if err != nil {
		log.Fatal(err)
	}
	ch, chErr, err := cns.Consume(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	var now time.Time

	for {
		select {
		case <-termSig:
			dur := time.Since(now)
			rps := float64(count) / dur.Seconds()
			log.Printf("count: %d dur: %v rps: %f", count, dur, rps)
			log.Print("exiting")
			cns.Close()
			return
		case <-ch:
			if count == 0 {
				now = time.Now()
				log.Print("starting")
			}
			count++
			if count%1000 == 0 {
				log.Printf("%d messages received", count)
			}
		case err = <-chErr:
			log.Fatal(err)
		}
	}
}
