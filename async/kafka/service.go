package kafka

import (
	"context"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

// Service implementation of a kafka consumer
type Service struct {
	p       async.Processor
	brokers []string
	topics  []string
	cfg     *sarama.Config
	ms      sarama.Consumer
}

// New returns a new client
func New(p async.Processor, clientID string, brokers []string, topics []string) (*Service, error) {
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

	return &Service{p, brokers, topics, config, nil}, nil
}

// Run starts the async processing
func (s *Service) Run(ctx context.Context) error {

	ms, err := sarama.NewConsumer(s.brokers, s.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create consumer")
	}
	s.ms = ms

	chMsg, chErr, err := s.consumers()
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
					err := s.p.Process(ctx, msg.Value)
					if err != nil {
						failCh <- errors.Wrap(err, "failed to process message")
					}
				}()
			case errMsg := <-chErr:
				failCh <- errors.Wrap(errMsg, "an error occurred during consumption")
			}
		}
	}()

	return <-failCh
}

// Shutdown the service
func (s *Service) Shutdown(ctx context.Context) error {
	return errors.Wrap(s.ms.Close(), "failed to close consumer")
}

func (s *Service) consumers() (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError, error) {
	chMsg := make(chan *sarama.ConsumerMessage)
	chErr := make(chan *sarama.ConsumerError)

	for _, topic := range s.topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}

		partitions, err := s.ms.Partitions(topic)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get partitions")
		}

		consumer, err := s.ms.ConsumePartition(topic, partitions[0], sarama.OffsetOldest)
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
