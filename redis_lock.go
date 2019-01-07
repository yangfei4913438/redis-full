package redis_full

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func lockKey(val string) string {
	return fmt.Sprintf("redislock:%s", val)
}

func (c RedisCache) Lock(key string, expires time.Duration) (ok bool, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	// redis.String 对返回值进行格式化
	_, err = redis.String(conn.Do("SET", lockKey(key), "Redis Lock", "EX", int(expires.Seconds()), "NX"))

	// 判断返回值是否为空，空表示已经存在了
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}

	// 加锁出错，直接返回错误, 这里的错误是正常的错误，不是空值错误，空值错误在上面已经处理了。
	if err != nil {
		return false, err
	}

	// 没有任何问题，就表示加锁成功
	return true, nil
}

func (c RedisCache) Unlock(key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", lockKey(key))
	return err
}
