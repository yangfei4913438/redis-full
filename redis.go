package redis_full

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

// Wraps the Redis client to meet the Cache interface.
type RedisCache struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

var ErrCacheMiss = errors.New("redis_full: key not found.")

// until redigo supports sharding/clustering, only one host will be in hostList
func NewRedisCache(host, password string, MaxIdle, MaxActive int, IdleTimeout, defaultExpiration time.Duration) RedisCache {
	var pool = &redis.Pool{
		MaxIdle:     MaxIdle,
		MaxActive:   MaxActive,
		IdleTimeout: IdleTimeout,
		Dial: func() (redis.Conn, error) {
			protocol := "tcp"
			c, err := redis.Dial(protocol, host)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}
	return RedisCache{pool, defaultExpiration}
}

func (c RedisCache) CheckRedis() error {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return err
	} else {
		return nil
	}
}

func (c RedisCache) Life(key string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return -1, err
	}
	if !ok {
		return -1, ErrCacheMiss
	}

	res, err := conn.Do("TTL", key)
	if err != nil {
		return -1, err
	} else {
		return int64(res.(int64)), nil
	}
}

func (c RedisCache) Type(key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrCacheMiss
	}

	res, err := conn.Do("TYPE", key)
	if err != nil {
		return "", err
	} else {
		return string(res.(string)), nil
	}
}

func (c RedisCache) Exists(key string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	res, err := conn.Do("EXISTS", key)
	if err != nil {
		return false, err
	} else {
		return int64(res.(int64)) == 1, nil
	}
}

func (c RedisCache) Flush() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

//some tools
func IsBody(s []interface{}, y interface{}) bool {
	m := make(map[interface{}]int)
	for _, v := range s {
		var x []interface{}
		for i := 0; i < len(s); i++ {
			if v == s[i] {
				x = append(x, v)
			}
		}
		m[v] = len(x)
	}
	return m[y] > 0
}
