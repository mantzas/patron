package kafka

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

var (
	tracingTypeSync  = opentracing.Tag{Key: "type", Value: "sync"}
	tracingTypeAsync = opentracing.Tag{Key: "type", Value: "async"}
)

// Message abstraction of a Kafka message.
type Message struct {
	topic string
	body  []byte
}

// NewMessage creates a new message.
func NewMessage(topic string, body []byte) *Message {
	return &Message{topic: topic, body: body}
}

// NewJSONMessage creates a new message with a JSON encoded body.
func NewJSONMessage(topic string, d interface{}) (*Message, error) {
	b, err := json.Encode(d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to JSON encode")
	}
	return &Message{topic: topic, body: b}, nil
}

// Result describes the result of a sent message.
type Result struct {
	Err       error
	Topic     string
	Partition int32
	Offset    int64
}

// Sender interface for Kafka.
type Sender interface {
	Send(ctx context.Context, msg *Message) error
	Results() <-chan *Result
	Close()
}

// Producer defines a async Kafka producer.
type Producer struct {
	cfg   *kafka.ConfigMap
	prod  *kafka.Producer
	tag   opentracing.Tag
	chRes chan *Result
}

// NewProducer creates a new sync producer with default configuration.
func NewProducer(brokers []string, oo ...OptionFunc) (*Producer, error) {
	return newProducer(brokers, tracingTypeSync, oo...)
}

// NewAsyncProducer creates a new async producer with default configuration.
func NewAsyncProducer(brokers []string, oo ...OptionFunc) (*Producer, error) {
	p, err := newProducer(brokers, tracingTypeAsync, oo...)
	if err != nil {
		return nil, err
	}
	p.chRes = make(chan *Result)
	go p.monitorErrorEvents()
	return p, nil
}

func newProducer(brokers []string, tag opentracing.Tag, oo ...OptionFunc) (*Producer, error) {

	if len(brokers) == 0 {
		return nil, errors.New("at least one broker must be provided")
	}

	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
	}

	kp := Producer{cfg: cfg, tag: tag}

	for _, o := range oo {
		err := o(&kp)
		if err != nil {
			return nil, err
		}
	}

	p, err := kafka.NewProducer(kp.cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sync producer")
	}
	kp.prod = p
	return &kp, nil
}

// Send a message to a topic.
func (p *Producer) Send(ctx context.Context, msg *Message) error {
	var err error
	csp, _ := trace.ChildSpan(
		ctx,
		trace.ComponentOpName(trace.KafkaAsyncProducerComponent, msg.topic),
		trace.KafkaAsyncProducerComponent,
		ext.SpanKindProducer,
		p.tag,
		opentracing.Tag{Key: "topic", Value: msg.topic},
	)
	pm, err := createProducerMessage(msg, csp)
	if err != nil {
		trace.SpanError(csp)
		return err
	}

	// checking the channel to determine if the producer is sync or async.
	if p.chRes == nil {
		err = p.sendSync(pm)
	} else {
		err = p.sendAsync(pm)
	}

	if err != nil {
		trace.SpanError(csp)
		return errors.Wrap(err, "failed to produce message")
	}
	trace.SpanSuccess(csp)
	return nil
}

func (p *Producer) sendSync(msg *kafka.Message) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err := p.prod.Produce(msg, deliveryChan)
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

func (p *Producer) sendAsync(msg *kafka.Message) error {
	p.prod.ProduceChannel() <- msg
	return nil
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

func (p *Producer) monitorErrorEvents() {
	go func() {
		for e := range p.prod.Events() {
			msg, ok := e.(*kafka.Message)
			if !ok {
				continue
			}
			if msg.TopicPartition.Error == nil {
				p.chRes <- &Result{
					Topic:     *msg.TopicPartition.Topic,
					Partition: msg.TopicPartition.Partition,
					Offset:    int64(msg.TopicPartition.Offset),
				}
			} else {
				p.chRes <- &Result{Err: msg.TopicPartition.Error}
			}
		}
	}()
}

// Results returns a result channel for monitoring published messages.
func (p *Producer) Results() <-chan *Result {
	return p.chRes
}

// Close the producer.
func (p *Producer) Close() {
	if p.prod == nil {
		return
	}
	p.prod.Close()
	if p.chRes != nil {
		close(p.chRes)
	}
}

type kafkaHeadersCarrier []kafka.Header

// Set implements Set() of opentracing.TextMapWriter.
func (c *kafkaHeadersCarrier) Set(key, val string) {
	*c = append(*c, kafka.Header{Key: key, Value: []byte(val)})
}
