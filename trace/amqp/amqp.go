package amqp

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/mantzas/patron/encoding/json"
	patronerrors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/trace"
)

// Message for publishing.
type Message struct {
	contentType string
	body        []byte
}

// NewMessage creates a new message for publishing.
func NewMessage(ct string, body []byte) *Message {
	return &Message{contentType: ct, body: body}
}

// NewJSONMessage creates a new message for publishing.
func NewJSONMessage(d interface{}) (*Message, error) {
	body, err := json.Encode(d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal to JSON")
	}
	return &Message{contentType: json.ContentType, body: body}, nil
}

// Publisher interface of a RabbitMQ publisher.
type Publisher interface {
	Publish(ctx context.Context, msg *Message) error
	Close(ctx context.Context) error
}

// TracedPublisher defines a RabbitMQ publisher with integrated tracing.
type TracedPublisher struct {
	cn     *amqp.Connection
	ch     *amqp.Channel
	exc    string
	opName string
	tag    opentracing.Tag
}

// NewPublisher creates a new publisher.
// The default exchange type used is fanout.
// Notifications are not handled at this point TBD.
func NewPublisher(url, exc string) (*TracedPublisher, error) {

	if url == "" {
		return nil, errors.New("url is required")
	}

	if exc == "" {
		return nil, errors.New("exchange is required")
	}

	p := TracedPublisher{
		exc:    exc,
		opName: "amqp PUB exchange " + exc,
		tag:    opentracing.Tag{Key: "exchange", Value: exc},
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open RabbitMq connection")
	}
	p.cn = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open RabbitMq channel")
	}
	p.ch = ch

	err = ch.ExchangeDeclare(
		exc,                 // name
		amqp.ExchangeFanout, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare exchange")
	}

	return &p, nil
}

// Publish a payload to a exchange.
func (tc *TracedPublisher) Publish(ctx context.Context, msg *Message) error {
	sp, _ := trace.StartChildSpan(ctx, tc.opName, trace.AMQPPublisherComponent, tc.tag)

	p := amqp.Publishing{
		Headers:     amqp.Table{},
		ContentType: msg.contentType,
		Body:        msg.body,
	}

	c := amqpHeadersCarrier(p.Headers)
	sp.Tracer().Inject(sp.Context(), opentracing.TextMap, c)

	err := tc.ch.Publish(
		tc.exc, // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		p,
	)
	if err != nil {
		trace.FinishSpan(sp, true)
		return errors.Wrap(err, "failed to publish message")
	}
	trace.FinishSpan(sp, false)
	return nil
}

// Close the connection and channel of the publisher.
func (tc *TracedPublisher) Close(ctx context.Context) error {
	aggError := patronerrors.New()

	aggError.Append(tc.ch.Close())
	aggError.Append(tc.cn.Close())

	if aggError.Count() > 0 {
		return aggError
	}
	return nil
}

type amqpHeadersCarrier map[string]interface{}

// Set implements Set() of opentracing.TextMapWriter.
func (c amqpHeadersCarrier) Set(key, val string) {
	c[key] = val
}
