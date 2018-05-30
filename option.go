package patron

import (
	"time"
)

const (
	minReportingPeriod = 1 * time.Second
)

// Option defines a option for the HTTP service
type Option func(*Server) error
