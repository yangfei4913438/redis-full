package redis_full

import (
	"time"
)

func (c RedisCache) APPEND(key string, value string, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	if times > 0 {
		_, err := conn.Do("APPEND", key, value)
		if err != nil {
			return err
		}
		_, err = conn.Do("EXPIRE", key, times)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("APPEND", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c RedisCache) STRLEN(key string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	res, err := conn.Do("STRLEN", key)
	if err != nil {
		return -1, err
	}

	return int64(res.(int64)), nil
}
