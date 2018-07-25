package async

import (
	"context"
	"sync"

	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

// Component implementation of a async component.
type Component struct {
	name string
	proc ProcessorFunc
	sync.Mutex
	cns Consumer
	cnl context.CancelFunc
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
	}, nil
}

// Run starts the consumer processing loop messages.
func (c *Component) Run(ctx context.Context) error {
	c.Lock()
	ctx, cnl := context.WithCancel(ctx)
	c.cnl = cnl
	c.Unlock()

	chMsg, chErr, err := c.cns.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get consumer channels")
	}

	failCh := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("canceling consuming messages requested")
				failCh <- nil
				return
			case msg := <-chMsg:
				go func() {
					err = c.proc(msg)
					if err != nil {
						agr := agr_errors.New()
						agr.Append(errors.Wrap(err, "failed to process message. Nack message"))
						agr.Append(errors.Wrap(msg.Nack(), "failed to NACK message"))
						failCh <- agr
						return
					}
					if err := msg.Ack(); err != nil {
						failCh <- err
					}
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
	c.Lock()
	defer c.Unlock()
	c.cnl()
	return c.cns.Close()
}
