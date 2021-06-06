package model

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock/v8"
	"github.com/gotway/gotway/log"
	logMocks "github.com/gotway/gotway/mocks/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockLogger() *logMocks.LoggerI {
	logger := new(logMocks.LoggerI)
	log.Logger = logger
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Fatalf", mock.Anything, mock.Anything).Return()
	return logger
}

func newRedisClientMock() (*redismock.ClientMock, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := newRedisClient(mr.Addr())
	return redismock.NewNiceMock(client), mr
}

func TestNewRedisClient(t *testing.T) {
	mockLogger()

	client, _ := newRedisClientMock()

	status, err := client.Ping(context.Background()).Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", status)
}

func TestInitFail(t *testing.T) {
	mockLogger()

	client := newRedisClient("foo")

	status, err := client.Ping(context.Background()).Result()
	assert.NotNil(t, err)
	assert.Equal(t, "", status)
}

func TestTTLcommands(t *testing.T) {
	mockLogger()

	client, miniredis := newRedisClientMock()
	TTL := time.Duration(1 * time.Second)
	assertTTL := func(t *testing.T, key string) {
		exists, err := client.Exists(ctx, key).Result()
		assert.Nil(t, err)
		assert.Equal(t, int64(1), exists)

		miniredis.FastForward(TTL)

		exists, err = client.Exists(ctx, key).Result()
		assert.Nil(t, err)
		assert.Equal(t, int64(0), exists)
	}

	key := "hash"
	values := map[string]interface{}{"foo": "bar"}
	err := hsetTTL(client, key, values, TTL)
	assert.Nil(t, err)
	assertTTL(t, key)

	key = "set"
	members := []interface{}{"foo"}
	err = saddTTL(client, key, TTL, members...)
	assert.Nil(t, err)
	assertTTL(t, key)
}

func TestAllExists(t *testing.T) {
	mockLogger()

	tests := []struct {
		name               string
		mockKeys           []string
		keys               []string
		wantAllExists      bool
		wantNotExistsIndex int
		wantErr            error
	}{
		{
			name:               "None exists",
			mockKeys:           []string{"a", "b", "c"},
			keys:               []string{"foo", "bar", "hello", "world"},
			wantAllExists:      false,
			wantNotExistsIndex: 0,
			wantErr:            nil,
		},
		{
			name:               "Not exists one",
			mockKeys:           []string{"foo", "hello", "world"},
			keys:               []string{"foo", "bar", "hello", "world"},
			wantAllExists:      false,
			wantNotExistsIndex: 1,
			wantErr:            nil,
		},
		{
			name:               "Many not exist",
			mockKeys:           []string{"foo", "hello", "world"},
			keys:               []string{"bar", "hello", "lol"},
			wantAllExists:      false,
			wantNotExistsIndex: 0,
			wantErr:            nil,
		},
		{
			name:               "All exist",
			mockKeys:           []string{"foo", "bar", "hello", "world"},
			keys:               []string{"foo", "bar", "hello", "world"},
			wantAllExists:      true,
			wantNotExistsIndex: -1,
			wantErr:            nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newRedisClientMock()
			for _, key := range tt.mockKeys {
				client.Set(context.Background(), key, key, 0)
			}

			allExist, notExistsIndex, err := allExists(client, tt.keys...)

			assert.Equal(t, tt.wantAllExists, allExist)
			assert.Equal(t, tt.wantNotExistsIndex, notExistsIndex)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestAllExistsError(t *testing.T) {
	mockLogger()

	client := newRedisClient("foo")

	allExist, notExistsIndex, err := allExists(client, "foo")

	assert.Equal(t, false, allExist)
	assert.Equal(t, -1, notExistsIndex)
	assert.NotNil(t, err)
}
