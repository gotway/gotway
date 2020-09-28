package log

import "testing"

func TestLoggerInit(t *testing.T) {
	Init()
	if Logger == nil {
		t.Error("Initialization failed")
	}
}
