package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/slavken/clean-architecture/pkg/logger"
	"github.com/slavken/clean-architecture/pkg/store"
)

type Store interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Del(keys ...string) error
	DelAll(pattern string) error
	store.Store
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

type rdb struct {
	cfg    *Config
	client *redis.Client
	log    logger.Logger
}

var ctx = context.Background()

func New(cfg *Config, log logger.Logger) Store {
	return &rdb{
		cfg: cfg,
		log: log,
	}
}

func (r *rdb) Open() error {
	rdb := redis.NewClient(&redis.Options{
		Addr: r.cfg.Addr,
		DB:   r.cfg.DB,
	})

	if err := rdb.Set(ctx, "key", "value", 0).Err(); err != nil {
		return err
	}

	r.client = rdb

	return nil
}

func (r *rdb) Close() error {
	return r.client.Close()
}

func (r *rdb) Get(key string) (string, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		r.log.Error("redis.Get: key does not exist")
		return "", err
	} else if err != nil {
		r.log.Errorf("redis.Get: %v", err)
		return "", err
	}

	return res, nil
}

func (r *rdb) Set(
	key string,
	value interface{},
	expiration time.Duration,
) error {
	if err := r.client.Set(
		ctx,
		key,
		value,
		expiration,
	).Err(); err != nil {
		r.log.Errorf("redis.Set: %v", err)
		return err
	}

	return nil
}

func (r *rdb) Del(keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		r.log.Errorf("redis.Del: %v", err)
		return err
	}

	return nil
}

func (r *rdb) DelAll(pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		go r.client.Del(ctx, iter.Val())
	}

	if err := iter.Err(); err != nil {
		r.log.Errorf("redis.DelAll: %v", err)
		return err
	}

	return nil
}
