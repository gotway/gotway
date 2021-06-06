package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/log"
	controllerMocks "github.com/gotway/gotway/mocks/controller"
	modelMocks "github.com/gotway/gotway/mocks/model"
	"github.com/stretchr/testify/assert"
)

func TestGetServices(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	catalogPath := "catalog"
	stockPath := "stock"
	routePath := "route.Route"
	servicePaths := []string{catalogPath, stockPath, routePath}
	catalog := core.Service{
		Type: core.ServiceTypeREST,
		Path: catalogPath,
	}
	stock := core.Service{
		Type: core.ServiceTypeREST,
		Path: stockPath,
	}
	route := core.Service{
		Type: core.ServiceTypeGRPC,
		Path: routePath,
	}

	serviceRepository.On("GetAllServiceKeys").Return(servicePaths)
	serviceRepository.On("GetServices", catalogPath).Return(
		[]core.Service{catalog}, nil,
	)
	serviceRepository.On("GetServices", stockPath, routePath).Return(
		[]core.Service{stock, route}, nil,
	)
	serviceRepository.On("GetServices", catalogPath, stockPath, routePath).Return(
		[]core.Service{catalog, stock, route}, nil,
	)

	tests := []struct {
		name            string
		offset          int
		limit           int
		wantServicePage core.ServicePage
		wantErr         error
	}{
		{
			name:            "Get services with invalid offset",
			offset:          10,
			limit:           1,
			wantServicePage: core.ServicePage{},
			wantErr:         core.ErrServiceNotFound,
		},
		{
			name:            "Get empty range of services",
			offset:          0,
			limit:           0,
			wantServicePage: core.ServicePage{},
			wantErr:         core.ErrServiceNotFound,
		},
		{
			name:   "Get first service",
			offset: 0,
			limit:  1,
			wantServicePage: core.ServicePage{
				Services:   []core.Service{catalog},
				TotalCount: len(servicePaths),
			},
			wantErr: nil,
		},
		{
			name:   "Get las 2 services",
			offset: 1,
			limit:  2,
			wantServicePage: core.ServicePage{
				Services:   []core.Service{stock, route},
				TotalCount: len(servicePaths),
			},
			wantErr: nil,
		},
		{
			name:   "Get all services",
			offset: 0,
			limit:  10,
			wantServicePage: core.ServicePage{
				Services:   []core.Service{catalog, stock, route},
				TotalCount: len(servicePaths),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			servicePage, err := controller.GetServices(tt.offset, tt.limit)

			assert.Equal(t, tt.wantServicePage, servicePage)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetServicesRepoError(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	serviceRepository.On("GetAllServiceKeys").Return([]string{"foo"})
	repoErr := errors.New("Error getting services")
	serviceRepository.On("GetServices", mock.Anything).Return(
		[]core.Service{}, repoErr,
	)

	servicePage, err := controller.GetServices(0, 1)

	assert.Equal(t, core.ServicePage{}, servicePage)
	assert.Equal(t, repoErr, err)
	serviceRepository.AssertExpectations(t)
}

func TestRegisterService(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	service := core.Service{
		Type: core.ServiceTypeREST,
		Path: "service",
	}
	cacheConfig := core.CacheConfig{
		TTL:      1,
		Statuses: []int{200},
		Tags:     []string{"catalog"},
	}
	serviceDetail := core.ServiceDetail{
		Service: service,
		Cache:   cacheConfig,
	}

	serviceRepository.On("StoreService", serviceDetail).Return(nil)

	err := controller.RegisterService(serviceDetail)

	assert.Nil(t, err)
}

func TestGetService(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	service := core.Service{
		Type: core.ServiceTypeREST,
		Path: "catalog",
	}
	serviceRepository.On("GetService", service.Path).Return(service, nil)

	result, err := controller.GetService(service.Path)

	assert.Equal(t, result, service)
	assert.Nil(t, err)
}

func TestGetServiceDetail(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	catalog := core.Service{
		Type: core.ServiceTypeREST,
		Path: "catalog",
	}
	cacheConfig := core.CacheConfig{
		TTL:      1,
		Statuses: []int{200},
		Tags:     []string{"catalog"},
	}
	serviceDetail := core.ServiceDetail{
		Service: catalog,
		Cache:   cacheConfig,
	}

	serviceRepository.On("GetServiceDetail", catalog.Path).Return(serviceDetail, nil)

	result, err := controller.GetServiceDetail(catalog.Path)

	assert.Equal(t, serviceDetail, result)
	assert.Nil(t, err)
}

func TestDeleteService(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	service := "service"

	serviceRepository.On("DeleteService", service).Return(nil)

	err := controller.DeleteService(service)

	assert.Nil(t, err)
}

func TestUpdateServiceStatus(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	serviceRepository.On("UpdateServiceStatus", mock.Anything, mock.Anything).Return(nil)

	err := controller.UpdateServiceStatus("catalog", core.ServiceStatusHealthy)
	assert.Nil(t, err)

	err = controller.UpdateServiceStatus("stock", core.ServiceStatusIdle)
	assert.Nil(t, err)

	serviceRepository.AssertNumberOfCalls(t, "UpdateServiceStatus", 2)
}

func TestReverseProxy(t *testing.T) {
	log.Init()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalog := core.Service{
		Type: core.ServiceTypeREST,
		URL:  "http://api.catalog.com",
		Path: "catalog",
	}
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%v/products", catalog.URL),
		httpmock.NewStringResponder(200, `[{"id": 1, "name": "T-Shirt"}]`))

	stock := core.Service{
		Type: core.ServiceTypeREST,
		URL:  "http://api.stock.com",
		Path: "stock",
	}
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%v/list", stock.URL),
		httpmock.NewStringResponder(200, `[{"id": 1, "stock": 10}]`))

	serviceRepository := new(modelMocks.ServiceRepositoryI)
	controller := newServiceController(serviceRepository)

	cache := new(controllerMocks.CacheControllerI)
	Cache = cache

	cache.On("HandleResponse", catalog.Path, mock.Anything).Return(nil)
	errHandleResponse := errors.New("Response error")
	cache.On("HandleResponse", stock.Path, mock.Anything).Return(errHandleResponse)

	catalogReq, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/catalog/products", nil)
	stockReq, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/stock/list", nil)

	tests := []struct {
		name           string
		req            *http.Request
		service        core.Service
		wantErr        error
		wantStatusCode int
	}{
		{
			name:           "Error creating reverse proxy",
			req:            nil,
			service:        core.Service{Type: "foo"},
			wantErr:        core.ErrInvalidServiceType,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Reverse proxy succeeds",
			req:            catalogReq,
			service:        catalog,
			wantErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Reverse proxy fails",
			req:            stockReq,
			service:        stock,
			wantErr:        nil,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := controller.ReverseProxy(rr, tt.req, tt.service)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantStatusCode, rr.Code)
		})
	}
}
