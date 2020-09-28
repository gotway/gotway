package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthPath(t *testing.T) {
	tests := []struct {
		name     string
		service  Service
		wantPath string
		wantErr  error
	}{
		{
			name: "Health path for unknown service",
			service: Service{
				Type: "foo",
			},
			wantPath: "",
			wantErr:  ErrInvalidServiceType,
		},
		{
			name: "Default health path for REST service",
			service: Service{
				Type: ServiceTypeREST,
			},
			wantPath: "health",
			wantErr:  nil,
		},
		{
			name: "Custom health path for REST service",
			service: Service{
				Type:       ServiceTypeREST,
				HealthPath: "ping",
			},
			wantPath: "ping",
			wantErr:  nil,
		},
		{
			name: "Default health path for gRPC service",
			service: Service{
				Type: ServiceTypeGRPC,
			},
			wantPath: "grpc.health.v1.Health/Check",
			wantErr:  nil,
		},
		{
			name: "Custom health path for gRPC service",
			service: Service{
				Type: ServiceTypeGRPC,
				Path: "ping",
			},
			wantPath: "grpc.health.v1.Health/Check",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := tt.service.HealthPathForType()

			assert.Equal(t, path, tt.wantPath)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestIsHealthy(t *testing.T) {
	tests := []struct {
		name          string
		service       Service
		wantIsHealthy bool
	}{
		{
			name: "Healthy status for a healthy service",
			service: Service{
				Status: ServiceStatusHealthy,
			},
			wantIsHealthy: true,
		},
		{
			name: "Healthy status for a idle service",
			service: Service{
				Status: ServiceStatusIdle,
			},
			wantIsHealthy: false,
		},
		{
			name: "Healthy status for a unknown status service",
			service: Service{
				Status: "foo",
			},
			wantIsHealthy: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isHealthy := tt.service.IsHealthy()

			assert.Equal(t, isHealthy, tt.wantIsHealthy)
		})
	}
}

func TestValidateService(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr error
	}{
		{
			name: "Validate valid service",
			service: Service{
				Type:   ServiceTypeREST,
				Status: ServiceStatusHealthy,
				URL:    "http://foo.bar",
				Path:   "foo",
			},
			wantErr: nil,
		},
		{
			name: "Validate service with invalid type",
			service: Service{
				Type: "foo",
			},
			wantErr: ErrInvalidServiceType,
		},
		{
			name: "Validate service with invalid status",
			service: Service{
				Type:   ServiceTypeREST,
				Status: "bar",
			},
			wantErr: ErrInvalidServiceStatus,
		},
		{
			name: "Validate service with invalid URL",
			service: Service{
				Type:   ServiceTypeREST,
				Status: ServiceStatusHealthy,
				URL:    "",
			},
			wantErr: errInvalidField("url"),
		},
		{
			name: "Validate service with invalid path",
			service: Service{
				Type:   ServiceTypeREST,
				Status: ServiceStatusHealthy,
				URL:    "http://foo.bar",
				Path:   "",
			},
			wantErr: errInvalidField("path"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()

			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestValidateServiceDetail(t *testing.T) {
	tests := []struct {
		name          string
		serviceDetail ServiceDetail
		wantErr       error
	}{
		{
			name: "Validate valid service detail",
			serviceDetail: ServiceDetail{
				Service: Service{
					Type:   ServiceTypeREST,
					Status: ServiceStatusHealthy,
					URL:    "http://foo.bar",
					Path:   "foo",
				},
				Cache: CacheConfig{
					TTL:      1,
					Statuses: []int{200},
					Tags:     []string{"foo"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Validate service detail with invalid service",
			serviceDetail: ServiceDetail{
				Service: Service{
					Type: "foo",
				},
				Cache: CacheConfig{
					TTL:      1,
					Statuses: []int{200},
					Tags:     []string{"foo"},
				},
			},
			wantErr: ErrInvalidServiceType,
		},
		{
			name: "Validate service detail with invalid cache",
			serviceDetail: ServiceDetail{
				Service: Service{
					Type:   ServiceTypeREST,
					Status: ServiceStatusHealthy,
					URL:    "http://foo.bar",
					Path:   "foo",
				},
				Cache: CacheConfig{
					TTL:      0,
					Statuses: []int{},
					Tags:     []string{"foo"},
				},
			},
			wantErr: ErrInvalidCacheConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.serviceDetail.Validate()

			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestStatusMarshal(t *testing.T) {
	status := ServiceStatus(ServiceStatusHealthy)
	bytes, err := status.MarshalBinary()

	if err != nil {
		t.Errorf("Got unexpected error: %w", err)
	}
	assert.Equal(t, string(bytes), string(ServiceStatusHealthy))
}

func TestTypesMarshal(t *testing.T) {
	serviceType := ServiceType(ServiceTypeREST)
	bytes, err := serviceType.MarshalBinary()

	if err != nil {
		t.Errorf("Got unexpected error: %w", err)
	}
	assert.Equal(t, string(bytes), string(ServiceTypeREST))
}
