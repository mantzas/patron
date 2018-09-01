package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mantzas/patron/sync"
	tracehttp "github.com/mantzas/patron/trace/http"
	"github.com/mantzas/patron/trace/kafka"
	"github.com/mantzas/patron/errors"
)

type httpComponent struct {
	prd   kafka.Producer
	topic string
	url   string
}

func newHTTPComponent(kafkabroker, topic, url string) (*httpComponent, error) {
	prd, err := kafka.NewAsyncProducer([]string{kafkaBroker}, "")
	if err != nil {
		return nil, err
	}
	return &httpComponent{prd: prd, topic: topic, url: url}, nil
}

func (hc *httpComponent) first(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	aud := Audit{Name: "first HTTP component", Started: time.Now()}
	secondRouteReq, err := http.NewRequest("GET", hc.url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create requestfor www.google.com")
	}
	secondRouteReq.Header.Add("Content-Type", "application/json")
	secondRouteReq.Header.Add("Accept", "application/json")
	rsp, err := tracehttp.NewClient(5*time.Second).Do(ctx, secondRouteReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get www.google.com")
	}

	ads := Audits{}
	ads.append(aud)

	return sync.NewResponse(fmt.Sprintf("got %s from second HTTP route", rsp.Status)), nil
}

func (hc *httpComponent) second(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	aud := Audit{Name: "second HTTP component", Started: time.Now()}
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
