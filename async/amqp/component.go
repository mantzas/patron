package amqp

import (
	"context"
	"fmt"
	"sync"

	"github.com/mantzas/patron/async"
	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Component implementation of a AMQP subscriber.
type Component struct {
	name  string
	url   string
	queue string
	proc  async.ProcessorFunc
	tag   string
	sync.Mutex
	ch   *amqp.Channel
	conn *amqp.Connection
}

// New returns a new AMQP subscriber.
func New(name, url, queue string, p async.ProcessorFunc) (*Component, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if url == "" {
		return nil, errors.New("RabbitMQ url is required")
	}

	if queue == "" {
		return nil, errors.New("RabbitMQ queue name is required")
	}

	if p == nil {
		return nil, errors.New("work processor is required")
	}

	return &Component{name: name, url: url, queue: queue, proc: p, tag: "", ch: nil, conn: nil}, nil
}

// Run starts AMQP subscription and async processing of messages.
func (c *Component) Run(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()
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
			sp, chCtx := trace.StartConsumerSpan(ctx, c.name, trace.AMQPConsumerComponent, mapHeader(d.Headers))

			dec, err := async.DetermineDecoder(d.ContentType)
			if err != nil {
				handlerMessageError(d, a, err, fmt.Sprintf("failed to determine encoding %s. Sending NACK", d.ContentType))
				trace.FinishSpanWithError(sp)
				return
			}
			err = c.proc(chCtx, async.NewMessage(d.Body, dec))
			if err != nil {
				handlerMessageError(d, a, err, fmt.Sprintf("failed to process message %s. Sending NACK", d.MessageId))
				trace.FinishSpanWithError(sp)
				return
			}
			err = d.Ack(false)
			if err != nil {
				a.Append(errors.Wrapf(err, "failed to ACK message %s", d.MessageId))
				trace.FinishSpanWithError(sp)
				return
			}
			trace.FinishSpanWithSuccess(sp)
		}(&d, agr)

		if agr.Count() > 0 {
			return agr
		}
	}

	return nil
}

// Shutdown the component by closing gracefully AMQP channel and connection.
func (c *Component) Shutdown(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()
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

func mapHeader(hh amqp.Table) map[string]string {
	mp := make(map[string]string)
	for k, v := range hh {
		mp[k] = v.(string)
	}
	return mp
}
