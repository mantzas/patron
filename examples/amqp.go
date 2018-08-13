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
}

func newAmqpComponent(name, url, queue, exchange string) (*amqpComponent, error) {

	amqpCmp := amqpComponent{}

	cns, err := amqp.New(name, url, queue, exchange)
	if err != nil {
		return nil, err
	}
	cmp, err := async.New(amqpCmp.Process, cns)
	if err != nil {
		return nil, err
	}
	amqpCmp.cmp = cmp

	return &amqpCmp, nil
}

func (ac *amqpComponent) Process(ctx context.Context, msg async.Message) error {
	var ads Audits

	err := msg.Decode(&ads)
	if err != nil {
		return err
	}

	ads.append(Audit{Name: "AMQP consumer", Started: time.Now()})

	for _, a := range ads {
		log.Infof("%s@ took %s", a.Name, a.Duration)
	}

	return nil
}
