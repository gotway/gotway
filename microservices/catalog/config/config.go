package config

import "os"

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var (
	// Port indicates the Catalog API service port. It uses default K8s service port env variable
	Port = getEnv("PORT", "9000")
)
