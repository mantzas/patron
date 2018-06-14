package kafka

import (
	"context"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"
)

// Component implementation of a kafka consumer.
type Component struct {
	name    string
	p       async.Processor
	brokers []string
	topics  []string
	cfg     *sarama.Config
	ms      sarama.Consumer
}

// New returns a new component.
func New(name string, p async.Processor, clientID string, brokers []string, topics []string) (*Component, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if p == nil {
		return nil, errors.New("work processor is required")
	}

	if clientID == "" {
		return nil, errors.New("client id is required")
	}

	if len(brokers) == 0 {
		return nil, errors.New("provide at least one broker")
	}

	if len(topics) == 0 {
		return nil, errors.New("provide at least one topic")
	}

	config := sarama.NewConfig()
	config.ClientID = clientID
	config.Consumer.Return.Errors = true

	return &Component{name, p, brokers, topics, config, nil}, nil
}

// Run starts the async processing.
func (c *Component) Run(ctx context.Context) error {

	ms, err := sarama.NewConsumer(c.brokers, c.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create consumer")
	}
	c.ms = ms

	chMsg, chErr, err := c.consumers()
	if err != nil {
		return errors.Wrap(err, "failed to get consumers")
	}

	failCh := make(chan error)
	go func() {
		for {
			select {
			case msg := <-chMsg:
				log.Debugf("data received from topic %s", msg.Topic)
				go func() {
					sp := trace.StartConsumerSpan(c.name, trace.KafkaConsumerComponent, mapHeader(msg.Headers))

					ct, err := determineContentType(msg.Headers)
					if err != nil {
						failCh <- errors.Wrap(err, "failed to determine content type")
						trace.FinishConsumerSpan(sp, true)
						return
					}

					dec, err := async.DetermineDecoder(ct)
					if err != nil {
						failCh <- errors.Wrapf(err, "failed to determine decoder for %s", ct)
						trace.FinishConsumerSpan(sp, true)
						return
					}

					err = c.p.Process(ctx, async.NewMessage(msg.Value, dec))
					if err != nil {
						failCh <- errors.Wrap(err, "failed to process message")
						trace.FinishConsumerSpan(sp, true)
						return
					}
					trace.FinishConsumerSpan(sp, false)
				}()
			case errMsg := <-chErr:
				failCh <- errors.Wrap(errMsg, "an error occurred during consumption")
			}
		}
	}()

	return <-failCh
}

// Shutdown the component.
func (c *Component) Shutdown(ctx context.Context) error {
	return errors.Wrap(c.ms.Close(), "failed to close consumer")
}

func (c *Component) consumers() (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError, error) {
	chMsg := make(chan *sarama.ConsumerMessage)
	chErr := make(chan *sarama.ConsumerError)

	for _, topic := range c.topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}

		partitions, err := c.ms.Partitions(topic)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get partitions")
		}

		consumer, err := c.ms.ConsumePartition(topic, partitions[0], sarama.OffsetOldest)
		if nil != err {
			return nil, nil, errors.Wrap(err, "failed to get partition consumer")
		}

		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					chErr <- consumerError

				case msg := <-consumer.Messages():
					chMsg <- msg
				}
			}
		}(topic, consumer)
	}

	return chMsg, chErr, nil
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
