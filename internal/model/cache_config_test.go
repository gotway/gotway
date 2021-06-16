package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name        string
		config      CacheConfig
		wantIsEmpty bool
	}{
		{
			name:        "Is empty cache config",
			config:      CacheConfig{},
			wantIsEmpty: true,
		},
		{
			name: "Is empty cache config with values",
			config: CacheConfig{
				TTL:      0,
				Statuses: []int{},
				Tags:     []string{},
			},
			wantIsEmpty: true,
		},
		{
			name:        "Is empty default cache config",
			config:      DefaultCacheConfig,
			wantIsEmpty: true,
		},
		{
			name: "Is empty non empty cache config",
			config: CacheConfig{
				TTL:      1,
				Statuses: []int{200},
				Tags:     []string{"foo"},
			},
			wantIsEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEmpty := tt.config.IsEmpty()

			assert.Equal(t, isEmpty, tt.wantIsEmpty)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  CacheConfig
		wantErr error
	}{
		{
			name:    "Validate empty config",
			config:  CacheConfig{},
			wantErr: nil,
		},
		{
			name:    "Validate default config",
			config:  DefaultCacheConfig,
			wantErr: nil,
		},
		{
			name: "Validate config with all values",
			config: CacheConfig{
				TTL:      1,
				Statuses: []int{200},
				Tags:     []string{"foo"},
			},
			wantErr: nil,
		},
		{
			name: "Validate config with some values",
			config: CacheConfig{
				Statuses: []int{200},
				Tags:     []string{"foo"},
			},
			wantErr: ErrInvalidCacheConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestSerializeStatuses(t *testing.T) {
	tests := []struct {
		name         string
		statuses     CacheConfigStatuses
		wantStatuses []string
	}{
		{
			name:         "Empty statuses",
			statuses:     []int{},
			wantStatuses: []string{},
		},
		{
			name:         "Statuses",
			statuses:     []int{200, 400, 404},
			wantStatuses: []string{"200", "400", "404"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statuses := tt.statuses.Serialize()

			assert.Equal(t, statuses, tt.wantStatuses)
		})
	}
}
