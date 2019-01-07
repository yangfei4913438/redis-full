package redis_full

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"strings"
	"time"
)

func (c RedisCache) LPUSH(key string, expires time.Duration, values ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	for _, value := range values {

		v, err := json.Marshal(value)
		if err != nil {
			return err
		}

		if times > 0 {
			_, err := conn.Do("LPUSH", key, v)
			if err != nil {
				return err
			}
			_, err = conn.Do("EXPIRE", key, times)
			if err != nil {
				return err
			}

		} else {
			_, err := conn.Do("LPUSH", key, v)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (c RedisCache) RPUSH(key string, expires time.Duration, values ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	for _, value := range values {

		v, err := json.Marshal(value)
		if err != nil {
			return err
		}

		if times > 0 {
			_, err := conn.Do("RPUSH", key, v)
			if err != nil {
				return err
			}
			_, err = conn.Do("EXPIRE", key, times)
			if err != nil {
				return err
			}

		} else {
			_, err := conn.Do("RPUSH", key, v)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (c RedisCache) LSET(key string, index int, value interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = conn.Do("LSET", key, index, v)
	if err != nil {
		return err
	}

	return nil
}

func (c RedisCache) LTRIM(key string, start, stop int) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	_, err = conn.Do("LTRIM", key, start, stop)
	if err != nil {
		return err
	}

	return nil
}

func (c RedisCache) LINSERT(key, direction string, pivot, value interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	if !IsBody([]interface{}{"AFTER", "BEFORE"}, strings.ToUpper(direction)) {
		return errors.New("Value only allowed after and before.")
	}

	v1, err := json.Marshal(pivot)
	if err != nil {
		return err
	}

	v2, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = conn.Do("LINSERT", key, strings.ToUpper(direction), v1, v2)
	if err != nil {
		return err
	}

	return nil
}

func (c RedisCache) LPOP(key string, result interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	raw, err := conn.Do("LPOP", key)
	if err != nil {
		return err
	}

	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(item, &result); err != nil {
		return err
	}

	return nil
}

func (c RedisCache) RPOP(key string, result interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	raw, err := conn.Do("RPOP", key)
	if err != nil {
		return err
	}

	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(item, &result); err != nil {
		return err
	}

	return nil
}

func (c RedisCache) LINDEX(key string, index int, result interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	raw, err := conn.Do("LINDEX", key, index)
	if err != nil {
		return err
	}

	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(item, &result); err != nil {
		return err
	}

	return nil
}

func (c RedisCache) LLEN(key string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return -1, err
	}
	if !ok {
		return -1, ErrCacheMiss
	}

	res, err := conn.Do("LLEN", key)
	if err != nil {
		return -1, err
	}

	return int64(res.(int64)), nil
}

//delete from left
func (c RedisCache) LLREM(key string, count, value int) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	_, err = conn.Do("LREM", key, 0-count, value)
	if err != nil {
		return err
	}

	return nil
}

//delete from right
func (c RedisCache) LRREM(key string, count, value int) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	_, err = conn.Do("LREM", key, count, value)
	if err != nil {
		return err
	}

	return nil
}

func (c RedisCache) LRANGE(key string, start, stop int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrCacheMiss
	}

	items, err := redis.Values(conn.Do("LRANGE", key, start, stop))
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

func (c RedisCache) RPOPLPUSH(key1, key2 string) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok1, err := c.Exists(key1)
	if err != nil {
		return err
	}
	if !ok1 {
		return ErrCacheMiss
	}

	ok2, err := c.Exists(key2)
	if err != nil {
		return err
	}
	if !ok2 {
		return ErrCacheMiss
	}

	_, err = conn.Do("RPOPLPUSH", key1, key2)
	if err != nil {
		return err
	}

	return nil
}
