package model

import "errors"

// ServiceType defines the services type
type ServiceType string

const (
	// ServiceTypeREST service type
	ServiceTypeREST ServiceType = "rest"
	// ServiceTypeGRPC service type
	ServiceTypeGRPC ServiceType = "grpc"
)

// Validate validates a service type
func (st ServiceType) Validate() error {
	switch st {
	case ServiceTypeREST, ServiceTypeGRPC:
		return nil
	}
	return ErrInvalidServiceType
}

// MarshalBinary serializes a service type
func (st ServiceType) MarshalBinary() ([]byte, error) {
	return []byte(st), nil
}

// ErrInvalidServiceType error for unknown service types
var ErrInvalidServiceType = errors.New("Invalid service type")
