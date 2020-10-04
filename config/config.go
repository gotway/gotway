package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gotway/gotway/cert"
)

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

func getIntEnv(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

var (
	// Port indicates the gotway API service port. It uses default K8s service port env variable
	Port = getEnv("GOTWAY_SERVICE_PORT", "8000")
	// Env indicates the environment name
	Env = getEnv("ENV", "development")
	// Database indicates which database is used for data storage
	Database = getEnv("DATABASE", "redis")
	// RedisServer indicates the URL for the redis client
	RedisServer = getEnv("REDIS_SERVER", "127.0.0.1:6379")
	// HealthCheckInterval is the interval between health checks
	HealthCheckInterval = time.Duration(getIntEnv("HEALTH_CHECK_INTERVAL_SECONDS", 10)) * time.Second
	// HealthCheckTimeout is the timeout for health check
	HealthCheckTimeout = time.Duration(getIntEnv("HEALTH_CHECK_TIMEOUT_SECONDS", 5)) * time.Second
	// TLS indicates if TLS is enabled
	TLS = getBoolEnv("TLS", true)
	// TLScert is the certificate file for TLS
	TLScert = getEnv("TLS_CERT", cert.Path("server.pem"))
	// TLSkey is the key file for TLS
	TLSkey = getEnv("TLS_KEY", cert.Path("server.key"))
)
