package controller

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

	"github.com/gotway/gotway/internal/core"
	"github.com/gotway/gotway/internal/mocks"
	"github.com/gotway/gotway/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListenResponses(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewCacheController(cacheRepo, serviceRepo, log.Log)

	body := ioutil.NopCloser(bytes.NewBufferString("{}"))
	url, _ := url.Parse("http://api.gotway.com/catalog/products?offset=0&limit=10")
	httpRes := &http.Response{
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    url,
		},
		Body: body,
	}
	res := response{
		serviceKey:   "catalog",
		httpResponse: httpRes,
		body:         &body,
	}
	errRes := response{
		serviceKey:   "foo",
		httpResponse: httpRes,
		body:         &body,
	}

	cacheRepo.On("StoreCache", mock.Anything, mock.Anything).Return(nil)
	serviceRepo.On("IsCacheableStatusCode", mock.Anything, mock.Anything).Return(true)
	serviceRepo.On("GetServiceCache", res.serviceKey).Return(core.CacheConfig{}, nil)
	errCacheConfig := errors.New("Error getting cache config")
	serviceRepo.On("GetServiceCache", errRes.serviceKey).Return(core.CacheConfig{}, errCacheConfig)

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
	controller := NewCacheController(cacheRepo, serviceRepo, log.Log)

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
	serviceRepo.On("GetServiceCache", mock.Anything).Return(core.CacheConfig{}, nil)

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
	controller := NewCacheController(cacheRepo, serviceRepo, log.Log)

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
	serviceRepo.On("GetServiceCache", mock.Anything).Return(core.CacheConfig{}, nil)

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
	controller := NewCacheController(cacheRepo, serviceRepo, log.Log)

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
	serviceRepo.On("GetServiceCache", mock.Anything).Return(core.CacheConfig{}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.ListenResponses(ctx)

	err := controller.HandleResponse("foo", res)
	assert.NotNil(t, err)
}

func TestCachePolicy(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	serviceRepo := new(mocks.ServiceRepo)
	controller := NewCacheController(cacheRepo, serviceRepo, log.Log)

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
