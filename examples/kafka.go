package main

import (
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
}

func newKafkaComponent(name, broker, topic, amqpURL, amqpExc string) (*kafkaComponent, error) {

	kafkaCmp := kafkaComponent{}

	cf, err := kafka.New(name, json.Type, topic, []string{broker})
	if err != nil {
		return nil, err
	}

	cmp, err := async.New(kafkaCmp.Process, cf, async.ConsumerRetry(10, 5*time.Second))
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

func (kc *kafkaComponent) Process(msg async.Message) error {
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

	err = kc.pub.Publish(msg.Context(), amqpMsg)
	if err != nil {
		return err
	}

	for _, a := range ads {
		log.Infof("%s@ took %s", a.Name, a.Duration)
	}

	return nil
}
