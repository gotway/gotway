package cache

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gotway/gotway/internal/mocks"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsCacheable(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	cacheableReq, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/service/foo", nil)
	notCacheableReq, _ := http.NewRequest(http.MethodPost, "http://api.gotway.com/service/foo", nil)

	tests := []struct {
		name     string
		req      *http.Request
		service  model.Service
		wantBool bool
	}{
		{
			name:     "GET REST",
			req:      cacheableReq,
			service:  model.Service{Type: model.ServiceTypeREST},
			wantBool: true,
		},
		{
			name:     "POST REST",
			req:      notCacheableReq,
			service:  model.Service{Type: model.ServiceTypeREST},
			wantBool: false,
		},
		{
			name:     "POST GRPC",
			req:      notCacheableReq,
			service:  model.Service{Type: model.ServiceTypeGRPC},
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCacheable := controller.IsCacheableRequest(tt.req, tt.service)

			assert.Equal(t, isCacheable, tt.wantBool)
		})
	}
}

func TestGetCache(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	reqCacheError, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/service/foo", nil)
	cacheError := errors.New("Cache not found")
	cacheRepo.On("Get", "/service/foo", "service").Return(model.Cache{}, cacheError)

	reqSuccess, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/catalog/products", nil)
	cache := model.Cache{
		Path:       "/products",
		StatusCode: 200,
	}
	cacheRepo.On("Get", "/catalog/products", "catalog").Return(cache, nil)

	catalog := model.Service{Name: "catalog"}

	tests := []struct {
		name      string
		req       *http.Request
		service   model.Service
		wantCache model.Cache
		wantErr   error
	}{
		{
			name:      "Cache not found error",
			req:       reqCacheError,
			service:   model.Service{},
			wantCache: model.Cache{},
			wantErr:   cacheError,
		},
		{
			name:      "Get cache successfully",
			req:       reqSuccess,
			service:   catalog,
			wantCache: cache,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := controller.GetCache(tt.req, tt.service)

			assert.Equal(t, tt.wantCache, cache)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetCacheDetail(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	reqCacheError, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/service/foo", nil)
	cacheError := errors.New("Cache not found")
	cacheRepo.On("GetCacheDetail", "/foo", "service").Return(model.Cache{}, cacheError)

	reqSuccess, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/catalog/products", nil)
	reqPrefix, _ := http.NewRequest(
		http.MethodGet,
		"http://api.gotway.com/api/cache/catalog/products",
		nil,
	)
	cacheDetail := model.Cache{
		Path:       "/products",
		StatusCode: 200,
		TTL:        10,
		Tags:       []string{"foo"},
	}
	cacheRepo.On("GetCacheDetail", "/products", "catalog").Return(cacheDetail, nil)

	tests := []struct {
		name            string
		req             *http.Request
		service         model.Service
		wantCacheDetail model.Cache
		wantErr         error
	}{
		{
			name:            "Cache detail not found error",
			req:             reqCacheError,
			service:         model.Service{Name: "service"},
			wantCacheDetail: model.Cache{},
			wantErr:         cacheError,
		},
		{
			name:            "Get cache successfully",
			req:             reqSuccess,
			service:         model.Service{Name: "catalog"},
			wantCacheDetail: cacheDetail,
			wantErr:         nil,
		},
		{
			name:            "Get cache successfully",
			req:             reqPrefix,
			service:         model.Service{Name: "catalog"},
			wantCacheDetail: cacheDetail,
			wantErr:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheDetail, err := controller.GetCache(tt.req, tt.service)

			assert.Equal(t, tt.wantCacheDetail, cacheDetail)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestDeleteCacheByPath(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	paths := []model.CachePath{
		{
			ServicePath: "service",
			Path:        "foo",
		},
	}
	cacheRepo.On("DeleteCacheByPath", paths).Return(nil)

	err := controller.DeleteCacheByPath(paths)

	assert.Nil(t, err)
	cacheRepo.AssertExpectations(t)
}

func TestDeleteCacheByTags(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	tags := []string{"foo"}
	cacheRepo.On("DeleteCacheByTags", tags).Return(nil)

	err := controller.DeleteCacheByTags(tags)

	assert.Nil(t, err)
	cacheRepo.AssertExpectations(t)
}

func TestListenResponses(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	bodyBytes := []byte("{}")
	body := ioutil.NopCloser(bytes.NewBufferString("{}"))
	url, _ := url.Parse("http://api.gotway.com/catalog/products?offset=0&limit=10")
	httpRes := &http.Response{
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    url,
		},
		Body: body,
	}
	catalog := model.Service{Name: "catalog"}
	res := response{
		httpResponse: httpRes,
		bodyBytes:    bodyBytes,
		service:      catalog,
	}
	errRes := response{
		httpResponse: httpRes,
		bodyBytes:    bodyBytes,
		service:      catalog,
	}

	cacheRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	serviceRepo.On("IsCacheableStatusCode", mock.Anything, mock.Anything).Return(true)
	serviceRepo.On("GetService", res.httpResponse, res.service).Return(model.CacheConfig{}, nil)
	errCacheConfig := errors.New("Error getting cache config")
	serviceRepo.On("GetServiceCache", errRes.serviceKey).Return(model.CacheConfig{}, errCacheConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.ListenResponses(ctx)

	for _, r := range []response{res, errRes} {
		if err := controller.HandleResponse(r.serviceKey, r.httpResponse); err != nil {
			t.Errorf("Got unexpected error: %w", err)
		}
	}

	time.Sleep(1 * time.Second)
	cacheRepo.AssertNumberOfCalls(t, "StoreCache", 1)
	serviceRepo.AssertNumberOfCalls(t, "IsCacheableStatusCode", 2)
	serviceRepo.AssertNumberOfCalls(t, "GetServiceCache", 2)
}

func TestListenCacheControlResponses(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	body := ioutil.NopCloser(bytes.NewBufferString("{}"))
	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	header := http.Header{}
	header.Set("Cache-Control", "s-maxage=10")
	TTLRes := response{
		serviceKey: "foo",
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			Header: header,
			Body:   body,
		},
		body: &body,
	}
	noTTLRes := response{
		serviceKey: "foo",
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			Body: body,
		},
		body: &body,
	}
	zeroTTLHeader := http.Header{}
	zeroTTLHeader.Set("Cache-Control", "s-maxage=0")
	zeroTTLRes := response{
		serviceKey: "foo",
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			Header: zeroTTLHeader,
			Body:   body,
		},
		body: &body,
	}

	cacheRepo.On("StoreCache", mock.Anything, mock.Anything).Return(nil)
	serviceRepo.On("IsCacheableStatusCode", mock.Anything, mock.Anything).Return(true)
	serviceRepo.On("GetServiceCache", mock.Anything).Return(model.CacheConfig{}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.ListenResponses(ctx)

	for _, r := range []response{TTLRes, noTTLRes, zeroTTLRes} {
		if err := controller.HandleResponse(r.serviceKey, r.httpResponse); err != nil {
			t.Errorf("Got unexpected error: %w", err)
		}
	}

	time.Sleep(1 * time.Second)
	cacheRepo.AssertNumberOfCalls(t, "StoreCache", 2)
	serviceRepo.AssertNumberOfCalls(t, "IsCacheableStatusCode", 2)
	serviceRepo.AssertNumberOfCalls(t, "GetServiceCache", 2)
}

func TestListenCacheTagsResponses(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	body := ioutil.NopCloser(bytes.NewBufferString("{}"))
	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	header := http.Header{}
	header.Set("X-Cache-Tags", "products")
	tagsRes := response{
		serviceKey: "foo",
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			Header: header,
			Body:   body,
		},
		body: &body,
	}
	noTagsRes := response{
		serviceKey: "foo",
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			Body: body,
		},
		body: &body,
	}

	cacheRepo.On("StoreCache", mock.Anything, mock.Anything).Return(nil)
	serviceRepo.On("IsCacheableStatusCode", mock.Anything, mock.Anything).Return(true)
	serviceRepo.On("GetServiceCache", mock.Anything).Return(model.CacheConfig{}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.ListenResponses(ctx)

	for _, r := range []response{tagsRes, noTagsRes} {
		if err := controller.HandleResponse(r.serviceKey, r.httpResponse); err != nil {
			t.Errorf("Got unexpected error: %w", err)
		}
	}

	time.Sleep(1 * time.Second)
	cacheRepo.AssertNumberOfCalls(t, "StoreCache", 2)
	serviceRepo.AssertNumberOfCalls(t, "IsCacheableStatusCode", 2)
	serviceRepo.AssertNumberOfCalls(t, "GetServiceCache", 2)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Error reading")
}

func TestErrReadingBody(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	testRequest := httptest.NewRequest(http.MethodPost, "/foo", errReader(0))
	body := testRequest.Body
	defer body.Close()

	res := &http.Response{
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    url,
		},
		Body: body,
	}

	cacheRepo.On("StoreCache", mock.Anything, mock.Anything).Return(nil)
	serviceRepo.On("IsCacheableStatusCode", mock.Anything, mock.Anything).Return(true)
	serviceRepo.On("GetServiceCache", mock.Anything).Return(model.CacheConfig{}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.ListenResponses(ctx)

	err := controller.HandleResponse("foo", res)
	assert.NotNil(t, err)
}

func TestCachePolicy(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewController(cacheRepo, serviceRepo, log.Log)

	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	notCacheableHeader := http.Header{}
	notCacheableHeader.Set("Cache-Control", "s-maxage=0")
	cacheableHeader := http.Header{}
	cacheableHeader.Set("Cache-Control", "s-maxage=10")
	invalidCacheableHeader := http.Header{}
	invalidCacheableHeader.Set("Cache-Control", "s-maxage")

	cacheableService := "catalog"
	notCacheableService := "stock"

	serviceRepo.On("IsCacheableStatusCode", cacheableService, mock.Anything).Return(true)
	serviceRepo.On("IsCacheableStatusCode", notCacheableService, mock.Anything).Return(false)

	tests := []struct {
		name            string
		serviceKey      string
		response        *http.Response
		wantIsCacheable bool
	}{
		{
			name:       "Not cacheable by method",
			serviceKey: cacheableService,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    url,
				},
				Header: cacheableHeader,
			},
			wantIsCacheable: false,
		},
		{
			name:       "Not cacheable by headers",
			serviceKey: cacheableService,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				Header: notCacheableHeader,
			},
			wantIsCacheable: false,
		},
		{
			name:       "Not cacheable by config",
			serviceKey: notCacheableService,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				Header: cacheableHeader,
			},
			wantIsCacheable: false,
		},
		{
			name:       "Cacheable",
			serviceKey: cacheableService,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				Header: cacheableHeader,
			},
			wantIsCacheable: true,
		},
		{
			name:       "Error parsing cache header",
			serviceKey: cacheableService,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				Header: invalidCacheableHeader,
			},
			wantIsCacheable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCacheable := controller.isCacheableResponse(tt.response, tt.serviceKey)

			assert.Equal(t, tt.wantIsCacheable, isCacheable)
		})
	}

}
