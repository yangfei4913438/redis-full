package redis_full

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"time"
)

func (c RedisCache) SetJSON(key string, value interface{}, expires time.Duration) error {
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

func (c RedisCache) MSetJSON(args map[string]interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	for k, v := range args {
		if err := c.SetJSON(k, v, expires); err != nil {
			return err
		}
	}
	return nil
}

func (c RedisCache) GetJSON(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	raw, err := conn.Do("GET", key)
	if err != nil {
		return err
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

func (c RedisCache) MGetJSON(keys ...string) (map[string]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	values := make(map[string]interface{}, len(keys))

	for _, v := range keys {
		var data interface{}
		if err := c.GetJSON(v, &data); err != nil {
			return nil, err
		} else {
			values[v] = data
		}
	}

	return values, nil
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

func (c RedisCache) MDEL(keys ...string) error {
	conn := c.pool.Get()
	defer conn.Close()

	for _, v := range keys {
		if err := c.Del(v); err != nil {
			return err
		}
	}

	return nil
}
