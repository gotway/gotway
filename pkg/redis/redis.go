package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

type Cmdable interface {
	redis.Cmdable
	OptimisticLockTx(
		ctx context.Context,
		maxRetries int,
		txFn RedisTxFn,
		keys ...string,
	) error
	HSetTTL(
		ctx context.Context,
		key string,
		values map[string]interface{},
		TTL time.Duration,
	) error
	SAddTTL(
		ctx context.Context,
		key string,
		TTL time.Duration,
		members ...interface{},
	) error
	AllExist(
		ctx context.Context,
		keys ...string,
	) (allExist bool, notExistsIndex int, err error)
}

type RedisTxFn = func(*redis.Tx) error

var (
	ErrWatchMaxRetries = errors.New("Maximum number of watch retries reached")
)

func (r *Client) OptimisticLockTx(
	ctx context.Context,
	maxRetries int,
	txFn RedisTxFn,
	keys ...string,
) error {
	wrappedTxFn := func(tx *redis.Tx) error {
		defer tx.Unwatch(ctx, keys...)
		return txFn(tx)
	}
	for i := 0; i < maxRetries; i++ {
		err := r.Watch(ctx, wrappedTxFn, keys...)
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}
	return ErrWatchMaxRetries
}

func (r *Client) HSetTTL(
	ctx context.Context,
	key string,
	values map[string]interface{},
	TTL time.Duration,
) error {
	pipe := r.TxPipeline()

	pipe.HSet(ctx, key, values)
	pipe.Expire(ctx, key, TTL)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Client) SAddTTL(
	ctx context.Context,
	key string,
	TTL time.Duration,
	members ...interface{},
) error {
	pipe := r.TxPipeline()

	pipe.SAdd(ctx, key, members...)
	pipe.Expire(ctx, key, TTL)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Client) AllExist(
	ctx context.Context,
	keys ...string,
) (allExist bool, notExistsIndex int, err error) {
	pipe := r.TxPipeline()
	cmds := make([]*redis.IntCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.Exists(ctx, key)
	}

	if _, err = pipe.Exec(ctx); err != nil {
		return false, -1, err
	}

	for index, cmd := range cmds {
		if cmd.Val() == 0 {
			return false, index, nil
		}
	}

	return true, -1, nil
}

func AnyEmptyErr(errs ...error) bool {
	for _, err := range errs {
		if err == redis.Nil {
			return true
		}
	}
	return false
}

func New(client *redis.Client) Cmdable {
	return &Client{client}
}
