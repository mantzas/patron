package confluent

import (
	"context"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

// KafkaProducer defines a async Kafka producer.
type KafkaProducer struct {
	cfg   *kafka.ConfigMap
	prod  *kafka.Producer
	tag   opentracing.Tag
	sync  bool
	chErr chan error
}

// NewSyncProducer creates a new sync producer with default configuration.
func NewSyncProducer(brokers []string, oo ...OptionFunc) (*KafkaProducer, error) {
	return newProducer(brokers, true, nil, oo...)
}

// NewAsyncProducer creates a new async producer with default configuration.
func NewAsyncProducer(brokers []string, ch chan error, oo ...OptionFunc) (*KafkaProducer, error) {
	return newProducer(brokers, false, ch, oo...)
}

func newProducer(brokers []string, sync bool, chErr chan error, oo ...OptionFunc) (*KafkaProducer, error) {

	if !sync && chErr == nil {
		return nil, errors.New("error chan needed for async producer")
	}

	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
	}

	sp := KafkaProducer{cfg: cfg, tag: opentracing.Tag{Key: "type", Value: "sync"}, sync: sync, chErr: chErr}

	for _, o := range oo {
		err := o(&sp)
		if err != nil {
			return nil, err
		}
	}

	p, err := kafka.NewProducer(sp.cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sync producer")
	}
	if !sync {

	}
	sp.prod = p
	return &sp, nil
}

// Send a message to a topic.
func (kp *KafkaProducer) Send(ctx context.Context, msg *Message) error {
	csp, _ := trace.ChildSpan(
		ctx,
		trace.ComponentOpName(trace.KafkaAsyncProducerComponent, msg.topic),
		trace.KafkaAsyncProducerComponent,
		ext.SpanKindProducer,
		kp.tag,
		opentracing.Tag{Key: "topic", Value: msg.topic},
	)
	pm, err := createProducerMessage(msg, csp)
	if err != nil {
		trace.SpanError(csp)
		return err
	}

	if kp.sync {
		err = kp.sendSync(pm)
	} else {
		err = kp.sendAsync(pm)
	}

	if err != nil {
		trace.SpanError(csp)
		return errors.Wrap(err, "failed to produce message")
	}
	trace.SpanSuccess(csp)
	return nil
}

func (kp *KafkaProducer) sendSync(msg *kafka.Message) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err := kp.prod.Produce(msg, deliveryChan)
	if err != nil {
		return errors.Wrap(err, "failed to produce message")
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return errors.Wrap(err, "failed to deliver message")
	}
	return nil
}

func (kp *KafkaProducer) sendAsync(msg *kafka.Message) error {
	return nil
}

func (kp *KafkaProducer) monitorAsyncErrors(chErr chan error) {
	go func() {
		for e := range kp.prod.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error == nil {
					continue
				}
				chErr <- errors.Wrap(m.TopicPartition.Error, "failed to produce message")
			}
		}
	}()
}

func createProducerMessage(msg *Message, sp opentracing.Span) (*kafka.Message, error) {
	c := kafkaHeadersCarrier{}
	err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, &c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to inject tracing headers")
	}
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &msg.topic, Partition: kafka.PartitionAny},
		Value:          msg.body,
		Headers:        c,
	}, nil
}

type kafkaHeadersCarrier []kafka.Header

// Set implements Set() of opentracing.TextMapWriter.
func (c *kafkaHeadersCarrier) Set(key, val string) {
	*c = append(*c, kafka.Header{Key: key, Value: []byte(val)})
}
