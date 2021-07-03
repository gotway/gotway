package config

import (
	"time"

	c "github.com/gotway/gotway/pkg/config"
	"github.com/gotway/gotway/pkg/tlstest"
)

var (
	// Port indicates the gotway API service port. It uses default K8s service port env variable
	Port = c.GetEnv("PORT", "11000")
	// Env indicates the environment name
	Env = c.GetEnv("ENV", "local")
	// LogLevel indicates the log level
	LogLevel = c.GetEnv("LOG_LEVEL", "debug")
	// RedisUrl indicates the URL for the redis client
	RedisUrl = c.GetEnv("REDIS_URL", "redis://localhost:6379/11")
	// GatewayTimeout is the timeout when requesting services
	GatewayTimeout = time.Duration(c.GetIntEnv("GATEWAY_TIMEOUT_SECONDS", 5)) * time.Second
	// HealthNumWorkers is the number of workers used to perform health check
	HealthNumWorkers = c.GetIntEnv("HEALTH_CHECK_NUM_WORKERS", 10)
	// HealthBufferSize is the size of the buffered channel used to perform health check
	HealthBufferSize = c.GetIntEnv("HEALTH_CHECK_BUFFER_SIZE", 10)
	// CacheNumWorkers is the number of workers used to perform health check
	CacheNumWorkers = c.GetIntEnv("CACHE_NUM_WORKERS", 10)
	// CacheBufferSize is the size of the buffered channel used to perform health check
	CacheBufferSize = c.GetIntEnv("CACHE_BUFFER_SIZE", 10)
	// HealthCheckInterval is the interval between health checks
	HealthCheckInterval = time.Duration(
		c.GetIntEnv("HEALTH_CHECK_INTERVAL_SECONDS", 10),
	) * time.Second
	// HealthCheckTimeout is the timeout for health check
	HealthCheckTimeout = time.Duration(c.GetIntEnv("HEALTH_CHECK_TIMEOUT_SECONDS", 5)) * time.Second
	// TLS indicates if TLS is enabled
	TLS = c.GetBoolEnv("TLS", true)
	// TLScert is the certificate file for TLS
	TLScert = c.GetEnv("TLS_CERT", tlstest.Cert())
	// TLSkey is the key file for TLS
	TLSkey = c.GetEnv("TLS_KEY", tlstest.Key())
	// Metrics indicates whether the metrics are enabled
	Metrics = c.GetBoolEnv("METRICS", true)
	// MetricsPath indices the metrics server path
	MetricsPath = c.GetEnv("METRICS_PATH", "/metrics")
	// MetricsPort indicates the metrics server port
	MetricsPort = c.GetEnv("METRICS_PORT", "2112")
	// PProf indicates whether profiling is enabled
	PProf = c.GetBoolEnv("PPROF", false)
	// PProfPort indicates the port of the profiling server
	PProfPort = c.GetEnv("PPROF_PORT", "6060")
)
