package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"queueJob/pkg/common/config"
)

var (
	RedisClient redis.UniversalClient
)

const (
	maxRetry = 10 // number of retries
)

func GetRds() redis.UniversalClient {
	return RedisClient
}

// NewRedis Initialize redis connection.
func NewRedis() error {
	if len(config.Config.Redis.Address) == 0 {
		return errors.New("redis address is empty")
	}
	var rdb redis.UniversalClient
	if len(config.Config.Redis.Address) > 1 || config.Config.Redis.ClusterMode {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:      config.Config.Redis.Address,
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password, // no password set
			PoolSize:   50,
			MaxRetries: maxRetry,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:       config.Config.Redis.Address[0],
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password,
			DB:         0,   // use default DB
			PoolSize:   100, // connection pool size
			MaxRetries: maxRetry,
		})
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis ping %w", err)
	}

	RedisClient = rdb
	return err
}

func CheckErr(err error) error {
	if err != nil && errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}
