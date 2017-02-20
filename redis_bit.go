package redis_full

import (
	"errors"
	"time"
)

func (c RedisCache) GETBIT(key string, offset int) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return -1, err
	}

	if ok {
		res, err := conn.Do("GETBIT", key, offset)
		if err != nil {
			return -1, err
		}
		return int64(res.(int64)), nil
	} else {
		return -1, ErrCacheMiss
	}

}

func (c RedisCache) SETBIT(key string, offset, value int, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	if !IntIsBody([]int{0, 1}, value) {
		return errors.New("Value is only allowed is 0 or 1.")
	}

	times := int(expires / time.Second)

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}

	if ok {
		if times > 0 {
			_, err := conn.Do("SETBIT", key, offset, value)
			if err != nil {
				return err
			}
			_, err = conn.Do("EXPIRE", key, times)
			if err != nil {
				return err
			}
		} else {
			_, err := conn.Do("SETBIT", key, offset, value)
			if err != nil {
				return err
			}
		}
	} else {
		return ErrCacheMiss
	}

	return nil
}
