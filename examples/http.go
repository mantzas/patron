package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mantzas/patron/sync"
	tracehttp "github.com/mantzas/patron/trace/http"
	"github.com/mantzas/patron/trace/kafka"
	"github.com/pkg/errors"
)

type httpComponent struct {
	prd   kafka.Producer
	topic string
}

func newHTTPComponent(kafkabroker, topic string) (*httpComponent, error) {

	prd, err := kafka.NewAsyncProducer([]string{kafkaBroker})
	if err != nil {
		return nil, err
	}
	return &httpComponent{prd: prd, topic: topic}, nil
}

func (hc *httpComponent) process(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	aud := Audit{Name: "HTTP component", Started: time.Now()}
	googleReq, err := http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create requestfor www.google.com")
	}
	rsp, err := tracehttp.NewClient(5*time.Second).Do(ctx, googleReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get www.google.com")
	}

	ads := Audits{}
	ads.append(aud)

	kafkaMsg, err := kafka.NewJSONMessage(hc.topic, ads)
	if err != nil {
		return nil, err
	}

	err = hc.prd.Send(ctx, kafkaMsg)
	if err != nil {
		return nil, err
	}

	return sync.NewResponse(fmt.Sprintf("got %s from google", rsp.Status)), nil
}
