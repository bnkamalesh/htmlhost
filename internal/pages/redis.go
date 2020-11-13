package pages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// RedisDriver has all the dependencies required to initialize redigo & exposes all the methods
// provided by redigo
type RedisDriver struct {
	pool *redis.Pool
}

// Conn returns a new redis.Conn
func (rd *RedisDriver) Conn(ctx context.Context) (redis.Conn, error) {
	conn, err := rd.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// newRedisDriver returns a new RedisDriver with all the required dependencies initialized
func newRedisDriver(cfg *Config) (*RedisDriver, error) {
	db, _ := strconv.Atoi(cfg.StoreName)
	rpool := &redis.Pool{
		MaxIdle:         cfg.PoolSize,
		MaxActive:       cfg.PoolSize,
		IdleTimeout:     cfg.IdleTimeoutSecs,
		Wait:            true,
		MaxConnLifetime: cfg.IdleTimeoutSecs * 2,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
				redis.DialReadTimeout(cfg.ReadTimeoutSecs),
				redis.DialWriteTimeout(cfg.WriteTimeoutSecs),
				redis.DialPassword(cfg.Password),
				redis.DialConnectTimeout(cfg.DialTimeoutSecs),
				redis.DialDatabase(db),
			)
		},
	}

	conn := rpool.Get()
	rep, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}

	pong, _ := rep.(string)
	if pong != "PONG" {
		return nil, errors.New("ping failed")
	}
	defer conn.Close()

	return &RedisDriver{
		pool: rpool,
	}, nil
}

// CacheSerialize is used to serialize data which is to be stored in the cache
func CacheSerialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// CacheDeserialize is used to deserialize data retrieved from cache
func CacheDeserialize(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
