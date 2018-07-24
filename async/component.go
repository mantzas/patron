package async

import (
	"context"

	"github.com/pkg/errors"
)

// MessageI interface for defining messages that are handled by the async component.
type MessageI interface {
	Context() context.Context
	Decode(v interface{}) error
	Ack() error
	Nack() error
}

// Consumer interface which every specific consumer has to implement.
type Consumer interface {
	Consume(context.Context) (<-chan MessageI, <-chan error, error)
	Close() error
}

// Component implementation of a async component.
type Component struct {
	name string
	proc ProcessorFunc
	cns  Consumer
	cnl  context.CancelFunc
}

// New returns a new async component.
func New(name string, p ProcessorFunc, cns Consumer) (*Component, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if p == nil {
		return nil, errors.New("work processor is required")
	}

	if cns == nil {
		return nil, errors.New("consumer is required")
	}

	return &Component{
		name: name,
		proc: p,
		cns:  cns,
		cnl:  nil,
	}, nil
}

// Run starts the consumer processing loop messages.
func (c *Component) Run(ctx context.Context) error {

	chCtx, cnl := context.WithCancel(ctx)
	c.cnl = cnl

	chMsg, chErr, err := c.cns.Consume(chCtx)
	if err != nil {
		return errors.Wrap(err, "failed to get consumer channels")
	}

	failCh := make(chan error)
	go func() {
		for {
			select {
			case <-chCtx.Done():
				failCh <- errors.Wrap(c.cns.Close(), "failed to close consumer")
				return
			case msg := <-chMsg:
				go func() {
					err = c.proc(msg)
					if err != nil {
						msg.Nack()
						failCh <- errors.Wrap(err, "failed to process message")
						return
					}
					msg.Ack()
				}()
			case errMsg := <-chErr:
				failCh <- errors.Wrap(errMsg, "an error occurred during message consumption")
				return
			}
		}
	}()

	return <-failCh
}

// Shutdown gracefully the component by closing the consumer.
func (c *Component) Shutdown(ctx context.Context) error {
	if c.cnl != nil {
		c.cnl()
	}
	return c.cns.Close()
}
