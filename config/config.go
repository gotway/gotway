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
	PORT = getEnv("MICROGATEWAY_SERVICE_PORT", "8080")
	// ENV indicates the environment name
	ENV = getEnv("ENV", "development")
	// REDIS_SERVER indicates the URL for the redis client
	REDIS_SERVER = getEnv("REDIS_SERVER", "127.0.0.1:6379")
	// DEFAULT_SERVICE_TTL standard ttl for a service that does not specify one
	DEFAULT_SERVICE_TTL = getEnv("DEFAULT_SERVICE_TTL", "300")
	// MAX_SERVICE_TTL maximum ttl possible for a service
	MAX_SERVICE_TTL = getEnv("MAX_SERVICE_TTL", "1800")
)
