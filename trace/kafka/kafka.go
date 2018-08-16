package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
)

// Message abstraction of a Kafka message.
type Message struct {
	topic string
	body  []byte
}

// NewMessage creates a new message.
func NewMessage(t string, b []byte) *Message {
	return &Message{topic: t, body: b}
}

// NewJSONMessage creates a new message with a JSON encoded body.
func NewJSONMessage(t string, d interface{}) (*Message, error) {

	b, err := json.Encode(d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to JSON encode")
	}
	return &Message{topic: t, body: b}, nil
}

// Producer interface for Kafka.
type Producer interface {
	Send(ctx context.Context, msg *Message) error
	Error() <-chan error
	Close() error
}

// AsyncProducer defines a async Kafka producer.
type AsyncProducer struct {
	prod  sarama.AsyncProducer
	chErr chan error
	tag   opentracing.Tag
}

// NewAsyncProducer creates a new async producer with default configuration.
func NewAsyncProducer(brokers []string, version string) (*AsyncProducer, error) {
	var v sarama.KafkaVersion
	var err error
	if version == "" {
		v = sarama.V0_11_0_0
	} else {
		v, err = sarama.ParseKafkaVersion(version)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse kafka version")
		}
	}

	config := sarama.NewConfig()
	config.Version = v

	prod, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sync producer")
	}
	ap := AsyncProducer{prod: prod, chErr: make(chan error), tag: opentracing.Tag{Key: "type", Value: "async"}}
	go ap.propagateError()
	return &ap, nil
}

// Send a message to a topic.
func (ap *AsyncProducer) Send(ctx context.Context, msg *Message) error {
	sp, _ := trace.ChildSpan(
		ctx,
		trace.ComponentOpName(trace.KafkaAsyncProducerComponent, msg.topic),
		trace.KafkaAsyncProducerComponent,
		ext.SpanKindProducer,
		ap.tag,
		opentracing.Tag{Key: "topic", Value: msg.topic},
	)
	pm, err := createProducerMessage(msg, sp)
	if err != nil {
		trace.SpanError(sp)
		return err
	}
	ap.prod.Input() <- pm
	trace.SpanSuccess(sp)
	return nil
}

// Error returns a chanel to monitor for errors.
func (ap *AsyncProducer) Error() <-chan error {
	return ap.chErr
}

// Close gracefully the producer.
func (ap *AsyncProducer) Close() error {
	return errors.Wrap(ap.prod.Close(), "failed to close sync producer")
}

func (ap *AsyncProducer) propagateError() {
	for pe := range ap.prod.Errors() {
		ap.chErr <- errors.Wrap(pe, "failed to send message")
	}
}

func createProducerMessage(msg *Message, sp opentracing.Span) (*sarama.ProducerMessage, error) {
	c := kafkaHeadersCarrier{}
	err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, &c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to inject tracing headers")
	}
	return &sarama.ProducerMessage{
		Topic:   msg.topic,
		Key:     nil,
		Value:   sarama.ByteEncoder(msg.body),
		Headers: c,
	}, nil
}

type kafkaHeadersCarrier []sarama.RecordHeader

// Set implements Set() of opentracing.TextMapWriter.
func (c *kafkaHeadersCarrier) Set(key, val string) {
	*c = append(*c, sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
}
