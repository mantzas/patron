package main

import (
	"context"
	"time"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/async/amqp"
	"github.com/mantzas/patron/log"
)

type amqpComponent struct {
	cmp patron.Component
	log log.Logger
}

func newAmqpComponent(name, url, queue, exchange string) (*amqpComponent, error) {

	amqpCmp := amqpComponent{}

	cns, err := amqp.New(name, url, queue, exchange, true, 1000)
	if err != nil {
		return nil, err
	}
	cmp, err := async.New(name, amqpCmp.Process, cns)
	if err != nil {
		return nil, err
	}
	amqpCmp.cmp = cmp

	return &amqpCmp, nil
}

func (ac *amqpComponent) Process(ctx context.Context, msg async.Message) error {
	ac.log = log.Create()
	var ads Audits

	err := msg.Decode(&ads)
	if err != nil {
		return err
	}

	ads.append(Audit{Name: "AMQP consumer", Started: time.Now()})

	for _, a := range ads {
		ac.log.Infof("%s@ took %s", a.Name, a.Duration)
	}

	return nil
}
