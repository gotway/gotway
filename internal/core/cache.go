package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Cache is a cached service response
type Cache struct {
	Path       string      `json:"path"`
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"headers"`
	Body       CacheBody   `json:"body"`
}

// CacheBody is a cache body
type CacheBody struct {
	io.Reader
}

// MarshalJSON serializes a cache body
func (c CacheBody) MarshalJSON() ([]byte, error) {
	bytes, err := ioutil.ReadAll(c)
	if err != nil || len(bytes) == 0 {
		var data struct{}
		return json.Marshal(data)
	}
	return bytes, nil
}

// CacheDetail provides extra information about a cache
type CacheDetail struct {
	Cache
	TTL  CacheTTL `json:"ttl"`
	Tags []string `json:"tags"`
}

// CacheTTL is the cache time to live in seconds
type CacheTTL time.Duration

// MarshalJSON serializes a TTL
func (c CacheTTL) MarshalJSON() ([]byte, error) {
	seconds := time.Duration(c) / time.Second
	return json.Marshal(seconds)
}

// NewCacheTTL creates a new cache TTL in seconds
func NewCacheTTL(seconds int64) CacheTTL {
	return CacheTTL(time.Duration(seconds) * time.Second)
}

// DeleteCache defines the params used to delete cache
type DeleteCache struct {
	Paths []CachePath `json:"paths"`
	Tags  []string    `json:"tags"`
}

// Validate checks if the payload is valid
func (p DeleteCache) Validate() error {
	if len(p.Paths) == 0 && len(p.Tags) == 0 {
		return ErrInvalidDeleteCache
	}
	if len(p.Paths) > 0 && len(p.Tags) > 0 {
		return ErrInvalidDeleteCache
	}
	return nil
}

// CachePath defines the paths that conform a cache
type CachePath struct {
	ServicePath string `json:"servicePath"`
	Path        string `json:"path"`
}

// ErrCachePathNotFound used when a cache defined by its path was not found
type ErrCachePathNotFound struct {
	CachePath
}

func (e *ErrCachePathNotFound) Error() string {
	return fmt.Sprintf("Cache path not found: %s%s", e.ServicePath, e.Path)
}

// ErrCacheNotFound error for not found cache
var ErrCacheNotFound = errors.New("Cache not found")

// ErrInvalidDeleteCache error for invalid delete cache objects
var ErrInvalidDeleteCache = errors.New("Paths or tags should be specified")
