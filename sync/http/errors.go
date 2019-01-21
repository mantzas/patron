package http

// ValidationError defines a validation error.
type ValidationError struct {
	err string
}

func (e *ValidationError) Error() string {
	return e.err
}

// NewValidationError creates a new validation error.
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{err: msg}
}

// UnauthorizedError defines a authorization error.
type UnauthorizedError struct {
	err string
}

func (e *UnauthorizedError) Error() string {
	return e.err
}

// NewUnauthorizedError creates a new unauthorized error.
func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{err: msg}
}

// ForbiddenError defines a access error.
type ForbiddenError struct {
	err string
}

func (e *ForbiddenError) Error() string {
	return e.err
}

// NewForbiddenError creates a new forbidden error.
func NewForbiddenError(msg string) *ForbiddenError {
	return &ForbiddenError{err: msg}
}

// NotFoundError defines a not found error.
type NotFoundError struct {
	err string
}

func (e *NotFoundError) Error() string {
	return e.err
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{err: msg}
}

// ServiceUnavailableError defines a service unavailable error.
type ServiceUnavailableError struct {
	err string
}

func (e *ServiceUnavailableError) Error() string {
	return e.err
}

// NewServiceUnavailableError creates a new service unavailable error.
func NewServiceUnavailableError(msg string) *ServiceUnavailableError {
	return &ServiceUnavailableError{err: msg}
}
