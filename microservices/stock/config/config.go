package config

import (
	"log"
	"os"
	"strconv"
)

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid non int value for env variable '%s': %s", key, value)
	}
	return valueInt
}

var (
	// Port indicates the Stock API service port. It uses default K8s service port env variable
	Port = getEnvInt("STOCK_SERVICE_PORT", 10000)
	// RedisURL indicates the URL of redis
	RedisURL = getEnv("REDIS_URL", "localhost:6379")
	// RedisURL indicates the database of redis
	RedisDatabase = getEnvInt("REDIS_DATABASE", 1)
	// RedisPrefix is the prefix to be added to the keys
	RedisPrefix = getEnv("REDIS_PREFIX", "stock::")
	// RedisTTLDefault indicates the default time to live of the keys in seconds. 5 minutes by default
	RedisTTLDefault = getEnvInt("REDIS_TTL_DEFAULT", 300)
	// RedisTTLMax indicates the maximum time to live of the keys in seconds. 1 day by default
	RedisTTLMax = getEnvInt("REDIS_TTL_MAX", 86400)
)
