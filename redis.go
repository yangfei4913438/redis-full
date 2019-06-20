package redis_full

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

// Wraps the Redis client to meet the Cache interface.
type RedisCache struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

var ErrCacheMiss = errors.New("redis_full: key not found.")

// until redigo supports sharding/clustering, only one host will be in hostList
func NewRedisCache(host, password string, database, MaxIdle, MaxActive int, IdleTimeout, defaultExpiration time.Duration) RedisCache {
	var pool = &redis.Pool{
		MaxIdle:     MaxIdle,
		MaxActive:   MaxActive,
		IdleTimeout: IdleTimeout,
		Dial: func() (redis.Conn, error) {
			protocol := "tcp"
			c, err := redis.Dial(
				protocol,
				host,
				redis.DialDatabase(database),
				redis.DialPassword(password),
				redis.DialConnectTimeout(time.Second*5),
				redis.DialWriteTimeout(time.Second*3),
				redis.DialReadTimeout(time.Second*3),
			)
			if err != nil {
				return nil, err
			}
			return c, err
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

func (c RedisCache) FlushALL() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

func (c RedisCache) FlushDB() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHDB")
	return err
}

func (c RedisCache) Keys() ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	res, err := conn.Do("KEYS", "*")
	if err != nil {
		return nil, err
	}

	result, _ := res.([]interface{})

	var send []string

	for _, v := range result {
		s, _ := v.([]uint8)
		send = append(send, string(s))
	}

	return send, nil
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
