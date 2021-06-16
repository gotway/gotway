package model

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Service defines the relevant info about a microservice
type Service struct {
	Type       ServiceType   `json:"type"`
	URL        string        `json:"url"`
	Path       string        `json:"path"`
	HealthPath string        `json:"healthPath"`
	Status     ServiceStatus `json:"status"`
}

// HealthPathForType returns the path used for health check for all service types
func (s Service) HealthPathForType() (string, error) {
	switch s.Type {
	case ServiceTypeREST:
		var path string
		if s.HealthPath != "" {
			path = s.HealthPath
		} else {
			path = "health"
		}
		return path, nil
	case ServiceTypeGRPC:
		path := "grpc.health.v1.Health/Check"
		return path, nil
	default:
		return "", ErrInvalidServiceType
	}
}

// IsHealthy returns whether a service is healthy
func (s Service) IsHealthy() bool {
	return s.Status == ServiceStatusHealthy
}

// Validate checks whether a service is valid
func (s Service) Validate() error {
	err := s.Type.Validate()
	if err != nil {
		return err
	}
	if s.Status != "" {
		if err := s.Status.Validate(); err != nil {
			return err
		}
	}
	if s.URL == "" {
		return errInvalidField("url")
	}
	if s.Path == "" {
		return errInvalidField("path")
	}
	return nil
}

// ServiceDetail model
type ServiceDetail struct {
	Service
	Cache CacheConfig `json:"cache"`
}

// Validate checks whether a service detail is valid
func (sd ServiceDetail) Validate() error {
	err := sd.Service.Validate()
	if err != nil {
		return err
	}
	return sd.Cache.Validate()
}

// ServicePage model
type ServicePage struct {
	Services   []Service `json:"services"`
	TotalCount int       `json:"totalCount"`
}

// GetServiceRelativePathPrefixed retrieves the relative path of a service that has a prefix
func GetServiceRelativePathPrefixed(
	r *http.Request,
	pathPrefix string,
	servicePath string,
) (string, error) {
	var b strings.Builder
	if r.URL.Scheme != "" && r.URL.Host != "" {
		root := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)
		b.WriteString(root)
	}
	if pathPrefix != "" {
		b.WriteString(fmt.Sprintf("/%s", pathPrefix))
	}
	if servicePath != "" {
		b.WriteString(fmt.Sprintf("/%s", servicePath))
	}

	urlString := r.URL.String()
	prefix := b.String()

	if !strings.HasPrefix(urlString, prefix) {
		return "", &ErrServiceNotFoundInURL{URL: r.URL, ServicePath: servicePath}
	}

	return strings.TrimPrefix(urlString, prefix), nil
}

// GetServiceRelativePath retrieves the relative path of a service
func GetServiceRelativePath(r *http.Request, servicePath string) (string, error) {
	return GetServiceRelativePathPrefixed(r, "", servicePath)
}

// ErrServiceNotFound error for not found service
var ErrServiceNotFound = errors.New("Service not found")

// ErrServiceAlreadyRegistered error for service already registered
var ErrServiceAlreadyRegistered = errors.New("Service is already registered")

// ErrServiceNotFoundInURL is returned when a service is not found in a URL
type ErrServiceNotFoundInURL struct {
	URL         *url.URL
	ServicePath string
}

func (e *ErrServiceNotFoundInURL) Error() string {
	return fmt.Sprintf("Service path '%s' not found in URL: %s", e.ServicePath, e.URL.String())
}

func errInvalidField(f string) error {
	return fmt.Errorf("Invalid field '%s'", f)
}
