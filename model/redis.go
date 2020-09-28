package model

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
)

var client *redis.Client

func initRedisClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.RedisServer,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Logger.Error("Error connecting to redis server")
		panic(err)
	}

	log.Logger.Info("Connected to redis server ", config.RedisServer)
}

func hsetTTL(key string, values map[string]interface{}, TTL time.Duration) error {
	err := client.HSet(context.Background(), key, values).Err()
	if err != nil {
		return err
	}
	err = client.Expire(context.Background(), key, TTL).Err()
	if err != nil {
		return err
	}
	return nil
}

func saddTTL(key string, TTL time.Duration, members ...interface{}) error {
	err := client.SAdd(context.Background(), key, members...).Err()
	if err != nil {
		return err
	}
	err = client.Expire(context.Background(), key, TTL).Err()
	if err != nil {
		return err
	}
	return nil
}

func indexedExists(keys ...string) (allExist bool, notExistsIndex int, err error) {
	pipe := client.Pipeline()
	var cmds []*redis.IntCmd

	for _, key := range keys {
		cmd := pipe.Exists(context.Background(), key)
		cmds = append(cmds, cmd)
	}

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return false, -1, err
	}

	for index, cmd := range cmds {
		if cmd.Val() == 0 {
			return false, index, nil
		}
	}

	return true, -1, nil
}
