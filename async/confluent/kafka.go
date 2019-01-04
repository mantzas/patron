package confluent

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/metric"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
)

var topicPartitionOffsetDiff *prometheus.GaugeVec

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
	trace.SpanSuccess(m.span)
	return nil
}

func (m *message) Nack() error {
	trace.SpanError(m.span)
	return nil
}

// Offset defines the offset of messages inside a topic.
type Offset string

const (
	// OffsetSmallest smallest offset.
	OffsetSmallest Offset = "smallest"
	// OffsetEarliest earliest offset.
	OffsetEarliest Offset = "earliest"
	// OffsetBeginning beginning offset.
	OffsetBeginning Offset = "beginning"
	// OffsetLargest largest offset.
	OffsetLargest Offset = "largest"
	// OffsetLatest latest offset.
	OffsetLatest Offset = "latest"
	// OffsetEnd end offset.
	OffsetEnd Offset = "end"
	// OffsetError error offset.
	OffsetError Offset = "error"
)

// Factory definition of a consumer factory.
type Factory struct {
	name    string
	ct      string
	topics  []string
	brokers []string
	oo      []OptionFunc
}

// New constructor.
func New(name, ct string, topics []string, brokers []string, oo ...OptionFunc) (*Factory, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(brokers) == 0 {
		return nil, errors.New("brokers are nil or empty")
	}

	if len(topics) == 0 {
		return nil, errors.New("topics are nil or empty")
	}

	return &Factory{name: name, ct: ct, topics: topics, brokers: brokers, oo: oo}, nil
}

// Create a new consumer.
func (f *Factory) Create() (async.Consumer, error) {

	host, err := os.Hostname()
	if err != nil {
		return nil, errors.New("failed to get hostname")
	}

	cfg := &kafka.ConfigMap{
		"client.id":                       fmt.Sprintf("%s-%s", host, f.name),
		"bootstrap.servers":               strings.Join(f.brokers, ","),
		"session.timeout.ms":              10000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"go.events.channel.size":          1000,
		"auto.offset.reset":               OffsetLatest,
	}

	c := &consumer{
		brokers:     f.brokers,
		topics:      f.topics,
		contentType: f.ct,
		buffer:      1000,
		info:        make(map[string]interface{}),
		cfg:         cfg,
	}

	for _, o := range f.oo {
		err = o(c)
		if err != nil {
			return nil, err
		}
	}

	err = setupMetrics()
	if err != nil {
		return nil, err
	}
	c.createInfo()
	return c, nil
}

type consumer struct {
	brokers     []string
	topic       string
	contentType string
	cnl         context.CancelFunc
	ms          sarama.Consumer
	cns         *kafka.Consumer
	cfg         *kafka.ConfigMap
	buffer      int
	topics      []string
	info        map[string]interface{}
}

// Info return the information of the consumer.
func (c *consumer) Info() map[string]interface{} {
	return c.info
}

// Consume starts consuming messages from a Kafka topic.
func (c *consumer) Consume(ctx context.Context) (<-chan async.Message, <-chan error, error) {
	ctx, cnl := context.WithCancel(ctx)
	c.cnl = cnl

	cns, err := kafka.NewConsumer(c.cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create new consumer")
	}
	c.cns = cns

	err = cns.SubscribeTopics(c.topics, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to subscribe to topics")
	}

	log.Infof("consuming messages for topic '%s'", c.topic)
	chMsg := make(chan async.Message, c.buffer)
	chErr := make(chan error)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("canceling consuming messages requested")
				closeConsumer(cns)
			case ev := <-cns.Events():
				switch e := ev.(type) {
				case kafka.AssignedPartitions:
					log.Infof("assigned partitions: %v", e)
					cns.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					log.Infof("revoking partitions: %v", e)
					cns.Unassign()
				case kafka.Error:
					// Errors should generally be considered as informational, the client will try to automatically recover
					log.Errorf("failure in message consumption: %v", e)
				case *kafka.Message:
					log.Debugf("data received from topic %d", e.TopicPartition)
					c.topicPartitionOffsetDiffGaugeSet(e.TopicPartition)
					go func(msg *kafka.Message) {
						sp, chCtx := trace.ConsumerSpan(
							ctx,
							trace.ComponentOpName(trace.KafkaConsumerComponent, *msg.TopicPartition.Topic),
							trace.KafkaConsumerComponent,
							mapHeader(msg.Headers),
						)
						var ct string
						if c.contentType != "" {
							ct = c.contentType
						} else {
							ct, err = determineContentType(msg.Headers)
							if err != nil {
								trace.SpanError(sp)
								chErr <- errors.Wrap(err, "failed to determine content type")
								return
							}
						}

						dec, err := async.DetermineDecoder(ct)
						if err != nil {
							trace.SpanError(sp)
							chErr <- errors.Wrapf(err, "failed to determine decoder for %s", ct)
							return
						}

						chCtx = log.WithContext(chCtx, log.Sub(map[string]interface{}{"messageID": uuid.New().String()}))

						chMsg <- &message{
							ctx:  chCtx,
							dec:  dec,
							span: sp,
							val:  msg.Value,
						}
					}(e)
				}
			}
		}
	}()

	return chMsg, chErr, nil
}

// Close handles closing consumer.
func (c *consumer) Close() error {
	if c.cnl != nil {
		c.cnl()
	}

	if c.ms == nil {
		return nil
	}

	return errors.Wrap(c.ms.Close(), "failed to close consumer")
}

func (c *consumer) createInfo() {
	c.info["type"] = "kafka-consumer"
	c.info["brokers"] = strings.Join(c.brokers, ",")
	c.info["topics"] = strings.Join(c.topics, ",")
	c.info["buffer"] = c.buffer
	c.info["content-type"] = c.contentType
	for k, v := range *c.cfg {
		c.info[k] = v
	}
}

func (c *consumer) topicPartitionOffsetDiffGaugeSet(tp kafka.TopicPartition) {
	_, high, err := c.cns.QueryWatermarkOffsets(*tp.Topic, tp.Partition, 10)
	if err != nil {
		log.Warnf("failed to query watermarks: %v", err)
	}
	topicPartitionOffsetDiff.WithLabelValues(*tp.Topic, strconv.FormatInt(int64(tp.Partition), 10)).Set(float64(high - int64(tp.Offset)))
}

func closeConsumer(cns *kafka.Consumer) {
	if cns == nil {
		return
	}
	err := cns.Close()
	if err != nil {
		log.Errorf("failed to close partition consumer: %v", err)
	}
}

func determineContentType(hdr []kafka.Header) (string, error) {
	for _, h := range hdr {
		if string(h.Key) == encoding.ContentTypeHeader {
			return string(h.Value), nil
		}
	}
	return "", errors.New("content type header is missing")
}

func mapHeader(hh []kafka.Header) map[string]string {
	mp := make(map[string]string)
	for _, h := range hh {
		mp[h.Key] = string(h.Value)
	}
	return mp
}

func setupMetrics() error {
	var err error
	topicPartitionOffsetDiff, err = metric.NewGauge(
		"kafka_consumer",
		"offset_diff",
		"Message offset difference with high watermark, classified by topic and partition",
		"topic",
		"partition",
	)
	if err != nil {
		return err
	}
	return nil
}
