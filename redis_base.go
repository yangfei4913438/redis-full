package redis_full

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"time"
)

func (c RedisCache) Set(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if times > 0 {
		_, err := conn.Do("SETEX", key, times, v)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("SET", key, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c RedisCache) Get(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("GET", key)
	if err != nil {
		return err
	} else if raw == nil {
		return ErrCacheMiss
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(item, &ptrValue); err != nil {
		return err
	}
	return nil
}

func (c RedisCache) Del(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	} else {
		return nil
	}
}
