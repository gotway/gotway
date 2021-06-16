package model

import (
	"fmt"
	"net/http"
	"net/url"
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

func TestGetServiceRelativePath(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/catalog/products", nil)
	pathReq, _ := http.NewRequest(http.MethodGet, "/catalog/products", nil)
	pathPrefixReq, _ := http.NewRequest(http.MethodGet, "/api/cache/catalog/products", nil)

	tests := []struct {
		name             string
		req              *http.Request
		pathPrefix       string
		servicePath      string
		wantRelativePath string
		wantErr          error
	}{
		{
			name:             "Service not found in URL error",
			req:              req,
			pathPrefix:       "",
			servicePath:      "foo",
			wantRelativePath: "",
			wantErr: &ErrServiceNotFoundInURL{
				URL:         req.URL,
				ServicePath: "foo",
			},
		},
		{
			name:             "Relative path with full URL",
			req:              req,
			pathPrefix:       "",
			servicePath:      "catalog",
			wantRelativePath: "/products",
			wantErr:          nil,
		},
		{
			name:             "Relative path with path URL",
			req:              pathReq,
			pathPrefix:       "",
			servicePath:      "catalog",
			wantRelativePath: "/products",
			wantErr:          nil,
		},
		{
			name:             "Relative path with path URL and path prefix",
			req:              pathPrefixReq,
			pathPrefix:       "api/cache",
			servicePath:      "catalog",
			wantRelativePath: "/products",
			wantErr:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relativePath string
			var err error
			if tt.pathPrefix != "" {
				relativePath, err = GetServiceRelativePathPrefixed(
					tt.req,
					tt.pathPrefix,
					tt.servicePath,
				)
			} else {
				relativePath, err = GetServiceRelativePath(tt.req, tt.servicePath)
			}

			assert.Equal(t, tt.wantRelativePath, relativePath)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestErrServiceNotFoundInURLFormat(t *testing.T) {
	urlString := "http://api.gotway.com/catalog/products"
	url, _ := url.Parse(urlString)
	servicePath := "foo"
	err := &ErrServiceNotFoundInURL{
		URL:         url,
		ServicePath: servicePath,
	}

	assert.EqualError(
		t,
		err,
		fmt.Sprintf("Service path '%s' not found in URL: %s", servicePath, urlString),
	)
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
