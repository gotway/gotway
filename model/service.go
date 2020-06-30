package model

import "errors"

// ServiceStatus defines the service current status
type ServiceStatus string

const (
	// Healthy service is responding to healtchecks
	Healthy ServiceStatus = "healthy"
	// Idle service is not responding to healthchecks
	Idle ServiceStatus = "idle"
)

// ErrInvalidServiceStatus error for unknown service statuses
var ErrInvalidServiceStatus = errors.New("Invalid service status")

// IsValid checks whether a service status is valid
func (ss ServiceStatus) validate() error {
	switch ss {
	case Healthy, Idle:
		return nil
	}
	return ErrInvalidServiceStatus
}

// MarshalBinary serializes a service status
func (ss ServiceStatus) MarshalBinary() ([]byte, error) {
	return []byte(ss), nil
}

// Service defines the relevant info about a microservice
type Service struct {
	Key       string
	URL       string
	HealthURL string
	Status    ServiceStatus
}

// IsHealthy returns whether a service is healthy
func (s Service) IsHealthy() bool {
	return s.Status == Healthy
}

// ErrServiceNotFound error for not found service
var ErrServiceNotFound = errors.New("Service not found")

// ErrServiceAlreadyRegistered error for service already registered
var ErrServiceAlreadyRegistered = errors.New("Service is already registered")
