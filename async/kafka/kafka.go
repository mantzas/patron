package kafka

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type message struct {
	span opentracing.Span
	ctx  context.Context
	dec  encoding.DecodeRawFunc
	val  []byte
}

func (m *message) Context() context.Context {
	return m.ctx
}

func (m *message) Decode(v interface{}) error {
	return m.dec(m.val, v)
}

func (m *message) Ack() error {
	trace.FinishSpanWithSuccess(m.span)
	return nil
}

func (m *message) Nack() error {
	trace.FinishSpanWithError(m.span)
	return nil
}

// Offset defines the offset of messages inside a topic.
type Offset int64

const (
	// OffsetNewest starts consuming from the newest available message in the topic.
	OffsetNewest Offset = -1
	// OffsetOldest starts consuming from the oldest available message in the topic.
	OffsetOldest Offset = -2
)

// Consumer definition of a Kafka consumer.
type Consumer struct {
	name        string
	brokers     []string
	topic       string
	buffer      int
	start       Offset
	cfg         *sarama.Config
	contentType string
	cnl         context.CancelFunc
	ms          sarama.Consumer
	log         log.Logger
}

// New creates a ew Kafka consumer.
func New(name, ct, topic string, brokers []string, buffer int, start Offset, dialTimeout time.Duration) (*Consumer, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(brokers) == 0 {
		return nil, errors.New("provide at least one broker")
	}

	if topic == "" {
		return nil, errors.New("topic is required")
	}

	if buffer < 0 {
		return nil, errors.New("buffer must greater or equal than 0")
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, errors.New("failed to get hostname")
	}

	config := sarama.NewConfig()
	config.ClientID = fmt.Sprintf("%s-%s", host, name)
	config.Consumer.Return.Errors = true
	config.Net.DialTimeout = dialTimeout

	return &Consumer{
		name:        name,
		brokers:     brokers,
		topic:       topic,
		cfg:         config,
		contentType: ct,
		buffer:      buffer,
		start:       start,
	}, nil
}

// Consume starts consuming messages from a Kafka topic.
func (c *Consumer) Consume(ctx context.Context) (<-chan async.Message, <-chan error, error) {
	c.log = log.Create()
	ctx, cnl := context.WithCancel(ctx)
	c.cnl = cnl

	pcs, err := c.consumers()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get partitions")
	}

	chMsg := make(chan async.Message, c.buffer)
	chErr := make(chan error, c.buffer)

	for _, pc := range pcs {
		go func(consumer sarama.PartitionConsumer) {
			for {
				select {
				case <-ctx.Done():
					c.log.Info("canceling consuming messages requested")
					return
				case consumerError := <-consumer.Errors():
					chErr <- consumerError
				case msg := <-consumer.Messages():
					c.log.Debugf("data received from topic %s", msg.Topic)
					go func() {
						sp, chCtx := trace.StartConsumerSpan(ctx, c.name, trace.KafkaConsumerComponent, mapHeader(msg.Headers))

						var ct string
						if c.contentType != "" {
							ct = c.contentType
						} else {
							ct, err = determineContentType(msg.Headers)
							if err != nil {
								trace.FinishSpanWithError(sp)
								chErr <- errors.Wrap(err, "failed to determine content type")
								return
							}
						}

						dec, err := async.DetermineDecoder(ct)
						if err != nil {
							trace.FinishSpanWithError(sp)
							chErr <- errors.Wrapf(err, "failed to determine decoder for %s", ct)
							return
						}

						chMsg <- &message{
							ctx:  chCtx,
							dec:  dec,
							span: sp,
							val:  msg.Value,
						}
					}()
				}
			}
		}(pc)
	}

	return chMsg, chErr, nil
}

// Close handles closing channel and connection of AMQP.
func (c *Consumer) Close() error {
	return errors.Wrap(c.ms.Close(), "failed to close consumer")
}

func (c *Consumer) consumers() ([]sarama.PartitionConsumer, error) {

	ms, err := sarama.NewConsumer(c.brokers, c.cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create consumer")
	}
	c.ms = ms

	partitions, err := c.ms.Partitions(c.topic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get partitions")
	}

	pcs := make([]sarama.PartitionConsumer, len(partitions))

	for i, partition := range partitions {

		pc, err := c.ms.ConsumePartition(c.topic, partition, int64(c.start))
		if nil != err {
			return nil, errors.Wrap(err, "failed to get partition consumer")
		}
		pcs[i] = pc
	}

	return pcs, nil
}

func determineContentType(hdr []*sarama.RecordHeader) (string, error) {
	for _, h := range hdr {
		if string(h.Key) == encoding.ContentTypeHeader {
			return string(h.Value), nil
		}
	}
	return "", errors.New("content type header is missing")
}

func mapHeader(hh []*sarama.RecordHeader) map[string]string {
	mp := make(map[string]string)
	for _, h := range hh {
		mp[string(h.Key)] = string(h.Value)
	}
	return mp
}
