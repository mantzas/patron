package patron

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mantzas/patron/log"
)

// Service base component
type Service struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

// New creates a new base service
func New() *Service {

	ctx, cancel := context.WithCancel(context.Background())
	s := Service{ctx, cancel}
	s.setupTermSignal()
	return &s
}

func (s *Service) setupTermSignal() {
	go func() {

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Info("term signal received, cancelling")
		s.Cancel()
	}()
}
