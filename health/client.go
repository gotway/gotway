package health

import (
	"errors"
	"net/url"

	"github.com/gosmo-devs/microgateway/core"
)

// Client interface
type Client interface {
	getHealthURL() (*url.URL, error)
	HealthCheck() error
}

// NewClient instanciates a new client
func NewClient(service core.Service) (Client, error) {
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
