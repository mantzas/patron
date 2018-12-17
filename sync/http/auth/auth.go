package auth

import (
	"net/http"
)

// Authenticator interface.
type Authenticator interface {
	Authenticate(req *http.Request) (bool, error)
}
