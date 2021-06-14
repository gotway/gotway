package config

import (
	c "github.com/gotway/gotway/pkg/config"
)

var (
	// Port indicates the server port
	Port = c.GetEnv("PORT", "14000")
)
