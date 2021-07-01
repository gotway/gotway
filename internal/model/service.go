package model

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Match struct {
	Scheme     string `json:"scheme"`
	Method     string `json:"method"`
	Host       string `json:"host"`
	Path       string `json:"path"`
	PathPrefix string `json:"pathPrefix"`
}

func (m Match) Validate() error {
	val := reflect.ValueOf(m)
	for i := 0; i < val.NumField(); i++ {
		if str := val.Field(i).String(); str != "" {
			return nil
		}
	}
	return errors.New("no match criterias provided")
}

type Backend struct {
	URL        string `json:"url"`
	HealthPath string `json:"healthPath"`
}

func (b Backend) Validate() error {
	if b.URL == "" {
		return errors.New("backend URL is mandatory")
	}
	return nil
}

// Service defines the relevant info about a microservice
type Service struct {
	Type    ServiceType   `json:"type"`
	Name    string        `json:"name"`
	Match   Match         `json:"match"`
	Backend Backend       `json:"backend"`
	Status  ServiceStatus `json:"status"`
	Cache   CacheConfig   `json:"cache"`
}

// HealthURL returns the URL used for health check for all service types
func (s Service) HealthURL() (*url.URL, error) {
	switch s.Type {
	case ServiceTypeREST:
		healthPath := s.Backend.HealthPath
		if healthPath == "" {
			healthPath = "/health"
		}
		return url.Parse(fmt.Sprintf("%s/%s", s.Backend.URL, s.Backend.HealthPath))
	case ServiceTypeGRPC:
		return url.Parse(s.Backend.URL)
	default:
		return nil, ErrInvalidServiceType
	}
}

// IsHealthy returns whether a service is healthy
func (s Service) IsHealthy() bool {
	return s.Status == ServiceStatusHealthy
}

func (s Service) MatchRequest(r *http.Request) bool {
	if s.Match.Scheme != "" && s.Match.Scheme != r.URL.Scheme {
		return false
	}
	if s.Match.Method != "" && s.Match.Method != r.Method {
		return false
	}
	if s.Match.Host != "" && s.Match.Host != r.Host {
		return false
	}
	if s.Match.Path != "" && s.Match.Path != r.URL.RawPath {
		return false
	}
	if s.Match.PathPrefix != "" && !strings.HasPrefix(r.URL.RawPath, s.Match.PathPrefix) {
		return false
	}
	return true
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
	if err := s.Match.Validate(); err != nil {
		return err
	}
	if err := s.Backend.Validate(); err != nil {
		return err
	}
	if !s.Cache.IsEmpty() {
		if err := s.Cache.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ErrServiceNotFound error for not found service
var ErrServiceNotFound = errors.New("Service not found")

// ErrServiceAlreadyRegistered error for service already registered
var ErrServiceAlreadyRegistered = errors.New("Service is already registered")
