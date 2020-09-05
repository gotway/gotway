package model

import (
	"errors"
	"fmt"
)

// ServiceType defines the services type
type ServiceType string

const (
	// REST service type
	REST ServiceType = "rest"
	// GRPC service type
	GRPC ServiceType = "grpc"
)

// ErrInvalidServiceType error for unknown service types
var ErrInvalidServiceType = errors.New("Invalid service type")

func (st ServiceType) validate() error {
	switch st {
	case REST, GRPC:
		return nil
	}
	return ErrInvalidServiceType
}

// MarshalBinary serializes a service type
func (st ServiceType) MarshalBinary() ([]byte, error) {
	return []byte(st), nil
}

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
	Type       ServiceType   `json:"type"`
	URL        string        `json:"url"`
	Path       string        `json:"path"`
	HealthPath string        `json:"healthPath"`
	Status     ServiceStatus `json:"status"`
}

// HealthPathForType returns the path used for health check for all service types
func (s Service) HealthPathForType() (*string, error) {
	switch s.Type {
	case REST:
		var path string
		if s.HealthPath != "" {
			path = s.HealthPath
		} else {
			path = "health"
		}
		return &path, nil
	case GRPC:
		path := "grpc.health.v1.Health/Check"
		return &path, nil
	default:
		return nil, ErrInvalidServiceType
	}
}

// IsHealthy returns whether a service is healthy
func (s Service) IsHealthy() bool {
	return s.Status == Healthy
}

// Validate checks whether a service is valid
func (s Service) Validate() error {
	err := s.Type.validate()
	if err != nil {
		return err
	}
	if s.URL == "" {
		return errInvalidField("url")
	}
	if s.Path == "" {
		return errInvalidField("path")
	}
	return nil
}

func errInvalidField(f string) error {
	return fmt.Errorf("Invalid field '%s'", f)
}

// ErrServiceNotFound error for not found service
var ErrServiceNotFound = errors.New("Service not found")

// ErrServiceAlreadyRegistered error for service already registered
var ErrServiceAlreadyRegistered = errors.New("Service is already registered")
