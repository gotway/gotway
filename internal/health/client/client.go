package client

import (
	"errors"
	"net/url"
	"time"

	"github.com/gotway/gotway/internal/model"
)

type Options struct {
	Timeout time.Duration
}

// Client interface
type Client interface {
	HealthCheck(url *url.URL) error
	Release()
}

var ErrServiceNotAvailable = errors.New("Service not available")

// New instanciates a new client
func New(serviceType model.ServiceType, options Options) (Client, error) {
	switch serviceType {
	case model.ServiceTypeREST:
		return newClientREST(options), nil
	case model.ServiceTypeGRPC:
		return newClientGRPC(options), nil
	default:
		return nil, model.ErrInvalidServiceType
	}
}
