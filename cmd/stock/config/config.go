package config

import (
	c "github.com/gotway/gotway/pkg/config"
)

var (
	// Port indicates the Stock API service port
	Port = c.GetIntEnv("PORT", 13000)
	// RedisURL indicates the URL of redis
	RedisURL = c.GetEnv("REDIS_URL", "localhost:6379")
	// RedisURL indicates the database of redis
	RedisDatabase = c.GetIntEnv("REDIS_DATABASE", 1)
	// RedisPrefix is the prefix to be added to the keys
	RedisPrefix = c.GetEnv("REDIS_PREFIX", "stock::")
	// RedisTTLDefault indicates the default time to live of the keys in seconds. 5 minutes by default
	RedisTTLDefault = c.GetIntEnv("REDIS_TTL_DEFAULT", 300)
	// RedisTTLMax indicates the maximum time to live of the keys in seconds. 1 day by default
	RedisTTLMax = c.GetIntEnv("REDIS_TTL_MAX", 86400)
)
