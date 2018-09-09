package async

import (
	"context"

	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/info"
	"github.com/mantzas/patron/log"
)

// Component implementation of a async component.
type Component struct {
	proc         ProcessorFunc
	failStrategy FailStrategy
	cns          Consumer
	info         info.Component
}

// New returns a new async component. The default behavior is to return a error of failure.
// Use options to change the default behavior.
func New(p ProcessorFunc, cns Consumer, oo ...OptionFunc) (*Component, error) {
	if p == nil {
		return nil, errors.New("work processor is required")
	}

	if cns == nil {
		return nil, errors.New("consumer is required")
	}

	c := &Component{
		proc:         p,
		cns:          cns,
		failStrategy: NackExitStrategy,
	}

	for _, o := range oo {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Run starts the consumer processing loop messages.
func (c *Component) Run(ctx context.Context) error {
	chMsg, chErr, err := c.cns.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get consumer channels")
	}

	failCh := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("closing consumer")
				failCh <- c.cns.Close()
				return
			case msg := <-chMsg:
				log.Debug("New message from consumer arrived")
				go c.processMessage(msg, failCh)
			case errMsg := <-chErr:
				failCh <- errors.Wrap(errMsg, "an error occurred during message consumption")
				return
			}
		}
	}()

	return <-failCh
}

func (c *Component) processMessage(msg Message, ch chan error) {
	err := c.proc(msg)
	if err != nil {
		err := c.executeFailureStrategy(msg, err)
		if err != nil {
			ch <- err
		}
		return
	}
	if err := msg.Ack(); err != nil {
		ch <- err
	}
}

func (c *Component) executeFailureStrategy(msg Message, err error) error {
	log.Errorf("failed to process message, failure strategy executed: %v", err)
	switch c.failStrategy {
	case NackExitStrategy:
		return errors.Aggregate(err, errors.Wrap(msg.Nack(), "failed to NACK message"))
	case NackStrategy:
		err := msg.Nack()
		if err != nil {
			return errors.Wrap(err, "nack failed when executing failure strategy")
		}
	case AckStrategy:
		err := msg.Ack()
		if err != nil {
			return errors.Wrap(err, "ack failed when executing failure strategy")
		}
	default:
		return errors.New("invalid failure strategy")
	}
	return nil
}
