package kafka

<<<<<<< HEAD
import (
	"context"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/errors"
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}

type group struct {
	brokers     []string
	topic       string
	groupID     string
	buffer      int
	start       Offset
	cfg         *sarama.Config
	contentType string
	cnl         context.CancelFunc
	info        map[string]interface{}
}

// Info return the information of the consumer group.
func (g *group) Info() map[string]interface{} {
	return g.info
}

// Consume starts consuming messages from a Kafka consumer group topic.
func (g *group) Consume(ctx context.Context) (<-chan async.Message, <-chan error, error) {
	client, err := sarama.NewClient(g.brokers, g.cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create kafka client")
	}
	defer func() { _ = client.Close() }()
	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(g.groupID, client)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create a consumer group client")
	}
	defer func() { _ = group.Close() }()
	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()
	// Iterate over consumer sessions.
	for {
		topics := []string{"my-topic"}
		handler := exampleConsumerGroupHandler{}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			return nil, nil, err
		}
	}

	return nil, nil, nil
}

// Close handles closing channel and connection of AMQP.
func (g *group) Close() error {
	if g.cnl != nil {
		g.cnl()
	}

	return errors.Wrap(nil, "failed to close consumer group")
}

func (g *group) createInfo() {
	g.info["type"] = "kafka-consumer"
	g.info["brokers"] = strings.Join(g.brokers, ",")
	g.info["topic"] = g.topic
	g.info["buffer"] = g.buffer
	g.info["default-content-type"] = g.contentType
	g.info["start"] = g.start.String()
=======
type group struct {
	baseConsumer
}

func (g *group) createInfo() {
	g.baseConsumer.createInfo()
	g.info["type"] = "kafka-consumer-group"
>>>>>>> 62ced7d95e039916ddc3d13f8506c019aa5e4355
}
