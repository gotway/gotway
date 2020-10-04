package log

import (
	"testing"

	"github.com/gotway/gotway/config"
)

func TestInit(t *testing.T) {
	zap := initZap()
	if zap == nil {
		t.Error("Expected zap to be initialized")
	}
}

func TestInitProduction(t *testing.T) {
	config.Env = "production"
	zap := initZap()
	if zap == nil {
		t.Error("Expected zap to be initialized")
	}
}