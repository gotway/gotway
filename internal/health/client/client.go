package client

import (
	"errors"
	"net/url"

	"github.com/gotway/gotway/internal/model"
)

// Client interface
type Client interface {
	getHealthURL() (*url.URL, error)
	HealthCheck() error
}

// New instanciates a new client
func New(service model.Service) (Client, error) {
	switch service.Type {
	case model.ServiceTypeREST:
		return clientREST{service}, nil
	case model.ServiceTypeGRPC:
		return clientGRPC{service}, nil
	default:
		return nil, model.ErrInvalidServiceType
	}
}

var errServiceNotAvailable = errors.New("Service not available")
