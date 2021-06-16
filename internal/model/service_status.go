package model

import "errors"

// ServiceStatus defines the service current status
type ServiceStatus string

const (
	// ServiceStatusHealthy service is responding to healtchecks
	ServiceStatusHealthy ServiceStatus = "healthy"
	// ServiceStatusIdle service is not responding to healthchecks
	ServiceStatusIdle ServiceStatus = "idle"
)

// Validate validates a service status
func (ss ServiceStatus) Validate() error {
	switch ss {
	case ServiceStatusHealthy, ServiceStatusIdle:
		return nil
	}
	return ErrInvalidServiceStatus
}

// MarshalBinary serializes a service status
func (ss ServiceStatus) MarshalBinary() ([]byte, error) {
	return []byte(ss), nil
}

// ErrInvalidServiceStatus error for unknown service statuses
var ErrInvalidServiceStatus = errors.New("Invalid service status")
