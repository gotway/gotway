package config

import (
	"time"

	"github.com/gotway/gotway/pkg/env"
	"github.com/gotway/gotway/pkg/tlstest"
)

var (
	// Port indicates the gotway API service port. It uses default K8s service port env variable
	Port = env.Get("PORT", "11000")
	// Env indicates the environment name
	Env = env.Get("ENV", "local")
	// LogLevel indicates the log level
	LogLevel = env.Get("LOG_LEVEL", "debug")
	// RedisUrl indicates the URL for the redis client
	RedisUrl = env.Get("REDIS_URL", "redis://localhost:6379/11")
	// GatewayTimeout is the timeout when requesting services
	GatewayTimeout = env.GetDuration("GATEWAY_TIMEOUT_SECONDS", 5) * time.Second

	// HealthNumWorkers is the number of workers used to perform health check
	HealthNumWorkers = env.GetInt("HEALTH_CHECK_NUM_WORKERS", 10)
	// HealthBufferSize is the size of the buffered channel used to perform health check
	HealthBufferSize = env.GetInt("HEALTH_CHECK_BUFFER_SIZE", 10)

	// CacheNumWorkers is the number of workers used to perform health check
	CacheNumWorkers = env.GetInt("CACHE_NUM_WORKERS", 10)
	// CacheBufferSize is the size of the buffered channel used to perform health check
	CacheBufferSize = env.GetInt("CACHE_BUFFER_SIZE", 10)

	// HealthCheckEnabled determines if health check is enabled
	HealthCheckEnabled = env.GetBool("HEALTH_CHECK_ENABLED", false)
	// HealthCheckInterval is the interval between health checks
	HealthCheckInterval = env.GetDuration("HEALTH_CHECK_INTERVAL_SECONDS", 10) * time.Second
	// HealthCheckTimeout is the timeout for health check
	HealthCheckTimeout = env.GetDuration("HEALTH_CHECK_TIMEOUT_SECONDS", 5) * time.Second

	// TLS indicates if TLS is enabled
	TLS = env.GetBool("TLS", true)
	// TLScert is the certificate file for TLS
	TLScert = env.Get("TLS_CERT", tlstest.Cert())
	// TLSkey is the key file for TLS
	TLSkey = env.Get("TLS_KEY", tlstest.Key())

	// Metrics indicates whether the metrics are enabled
	Metrics = env.GetBool("METRICS", true)
	// MetricsPath indices the metrics server path
	MetricsPath = env.Get("METRICS_PATH", "/metrics")
	// MetricsPort indicates the metrics server port
	MetricsPort = env.Get("METRICS_PORT", "2112")

	// PProf indicates whether profiling is enabled
	PProf = env.GetBool("PPROF", false)
	// PProfPort indicates the port of the profiling server
	PProfPort = env.Get("PPROF_PORT", "6060")
)
