package client

import (
	"errors"
	"net/url"

	"github.com/gotway/gotway/internal/model"
)

// Client interface
type Client interface {
	HealthCheck(url *url.URL) error
}

var ErrServiceNotAvailable = errors.New("Service not available")

// New instanciates a new client
func New(serviceType model.ServiceType) (Client, error) {
	switch serviceType {
	case model.ServiceTypeREST:
		return newClientREST(), nil
	case model.ServiceTypeGRPC:
		return newClientGRPC(), nil
	default:
		return nil, model.ErrInvalidServiceType
	}
}
