package model

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBodyMarshal(t *testing.T) {
	tests := []struct {
		name       string
		body       CacheBody
		wantString string
	}{
		{
			name: "Marshal an empty body",
			body: CacheBody{
				Reader: ioutil.NopCloser(bytes.NewBufferString("")),
			},
			wantString: "{}",
		},
		{
			name: "Marshal a body with content",
			body: CacheBody{
				Reader: ioutil.NopCloser(bytes.NewBufferString("{\"foo\":\"bar\"}")),
			},
			wantString: "{\"foo\":\"bar\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := tt.body.MarshalJSON()

			if err != nil {
				t.Errorf("Got unexpected error: %w", err)
			}
			assert.Equal(t, string(bytes), tt.wantString)
		})
	}
}

func TestTTLMarshal(t *testing.T) {
	tests := []struct {
		name       string
		ttl        CacheTTL
		wantString string
	}{
		{
			name:       "Marshal a TTL of zero seconds",
			ttl:        NewCacheTTL(0),
			wantString: "0",
		},
		{
			name:       "Marshal a TTL in seconds",
			ttl:        NewCacheTTL(300),
			wantString: "300",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := tt.ttl.MarshalJSON()

			if err != nil {
				t.Errorf("Got unexpected error: %w", err)
			}
			assert.Equal(t, string(bytes), tt.wantString)
		})
	}
}

func TestDeleteValidate(t *testing.T) {
	tests := []struct {
		name    string
		delete  DeleteCache
		wantErr error
	}{
		{
			name: "Validate empty delete",
			delete: DeleteCache{
				Paths: []CachePath{},
				Tags:  []string{},
			},
			wantErr: ErrInvalidDeleteCache,
		},
		{
			name: "Validate delete with paths and tags",
			delete: DeleteCache{
				Paths: []CachePath{
					{
						ServicePath: "catalog",
						Path:        "/products",
					},
				},
				Tags: []string{"catalog"},
			},
			wantErr: ErrInvalidDeleteCache,
		},
		{
			name: "Validate valid delete",
			delete: DeleteCache{
				Tags: []string{"catalog"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.delete.Validate()

			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestErrCachePathFormat(t *testing.T) {
	err := &ErrCachePathNotFound{
		CachePath: CachePath{
			ServicePath: "catalog",
			Path:        "/products",
		},
	}

	assert.EqualError(t, err, "Cache path not found: catalog/products")
}
