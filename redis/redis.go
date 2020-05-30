package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
)

var client *redis.Client

// Init Initialize a redis client
func Init() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_SERVER,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error("Error connecting to redis server")
		panic(err)
	}

	log.Info("Connected to redis server ", config.REDIS_SERVER)
}

// Store a key and a value in Redis
func Store(key string, url string, healthEndpoint string, ttl int) bool {

	redisKey := "service:" + key

	saved, err := client.HSet(context.Background(), redisKey, "url", url, "healthEndpoint", healthEndpoint, "ttl", ttl).Result()
	if err != nil {
		return false
	}

	// TODO do we have to check if ttl was correctly set?
	client.Expire(context.Background(), redisKey, time.Duration(ttl)*time.Second)

	return saved != 0
}

// DoesKeyExists checks if a certain key exists in Redis
func DoesKeyExists(key string, field string) bool {
	found, err := client.HGet(context.Background(), key, field).Result()
	if err != nil {
		return false
	}

	return found != ""
}

// Crear servicio: HMSET service:<endpoint> health "<health_endpoint>" type "<type>"
// Actualizar campo del servicio: HMSET service:<endpoint> health "<health_endpoint>"
// Obtener campo del servicio: HMGET service:<endpoint> health
// Comprobar si existe un servicio: HLEN service:<endpoint> (si devuelve > 0 ya existe el servicio)
// Borrar servicio: DEL service:<endpoint>
