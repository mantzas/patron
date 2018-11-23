package kafka

import (
	"context"
	"strings"

	"github.com/Shopify/sarama"
)

type baseConsumer struct {
	brokers     []string
	topic       string
	buffer      int
	cfg         *sarama.Config
	contentType string
	cnl         context.CancelFunc
	info        map[string]interface{}
}

// Info return the information of the consumer.
func (b *baseConsumer) Info() map[string]interface{} {
	return b.info
}

func (b *baseConsumer) createInfo() {
	b.info["brokers"] = strings.Join(b.brokers, ",")
	b.info["topic"] = b.topic
	b.info["buffer"] = b.buffer
	b.info["default-content-type"] = b.contentType
}
