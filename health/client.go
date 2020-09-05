package health

import (
	"errors"
	"net/url"

	"github.com/gosmo-devs/microgateway/model"
)

// Client interface
type Client interface {
	getHealthURL() (*url.URL, error)
	HealthCheck() error
}

var errServiceNotAvailable = errors.New("Service not available")

// NewClient instanciates a new client
func NewClient(service *model.Service) (Client, error) {
	switch service.Type {
	case model.REST:
		return restClient{service}, nil
	case model.GRPC:
		return grpcClient{service}, nil
	default:
		return nil, model.ErrInvalidServiceType
	}
}
