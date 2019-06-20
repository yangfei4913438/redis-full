package redis_full

import (
	"time"
)

// 查询剩余生存时间
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

// 设置 key 的剩余生存时间
func (c RedisCache) SetLife(key string, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCacheMiss
	}

	times := int(expires / time.Second)

	if times > 0 {
		_, err = conn.Do("EXPIRE", key, times)
		if err != nil {
			return err
		}
	}

	return nil
}
