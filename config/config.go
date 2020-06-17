package config

import "os"

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var (
	// Port indicates the Microgateway API service port. It uses default K8s service port env variable
	Port = getEnv("MICROGATEWAY_SERVICE_PORT", "8080")
	// Env indicates the environment name
	Env = getEnv("ENV", "development")
	// Database indicates which database is used for data storage
	Database = getEnv("DATABASE", "redis")
	// RedisServer indicates the URL for the redis client
	RedisServer = getEnv("REDIS_SERVER", "127.0.0.1:6379")
	// ServiceCheckInterval seconds between service health checks
	ServiceCheckInterval = getEnv("SERVICE_CHECK_INTERVAL", "10")
)
