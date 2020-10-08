package model

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gotway/gotway/config"
	"github.com/gotway/gotway/log"
)

const maxTxRetries = 1000

var (
	ctx             = context.Background()
	errTxMaxRetries = errors.New("Maximum number of transaction retries reached")
)

func newRedisClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Logger.Fatalf("Error connecting to redis server %s", err)
	}
	log.Logger.Infof("Connected to redis server %s", config.RedisServer)
	return client
}

func transaction(client *redis.Client, txFn func(*redis.Tx) error, keys ...string) error {
	for i := 0; i < maxTxRetries; i++ {
		err := client.Watch(ctx, txFn, keys...)
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}
	return errTxMaxRetries
}

func hsetTTL(client redis.Cmdable, key string, values map[string]interface{}, TTL time.Duration) error {
	pipe := client.TxPipeline()

	pipe.HSet(ctx, key, values)
	pipe.Expire(ctx, key, TTL)

	_, err := pipe.Exec(ctx)

	return err
}

func saddTTL(client redis.Cmdable, key string, TTL time.Duration, members ...interface{}) error {
	pipe := client.TxPipeline()

	pipe.SAdd(ctx, key, members...)
	pipe.Expire(ctx, key, TTL)

	_, err := pipe.Exec(ctx)

	return err
}

func allExists(client redis.Cmdable, keys ...string) (allExist bool, notExistsIndex int, err error) {
	pipe := client.TxPipeline()
	var cmds []*redis.IntCmd

	for _, key := range keys {
		cmd := pipe.Exists(ctx, key)
		cmds = append(cmds, cmd)
	}

	_, err = pipe.Exec(ctx)
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

func anyRedisNil(errs ...error) bool {
	for _, err := range errs {
		if err == redis.Nil {
			return true
		}
	}
	return false
}
