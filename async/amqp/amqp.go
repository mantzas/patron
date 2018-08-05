package amqp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/encoding"
	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var (
	defaultCfg = amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 30*time.Second)
		},
	}
)

type message struct {
	span    opentracing.Span
	ctx     context.Context
	del     *amqp.Delivery
	dec     encoding.DecodeRawFunc
	requeue bool
}

func (m *message) Context() context.Context {
	return m.ctx
}

func (m *message) Decode(v interface{}) error {
	return m.dec(m.del.Body, v)
}

func (m *message) Ack() error {
	err := m.del.Ack(false)
	trace.SpanSuccess(m.span)
	return err
}

func (m *message) Nack() error {
	err := m.del.Nack(false, m.requeue)
	trace.SpanError(m.span)
	return err
}

// Consumer defines a AMQP subscriber.
type Consumer struct {
	url      string
	queue    string
	exchange string
	requeue  bool
	tag      string
	buffer   int
	cfg      amqp.Config
	ch       *amqp.Channel
	conn     *amqp.Connection
}

// New creates a new AMQP consumer with some defaults. Use option to change.
func New(url, queue, exchange string, oo ...OptionFunc) (*Consumer, error) {

	if url == "" {
		return nil, errors.New("RabbitMQ url is required")
	}

	if queue == "" {
		return nil, errors.New("RabbitMQ queue name is required")
	}

	if exchange == "" {
		return nil, errors.New("RabbitMQ exchange name is required")
	}

	c := &Consumer{
		url:      url,
		queue:    queue,
		exchange: exchange,
		requeue:  true,
		cfg:      defaultCfg,
		buffer:   1000,
	}

	for _, o := range oo {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Consume starts of consuming a AMQP queue.
func (c *Consumer) Consume(ctx context.Context) (<-chan async.Message, <-chan error, error) {
	deliveries, err := c.consumer()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed initialize consumer")
	}

	chMsg := make(chan async.Message, c.buffer)
	chErr := make(chan error, c.buffer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("canceling consuming messages requested")
				return
			case d := <-deliveries:
				log.Debugf("processing message %d", d.DeliveryTag)
				sp, chCtx := trace.ConsumerSpan(
					ctx,
					fmt.Sprintf("%s %s", trace.AMQPConsumerComponent, c.queue),
					trace.AMQPConsumerComponent,
					mapHeader(d.Headers),
				)
				dec, err := async.DetermineDecoder(d.ContentType)
				if err != nil {
					agr := agr_errors.New()
					agr.Append(errors.Wrapf(err, "failed to determine encoding %s. Nack message", d.ContentType))
					agr.Append(errors.Wrap(d.Nack(false, c.requeue), "failed to NACK message"))
					trace.SpanError(sp)
					chErr <- agr
					return
				}

				chMsg <- &message{
					ctx:     chCtx,
					dec:     dec,
					del:     &d,
					span:    sp,
					requeue: c.requeue,
				}
			}
		}
	}()

	return chMsg, chErr, nil
}

// Close handles closing channel and connection of AMQP.
func (c *Consumer) Close() error {
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

func (c *Consumer) consumer() (<-chan amqp.Delivery, error) {
	conn, err := amqp.DialConfig(c.url, c.cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial @ %s", c.url)
	}
	c.conn = conn

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed get channel")
	}
	c.ch = ch

	c.tag = uuid.New().String()
<<<<<<< HEAD
	log.Infof("consuming messages for tag %s", c.tag)

	err = ch.ExchangeDeclare(c.exchange, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare exchange")
	}
=======
	c.log.Infof("consuming messages for queue '%s' tag '%s'", c.queue, c.tag)
>>>>>>> Minor fixes

	q, err := ch.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare queue")
	}

	err = ch.QueueBind(q.Name, "", c.exchange, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bind queue to exchange queue")
	}

	deliveries, err := ch.Consume(c.queue, c.tag, false, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize consumer")
	}

	return deliveries, nil
}

func mapHeader(hh amqp.Table) map[string]string {
	mp := make(map[string]string)
	for k, v := range hh {
		mp[k] = fmt.Sprint(v)
	}
	return mp
}
