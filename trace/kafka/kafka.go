package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// Message definition for Kafka.
type Message struct {
	topic string
	body  []byte
}

// NewMessage creates a new Kafka message.
func NewMessage(t string, b []byte) *Message {
	return &Message{topic: t, body: b}
}

// NewJSONMessage creates a new Kafka JSON message from a model.
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

// AsyncProducer definition of a sync Kafka producer.
type AsyncProducer struct {
	prod  sarama.AsyncProducer
	chErr chan error
	tag   opentracing.Tag
}

// NewAsyncProducer creates a new Kafka sync producer with default config.
func NewAsyncProducer(brokers []string) (*AsyncProducer, error) {

	prod, err := sarama.NewAsyncProducer(brokers, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sync producer")
	}
	ap := AsyncProducer{prod: prod, chErr: make(chan error), tag: opentracing.Tag{Key: "type", Value: "async"}}
	go ap.propagateError()
	return &ap, nil
}

// SendMessage to a Kafka topic.
func (ap *AsyncProducer) SendMessage(ctx context.Context, msg *Message) {
	sp, _ := trace.StartChildSpan(ctx, "kafka PROD topic "+msg.topic, trace.KafkaAsyncProducerComponent,
		ap.tag, opentracing.Tag{Key: "topic", Value: msg.topic})
	defer trace.FinishSpan(sp, false)
	ap.prod.Input() <- createProducerMessage(msg, sp)
}

func (ap *AsyncProducer) Error() <-chan error {
	return ap.chErr
}

// Close a existing producer.
func (ap *AsyncProducer) Close() error {
	return errors.Wrap(ap.prod.Close(), "failed to close sync producer")
}

func (ap *AsyncProducer) propagateError() {
	for pe := range ap.prod.Errors() {
		ap.chErr <- errors.Wrap(pe, "failed to send message")
	}
}

func createProducerMessage(msg *Message, sp opentracing.Span) *sarama.ProducerMessage {
	c := kafkaHeadersCarrier{[]sarama.RecordHeader{}}
	sp.Tracer().Inject(sp.Context(), opentracing.TextMap, &c)
	return &sarama.ProducerMessage{
		Topic:   msg.topic,
		Key:     nil,
		Value:   sarama.ByteEncoder(msg.body),
		Headers: c.hdr,
	}
}

type kafkaHeadersCarrier struct {
	hdr []sarama.RecordHeader
}

// Set implements Set() of opentracing.TextMapWriter.
func (c *kafkaHeadersCarrier) Set(key, val string) {
	c.hdr = append(c.hdr, sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
}
