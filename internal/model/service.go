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
	u, err := url.Parse(b.URL)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return errors.New("invalid backend url")
	}
	return nil
}

// Service defines the relevant info about a microservice
type Service struct {
	ID      string        `json:"id"`
	Match   Match         `json:"match"`
	Backend Backend       `json:"backend"`
	Status  ServiceStatus `json:"status"`
	Cache   CacheConfig   `json:"cache"`
}

// HealthURL returns the URL used for health check for all service types
func (s Service) HealthURL() (*url.URL, error) {
	healthPath := s.Backend.HealthPath
	if healthPath == "" {
		healthPath = "/health"
	}
	return url.Parse(fmt.Sprintf("%s/%s", s.Backend.URL, healthPath))
}

// IsHealthy returns whether a service is healthy
func (s Service) IsHealthy() bool {
	return s.Status == ServiceStatusHealthy
}

func (s Service) MatchRequest(r *http.Request) bool {
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
	if s.ID == "" {
		return errors.New("service id is mandatory")
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
