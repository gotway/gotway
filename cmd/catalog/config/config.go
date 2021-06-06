package config

import c "github.com/gotway/gotway/pkg/config"

var (
	// Port indicates the Catalog API service port. It uses default K8s service port env variable
	Port = c.GetEnv("PORT", "12000")
)
