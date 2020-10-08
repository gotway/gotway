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
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

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
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

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
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

	service := core.Service{
		Type: core.ServiceTypeREST,
		Path: "service",
	}
	invalidService := core.Service{
		Type: "foo",
		Path: "foo",
	}
	errService := errors.New("Error storing service")
	serviceRepository.On("StoreService", service).Return(nil)
	serviceRepository.On("StoreService", invalidService).Return(errService)

	cacheConfig := core.CacheConfig{
		TTL:      1,
		Statuses: []int{200},
		Tags:     []string{"catalog"},
	}
	invalidCacheConfig := core.CacheConfig{
		TTL:      -1,
		Statuses: []int{418},
	}
	errCacheConfig := errors.New("Error storing cache config")
	cacheConfigRepository.On("StoreConfig", cacheConfig, service.Path).Return(nil)
	cacheConfigRepository.On("StoreConfig", invalidCacheConfig, service.Path).Return(errCacheConfig)
	cacheConfigRepository.On("StoreConfig", invalidCacheConfig, invalidService.Path).Return(errCacheConfig)

	tests := []struct {
		name          string
		serviceDetail core.ServiceDetail
		wantErr       error
	}{
		{
			name: "Stores a invalid service",
			serviceDetail: core.ServiceDetail{
				Service: invalidService,
				Cache:   cacheConfig,
			},
			wantErr: errService,
		},
		{
			name: "Stores a service with invalid config",
			serviceDetail: core.ServiceDetail{
				Service: service,
				Cache:   invalidCacheConfig,
			},
			wantErr: errCacheConfig,
		},
		{
			name: "Stores a service without cache",
			serviceDetail: core.ServiceDetail{
				Service: service,
				Cache:   core.DefaultCacheConfig,
			},
			wantErr: nil,
		},
		{
			name: "Stores a service",
			serviceDetail: core.ServiceDetail{
				Service: service,
				Cache:   cacheConfig,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := controller.RegisterService(tt.serviceDetail)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetService(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

	catalog := core.Service{
		Type: core.ServiceTypeREST,
		Path: "catalog",
	}
	serviceRepository.On("GetService", catalog.Path).Return(catalog, nil)

	errServiceKey := "foo"
	errService := errors.New("Error getting service")
	serviceRepository.On("GetService", errServiceKey).Return(core.Service{}, errService)

	tests := []struct {
		name        string
		key         string
		wantService core.Service
		wantErr     error
	}{
		{
			name:        "Get an invalid service",
			key:         errServiceKey,
			wantService: core.Service{},
			wantErr:     errService,
		},
		{
			name:        "Get a service",
			key:         catalog.Path,
			wantService: catalog,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := controller.GetService(tt.key)

			assert.Equal(t, tt.wantService, service)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetServiceDetail(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

	catalog := core.Service{
		Type: core.ServiceTypeREST,
		Path: "catalog",
	}
	stock := core.Service{
		Type: core.ServiceTypeREST,
		Path: "stock",
	}
	cacheConfig := core.CacheConfig{
		TTL:      1,
		Statuses: []int{200},
		Tags:     []string{"catalog"},
	}

	errServiceKey := "foo"
	errService := errors.New("Error getting service")
	errCacheServiceKey := "bar"
	errCache := errors.New("Error getting cache")

	serviceRepository.On("GetService", catalog.Path).Return(catalog, nil)
	serviceRepository.On("GetService", stock.Path).Return(stock, nil)
	serviceRepository.On("GetService", errServiceKey).Return(core.Service{}, errService)
	serviceRepository.On("GetService", errCacheServiceKey).Return(core.Service{}, nil)

	cacheConfigRepository.On("GetConfig", catalog.Path).Return(cacheConfig, nil)
	cacheConfigRepository.On("GetConfig", stock.Path).Return(core.CacheConfig{}, core.ErrCacheNotFound)
	cacheConfigRepository.On("GetConfig", errCacheServiceKey).Return(core.CacheConfig{}, errCache)

	tests := []struct {
		name              string
		key               string
		wantServiceDetail core.ServiceDetail
		wantErr           error
	}{
		{
			name:              "Error getting service",
			key:               errServiceKey,
			wantServiceDetail: core.ServiceDetail{},
			wantErr:           errService,
		},
		{
			name:              "Error getting cache",
			key:               errCacheServiceKey,
			wantServiceDetail: core.ServiceDetail{},
			wantErr:           errCache,
		},
		{
			name: "Get service detail without cache",
			key:  stock.Path,
			wantServiceDetail: core.ServiceDetail{
				Service: stock,
				Cache:   core.DefaultCacheConfig,
			},
			wantErr: nil,
		},
		{
			name: "Get service detail",
			key:  catalog.Path,
			wantServiceDetail: core.ServiceDetail{
				Service: catalog,
				Cache:   cacheConfig,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := controller.GetServiceDetail(tt.key)

			assert.Equal(t, tt.wantServiceDetail, service)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestDeleteService(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

	successService := "service"
	errDeleteService := "err-deleting-service"
	notExistsService := "not-exists-service"

	errServiceNotExist := errors.New("Service does not exist")
	errDeletingService := errors.New("Error deleting service")

	serviceRepository.On("ExistService", successService).Return(nil)
	serviceRepository.On("ExistService", errDeleteService).Return(nil)
	serviceRepository.On("ExistService", notExistsService).Return(errServiceNotExist)

	serviceRepository.On("DeleteService", successService).Return(nil)
	serviceRepository.On("DeleteService", errDeleteService).Return(errDeletingService)

	cacheConfigRepository.On("DeleteConfig", mock.Anything).Return(nil)

	tests := []struct {
		name    string
		key     string
		wantErr error
	}{
		{
			name:    "Service does not exist",
			key:     notExistsService,
			wantErr: errServiceNotExist,
		},
		{
			name:    "Error deleting service",
			key:     errDeleteService,
			wantErr: errDeletingService,
		},
		{
			name:    "Delete service successfully",
			key:     successService,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := controller.DeleteService(tt.key)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUpdateServiceStatus(t *testing.T) {
	serviceRepository := new(modelMocks.ServiceRepositoryI)
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

	serviceRepository.On("UpdateServiceStatus", mock.Anything, mock.Anything)

	controller.UpdateServiceStatus("catalog", core.ServiceStatusHealthy)
	controller.UpdateServiceStatus("stock", core.ServiceStatusIdle)

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
	cacheConfigRepository := new(modelMocks.CacheConfigRepositoryI)
	controller := NewServiceController(serviceRepository, cacheConfigRepository)

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
