package config

import (
	"time"

	"github.com/gotway/gotway/pkg/env"
	"github.com/gotway/gotway/pkg/tlstest"
)

type Kubernetes struct {
	KubeConfig   string
	Namespace    string
	ResyncPeriod time.Duration
}

type HealthCheck struct {
	Enabled    bool
	NumWorkers int
	BufferSize int
	Interval   time.Duration
	Timeout    time.Duration
}

type Cache struct {
	Enabled    bool
	NumWorkers int
	BufferSize int
}

type TLS struct {
	Enabled bool
	Cert    string
	Key     string
}

type HA struct {
	Enabled            bool
	NodeId             string
	LeaseLockName      string
	LeaseLockNamespace string
	LeaseDuration      time.Duration
	RenewDeadline      time.Duration
	RetryPeriod        time.Duration
}

type Metrics struct {
	Enabled bool
	Path    string
	Port    string
}

type PProf struct {
	Enabled bool
	Port    string
}

type Config struct {
	Port           string
	Env            string
	LogLevel       string
	RedisUrl       string
	GatewayTimeout time.Duration

	Kubernetes  Kubernetes
	TLS         TLS
	HA          HA
	HealthCheck HealthCheck
	Cache       Cache
	Metrics     Metrics
	PProf       PProf
}

func GetConfig() (Config, error) {
	return Config{
		Port:           env.Get("PORT", "9111"),
		Env:            env.Get("ENV", "local"),
		LogLevel:       env.Get("LOG_LEVEL", "debug"),
		RedisUrl:       env.Get("REDIS_URL", "redis://localhost:6379/11"),
		GatewayTimeout: env.GetDuration("GATEWAY_TIMEOUT_SECONDS", 5) * time.Second,

		Kubernetes: Kubernetes{
			KubeConfig:   env.Get("KUBECONFIG", ""),
			Namespace:    env.Get("KUBERNETES_NAMESPACE", "default"),
			ResyncPeriod: env.GetDuration("KUBERNETES_RESYNC_PERIOD_SECONDS", 5) * time.Second,
		},
		TLS: TLS{
			Enabled: env.GetBool("TLS", true),
			Cert:    env.Get("TLS_CERT", tlstest.Cert()),
			Key:     env.Get("TLS_KEY", tlstest.Key()),
		},
		HA: HA{
			Enabled:            env.GetBool("HA", false),
			NodeId:             env.Get("HA_NODE_ID", ""),
			LeaseLockName:      env.Get("HA_LEASE_LOCK_NAME", "gotway"),
			LeaseLockNamespace: env.Get("HA_LEASE_LOCK_NAMESPACE", "gotway-system"),
			LeaseDuration:      env.GetDuration("HA_LEASE_DURATION_SECONDS", 15) * time.Second,
			RenewDeadline:      env.GetDuration("HA_RENEW_DEADLINE_SECONDS", 10) * time.Second,
			RetryPeriod:        env.GetDuration("HA_RETRY_PERIOD_SECONDS", 2) * time.Second,
		},
		HealthCheck: HealthCheck{
			Enabled:    env.GetBool("HEALTH_CHECK", true),
			NumWorkers: env.GetInt("HEALTH_CHECK_NUM_WORKERS", 10),
			BufferSize: env.GetInt("HEALTH_CHECK_BUFFER_SIZE", 10),
			Interval:   env.GetDuration("HEALTH_CHECK_INTERVAL_SECONDS", 10) * time.Second,
			Timeout:    env.GetDuration("HEALTH_CHECK_TIMEOUT_SECONDS", 5) * time.Second,
		},
		Cache: Cache{
			Enabled:    env.GetBool("CACHE", true),
			NumWorkers: env.GetInt("CACHE_NUM_WORKERS", 10),
			BufferSize: env.GetInt("CACHE_BUFFER_SIZE", 10),
		},
		Metrics: Metrics{
			Enabled: env.GetBool("METRICS", true),
			Path:    env.Get("METRICS_PATH", "/metrics"),
			Port:    env.Get("METRICS_PORT", "2112"),
		},
		PProf: PProf{
			Enabled: env.GetBool("PPROF", false),
			Port:    env.Get("PPROF_PORT", "6060"),
		},
	}, nil
}
