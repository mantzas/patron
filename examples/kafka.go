package main

import (
	"context"
	"time"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/async/kafka"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace/amqp"
)

type kafkaComponent struct {
	cmp patron.Component
	pub amqp.Publisher
	log log.Logger
}

func newKafkaComponent(name, broker, topic, amqpURL, amqpExc string) (*kafkaComponent, error) {

	kafkaCmp := kafkaComponent{}

	cns, err := kafka.New(name, "kafka-consumer", json.ContentType, topic, []string{broker}, 1000, kafka.OffsetNewest)
	if err != nil {
		return nil, err
	}
	cmp, err := async.New(kafkaCmp.Process, cns)
	if err != nil {
		return nil, err
	}
	kafkaCmp.cmp = cmp

	pub, err := amqp.NewPublisher(amqpURL, amqpExc)
	if err != nil {
		return nil, err
	}
	kafkaCmp.pub = pub

	return &kafkaCmp, nil
}

func (kc *kafkaComponent) Process(ctx context.Context, msg async.Message) error {
	kc.log = log.Create()
	var ads Audits

	err := msg.Decode(&ads)
	if err != nil {
		return err
	}

	ads.append(Audit{Name: "Kafka consumer", Started: time.Now()})

	amqpMsg, err := amqp.NewJSONMessage(ads)
	if err != nil {
		return err
	}

	err = kc.pub.Publish(ctx, amqpMsg)
	if err != nil {
		return err
	}

	for _, a := range ads {
		kc.log.Infof("%s@ took %s", a.Name, a.Duration)
	}

	return nil
}
