package amqp

import (
	"context"
	"fmt"

	"github.com/mantzas/patron/async"
	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Component implementation of a AMQP client
type Component struct {
	url   string
	queue string
	p     async.Processor
	tag   string
	ch    *amqp.Channel
	conn  *amqp.Connection
}

// New returns a new client
func New(url, queue string, p async.Processor) (*Component, error) {

	if url == "" {
		return nil, errors.New("RabbitMQ url is required")
	}

	if queue == "" {
		return nil, errors.New("RabbitMQ queue name is required")
	}

	if p == nil {
		return nil, errors.New("work processor is required")
	}

	return &Component{url, queue, p, "", nil, nil}, nil
}

// Run starts the async processing.
func (s *Component) Run(ctx context.Context) error {

	conn, err := amqp.Dial(s.url)
	if err != nil {
		return errors.Wrapf(err, "failed to dial @ %s", s.url)
	}
	s.conn = conn

	ch, err := s.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed get channel")
	}
	s.ch = ch

	s.tag = uuid.New().String()
	log.Infof("consuming messages for tag %s", s.tag)

	deliveries, err := ch.Consume(s.queue, s.tag, false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed initialize consumer")
	}

	agr := agr_errors.New()

	select {
	case <-ctx.Done():
		log.Info("canceling requested")
		return nil
	case d := <-deliveries:
		log.Infof("processing message %s", d.MessageId)

		go func(d *amqp.Delivery, a *agr_errors.Aggregate) {
			dec, err := async.DetermineDecoder(d.ContentType)
			if err != nil {
				s.handlerMessageError(d, a, err, fmt.Sprintf("failed to determine encoding %s. Sending NACK", d.ContentType))
				return
			}
			err = s.p.Process(ctx, async.NewMessage(d.Body, dec))
			if err != nil {
				s.handlerMessageError(d, a, err, fmt.Sprintf("failed to process message %s. Sending NACK", d.MessageId))
				return
			}
			err = d.Ack(false)
			if err != nil {
				a.Append(errors.Wrapf(err, "failed to ACK message %s", d.MessageId))
				return
			}
		}(&d, agr)

		if agr.Count() > 0 {
			return agr
		}
	}

	return nil
}

// Shutdown the component.
func (s *Component) Shutdown(ctx context.Context) error {

	agr := agr_errors.New()

	if s.ch != nil {
		err := s.ch.Cancel(s.tag, true)
		agr.Append(errors.Wrapf(err, "failed to cancel channel of consumer %s", s.tag))
	}

	if s.conn != nil {
		err := s.conn.Close()
		agr.Append(errors.Wrap(err, "failed to close connection"))
	}

	if agr.Count() > 0 {
		return agr
	}
	return nil
}

func (s *Component) handlerMessageError(d *amqp.Delivery, a *agr_errors.Aggregate, err error, msg string) {
	a.Append(errors.Wrap(err, msg))
	err = d.Nack(false, true)
	if err != nil {
		a.Append(errors.Wrapf(err, "failed to NACK message", d.MessageId))
	}
}
