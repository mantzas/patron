package kafka

import (
	"context"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"
)

// Component implementation of a kafka consumer.
type Component struct {
	name        string
	proc        async.ProcessorFunc
	brokers     []string
	topics      []string
	cfg         *sarama.Config
	contentType string
	sync.Mutex
	ms sarama.Consumer
}

// New returns a new kafka consumer component.
func New(name string, p async.ProcessorFunc, clientID, ct string, brokers, topics []string) (*Component, error) {
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

	return &Component{
		name:        name,
		proc:        p,
		brokers:     brokers,
		topics:      topics,
		cfg:         config,
		ms:          nil,
		contentType: ct,
	}, nil
}

// Run starts the kafka consumer processing messages.
func (c *Component) Run(ctx context.Context) error {

	ms, err := sarama.NewConsumer(c.brokers, c.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create consumer")
	}
	c.Lock()
	c.ms = ms
	c.Unlock()

	chMsg, chErr, err := c.consumers()
	if err != nil {
		return errors.Wrap(err, "failed to get consumers")
	}

	failCh := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				failCh <- errors.New("canceling requested")
				return
			case msg := <-chMsg:
				log.Debugf("data received from topic %s", msg.Topic)
				go func() {
					sp, chCtx := trace.StartConsumerSpan(ctx, c.name, trace.KafkaConsumerComponent,
						mapHeader(msg.Headers))

					var ct string
					if c.contentType != "" {
						ct = c.contentType
					} else {
						ct, err = determineContentType(msg.Headers)
						if err != nil {
							failCh <- errors.Wrap(err, "failed to determine content type")
							trace.FinishSpanWithError(sp)
							return
						}
					}

					dec, err := async.DetermineDecoder(ct)
					if err != nil {
						failCh <- errors.Wrapf(err, "failed to determine decoder for %s", ct)
						trace.FinishSpanWithError(sp)
						return
					}

					err = c.proc(chCtx, async.NewMessage(msg.Value, dec))
					if err != nil {
						failCh <- errors.Wrap(err, "failed to process message")
						trace.FinishSpanWithError(sp)
						return
					}
					trace.FinishSpanWithSuccess(sp)
				}()
			case errMsg := <-chErr:
				failCh <- errors.Wrap(errMsg, "an error occurred during consumption")
				return
			}
		}
	}()

	return <-failCh
}

// Shutdown gracefully the component by closing the kafka consumer.
func (c *Component) Shutdown(ctx context.Context) error {
	c.Lock()
	c.Unlock()
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
