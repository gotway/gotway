package config

import "os"

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var (
	// PORT indicates the Microgateway API service port. It uses default K8s service port env variable
	PORT = getEnv("MICROGATEWAY_SERVICE_PORT", "8000")
	// ENV indicates the environment name
	ENV = getEnv("ENV", "development")
)
