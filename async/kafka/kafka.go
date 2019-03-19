package kafka

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
)

var topicPartitionOffsetDiff *prometheus.GaugeVec

func init() {
	topicPartitionOffsetDiff = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "component",
			Subsystem: "kafka_consumer",
			Name:      "offset_diff",
			Help:      "Message offset difference with high watermark, classified by topic and partition",
		},
		[]string{"topic", "partition"},
	)
	prometheus.MustRegister(topicPartitionOffsetDiff)
}

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
type Offset int64

const (
	// OffsetNewest starts consuming from the newest available message in the topic.
	OffsetNewest Offset = -1
	// OffsetOldest starts consuming from the oldest available message in the topic.
	OffsetOldest Offset = -2
)

func (o Offset) String() string {
	switch o {
	case OffsetNewest:
		return "OffsetNewest"
	case OffsetOldest:
		return "OffsetOldest"
	default:
		return strconv.FormatInt(int64(o), 10)
	}
}

// Factory definition of a consumer factory.
type Factory struct {
	name    string
	ct      string
	topic   string
	brokers []string
	oo      []OptionFunc
}

// New constructor.
func New(name, ct, topic string, brokers []string, oo ...OptionFunc) (*Factory, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(brokers) == 0 {
		return nil, errors.New("provide at least one broker")
	}

	if topic == "" {
		return nil, errors.New("topic is required")
	}

	return &Factory{name: name, ct: ct, topic: topic, brokers: brokers, oo: oo}, nil
}

// Create a new consumer.
func (f *Factory) Create() (async.Consumer, error) {

	host, err := os.Hostname()
	if err != nil {
		return nil, errors.New("failed to get hostname")
	}

	config := sarama.NewConfig()
	config.ClientID = fmt.Sprintf("%s-%s", host, f.name)
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_0

	c := &consumer{
		brokers:     f.brokers,
		topic:       f.topic,
		cfg:         config,
		contentType: f.ct,
		buffer:      1000,
		start:       OffsetNewest,
		info:        make(map[string]interface{}),
	}

	for _, o := range f.oo {
		err = o(c)
		if err != nil {
			return nil, err
		}
	}

	c.createInfo()
	return c, nil
}

type consumer struct {
	brokers     []string
	topic       string
	buffer      int
	start       Offset
	cfg         *sarama.Config
	contentType string
	cnl         context.CancelFunc
	ms          sarama.Consumer
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

	pcs, err := c.consumers()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get partitions")
	}
	log.Infof("consuming messages for topic '%s'", c.topic)
	chMsg := make(chan async.Message, c.buffer)
	chErr := make(chan error, c.buffer)

	for _, pc := range pcs {
		go func(consumer sarama.PartitionConsumer) {
			for {
				select {
				case <-ctx.Done():
					log.Info("canceling consuming messages requested")
					closeConsumer(consumer)
					return
				case consumerError := <-consumer.Errors():
					closeConsumer(consumer)
					chErr <- consumerError
					return
				case m := <-consumer.Messages():
					log.Debugf("data received from topic %s", m.Topic)
					topicPartitionOffsetDiffGaugeSet(m.Topic, m.Partition, consumer.HighWaterMarkOffset(), m.Offset)
					go func(msg *sarama.ConsumerMessage) {
						sp, chCtx := trace.ConsumerSpan(
							ctx,
							trace.ComponentOpName(trace.KafkaConsumerComponent, msg.Topic),
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
					}(m)
				}
			}
		}(pc)
	}

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

func (c *consumer) consumers() ([]sarama.PartitionConsumer, error) {

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

func (c *consumer) createInfo() {
	c.info["type"] = "kafka-consumer"
	c.info["brokers"] = strings.Join(c.brokers, ",")
	c.info["topic"] = c.topic
	c.info["buffer"] = c.buffer
	c.info["default-content-type"] = c.contentType
	c.info["start"] = c.start.String()
}

func closeConsumer(cns sarama.PartitionConsumer) {
	if cns == nil {
		return
	}
	err := cns.Close()
	if err != nil {
		log.Errorf("failed to close partition consumer: %v", err)
	}
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

func topicPartitionOffsetDiffGaugeSet(topic string, partition int32, high, offset int64) {
	topicPartitionOffsetDiff.WithLabelValues(topic, strconv.FormatInt(int64(partition), 10)).Set(float64(high - offset))
}
