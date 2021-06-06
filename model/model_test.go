package model

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gotway/gotway/config"
	"github.com/gotway/gotway/log"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	log.Init()
	config.RedisServer = s.Addr()
	Init()

	assert.NotNil(t, ServiceRepository)
	assert.NotNil(t, CacheRepository)
}
