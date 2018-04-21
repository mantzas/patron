package patron

import "context"

// Service interface for implementing services
type Service interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
