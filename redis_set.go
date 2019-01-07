package redis_full

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"time"
)

func (c RedisCache) SADD(key string, expires time.Duration, args ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	for _, arg := range args {

		v, err := json.Marshal(arg)
		if err != nil {
			return err
		}

		if times > 0 {
			_, err := conn.Do("SADD", key, v)
			if err != nil {
				return err
			}
			_, err = conn.Do("EXPIRE", key, times)
			if err != nil {
				return err
			}
		} else {
			_, err := conn.Do("SADD", key, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c RedisCache) SREM(key string, args ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	for _, v := range args {
		_, err := conn.Do("SREM", key, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c RedisCache) SMEMBERS(key string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrCacheMiss
	}

	items, err := redis.Values(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}

	var results []interface{}
	for _, v := range items {
		item, err := redis.Bytes(v, nil)
		if err != nil {
			return nil, err
		} else {
			var x interface{}
			if err := json.Unmarshal(item, &x); err != nil {
				return nil, err
			}
			results = append(results, x)
		}
	}

	return results, nil
}

func (c RedisCache) SISMEMBER(key string, value interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, ErrCacheMiss
	}

	v, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	res, err := conn.Do("SISMEMBER", key, v)
	if err != nil {
		return false, err
	}

	return int64(res.(int64)) == 1, nil
}

func (c RedisCache) SCARD(key string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return -1, err
	}
	if !ok {
		return -1, ErrCacheMiss
	}

	res, err := conn.Do("SCARD", key)
	if err != nil {
		return -1, err
	}

	return int64(res.(int64)), nil
}
