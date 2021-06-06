package client

import (
	"errors"
	"net/url"

	"github.com/gotway/gotway/internal/core"
)

// Client interface
type Client interface {
	getHealthURL() (*url.URL, error)
	HealthCheck() error
}

// New instanciates a new client
func New(service core.Service) (Client, error) {
	switch service.Type {
	case core.ServiceTypeREST:
		return clientREST{service}, nil
	case core.ServiceTypeGRPC:
		return clientGRPC{service}, nil
	default:
		return nil, core.ErrInvalidServiceType
	}
}

var errServiceNotAvailable = errors.New("Service not available")
