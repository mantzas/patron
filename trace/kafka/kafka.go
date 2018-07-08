package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
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
	SendMessage(ctx context.Context, msg *Message) error
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
func NewAsyncProducer(brokers []string) (*AsyncProducer, error) {

	prod, err := sarama.NewAsyncProducer(brokers, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sync producer")
	}
	ap := AsyncProducer{prod: prod, chErr: make(chan error), tag: opentracing.Tag{Key: "type", Value: "async"}}
	go ap.propagateError()
	return &ap, nil
}

// Send a message to a topic.
func (ap *AsyncProducer) Send(ctx context.Context, msg *Message) error {
	sp, _ := trace.StartChildSpan(ctx, "kafka PROD topic "+msg.topic, trace.KafkaAsyncProducerComponent,
		ap.tag, opentracing.Tag{Key: "topic", Value: msg.topic})
	pm, err := createProducerMessage(msg, sp)
	if err != nil {
		trace.FinishSpanWithError(sp)
		return err
	}
	ap.prod.Input() <- pm
	trace.FinishSpanWithSuccess(sp)
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
	c := kafkaHeadersCarrier{hdr: []sarama.RecordHeader{}}
	err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, &c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to inject tracing headers")
	}
	return &sarama.ProducerMessage{
		Topic:   msg.topic,
		Key:     nil,
		Value:   sarama.ByteEncoder(msg.body),
		Headers: c.hdr,
	}, nil
}

type kafkaHeadersCarrier struct {
	hdr []sarama.RecordHeader
}

// Set implements Set() of opentracing.TextMapWriter.
func (c *kafkaHeadersCarrier) Set(key, val string) {
	c.hdr = append(c.hdr, sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
}
