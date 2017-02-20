package redis_full

import (
	"time"
)

func (c RedisCache) SADD(key string, expires time.Duration, args ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	for _, v := range args {
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

func (c RedisCache) SISMEMBER(key string, value interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	res, err := conn.Do("SISMEMBER", key, value)
	if err != nil {
		return false, err
	}

	return int64(res.(int64)) == 1, nil
}

func (c RedisCache) SCARD(key string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	res, err := conn.Do("SCARD", key)
	if err != nil {
		return -1, err
	}

	return int64(res.(int64)), nil
}
