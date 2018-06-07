package amqp

import (
	"context"
	"fmt"

	"github.com/mantzas/patron/async"
	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type amqpContextKey string

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
func (c *Component) Run(ctx context.Context, tr opentracing.Tracer) error {

	conn, err := amqp.Dial(c.url)
	if err != nil {
		return errors.Wrapf(err, "failed to dial @ %s", c.url)
	}
	c.conn = conn

	ch, err := c.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed get channel")
	}
	c.ch = ch

	c.tag = uuid.New().String()
	log.Infof("consuming messages for tag %s", c.tag)

	deliveries, err := ch.Consume(c.queue, c.tag, false, false, false, false, nil)
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

			chCtx, _ := createContext(ctx, d.Headers)

			dec, err := async.DetermineDecoder(d.ContentType)
			if err != nil {
				handlerMessageError(d, a, err, fmt.Sprintf("failed to determine encoding %s. Sending NACK", d.ContentType))
				return
			}
			err = c.p.Process(chCtx, async.NewMessage(d.Body, dec))
			if err != nil {
				handlerMessageError(d, a, err, fmt.Sprintf("failed to process message %s. Sending NACK", d.MessageId))
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
func (c *Component) Shutdown(ctx context.Context) error {

	agr := agr_errors.New()

	if c.ch != nil {
		err := c.ch.Cancel(c.tag, true)
		agr.Append(errors.Wrapf(err, "failed to cancel channel of consumer %s", c.tag))
	}

	if c.conn != nil {
		err := c.conn.Close()
		agr.Append(errors.Wrap(err, "failed to close connection"))
	}

	if agr.Count() > 0 {
		return agr
	}
	return nil
}

func handlerMessageError(d *amqp.Delivery, a *agr_errors.Aggregate, err error, msg string) {
	a.Append(errors.Wrap(err, msg))
	err = d.Nack(false, true)
	if err != nil {
		a.Append(errors.Wrapf(err, "failed to NACK message %s", d.MessageId))
	}
}

func createContext(ctx context.Context, hdr amqp.Table) (context.Context, context.CancelFunc) {
	chCtx, cnl := context.WithCancel(ctx)

	for k, v := range hdr {
		chCtx = context.WithValue(chCtx, amqpContextKey(k), v)
	}

	return chCtx, cnl
}
