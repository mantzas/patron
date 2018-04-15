package amqp

import (
	"context"

	"github.com/mantzas/patron/async"
	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Service implementation of a AMQP client
type Service struct {
	url   string
	queue string
	mp    async.MessageProcessor
	tag   string
	ch    *amqp.Channel
	conn  *amqp.Connection
}

// New returns a new client
func New(url, queue string, mp async.MessageProcessor) (*Service, error) {

	if url == "" {
		return nil, errors.New("rabbitmq url is required")
	}

	if queue == "" {
		return nil, errors.New("rabbitmq queue name is required")
	}

	if mp == nil {
		return nil, errors.New("work processor is required")
	}

	return &Service{url, queue, mp, "", nil, nil}, nil
}

// Run starts the async processing
func (s *Service) Run(ctx context.Context) error {

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

			err := s.mp.Process(ctx, d.Body)
			if err != nil {
				a.Append(errors.Wrapf(err, "failed to process message %s", d.MessageId))
				return
			}
			d.Ack(false)
		}(&d, agr)

		if agr.Count() > 0 {
			return agr
		}
	}

	return nil
}

// Shutdown the service
func (s *Service) Shutdown(ctx context.Context) error {

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
