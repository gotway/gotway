package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		name        string
		service     Service
		wantIsValid bool
	}{
		{
			name: "Validate valid service",
			service: Service{
				ID: "foo",
				Match: Match{
					Host: "foo",
				},
				Backend: Backend{
					URL: "http://foo.bar",
				},
				Status: ServiceStatusHealthy,
			},
			wantIsValid: true,
		},
		{
			name: "Validate service with invalid status",
			service: Service{
				ID: "foo",
				Match: Match{
					Host: "foo",
				},
				Backend: Backend{
					URL: "http://foo.bar",
				},
				Status: "foo",
			},
			wantIsValid: false,
		},
		{
			name: "Validate service with invalid match",
			service: Service{
				ID:    "foo",
				Match: Match{},
				Backend: Backend{
					URL: "http://foo.bar",
				},
			},
			wantIsValid: false,
		},
		{
			name: "Validate service with invalid backend",
			service: Service{
				ID: "foo",
				Match: Match{
					Host: "foo",
				},
				Backend: Backend{},
			},
			wantIsValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.service.Validate() == nil, tt.wantIsValid)
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
