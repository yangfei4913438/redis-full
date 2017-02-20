package redis_full

import (
	"time"
)

func (c RedisCache) INCR(key string, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	if times > 0 {
		_, err := conn.Do("INCR", key)
		if err != nil {
			return err
		}
		_, err = conn.Do("EXPIRE", key, times)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("INCR", key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c RedisCache) INCRBY(key string, value int, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	if times > 0 {
		_, err := conn.Do("INCRBY", key, value)
		if err != nil {
			return err
		}
		_, err = conn.Do("EXPIRE", key, times)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("INCRBY", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c RedisCache) DECR(key string, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	if times > 0 {
		_, err := conn.Do("DECR", key)
		if err != nil {
			return err
		}
		_, err = conn.Do("EXPIRE", key, times)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("DECR", key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c RedisCache) DECRBY(key string, value int, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	if times > 0 {
		_, err := conn.Do("DECRBY", key, value)
		if err != nil {
			return err
		}
		_, err = conn.Do("EXPIRE", key, times)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("DECRBY", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c RedisCache) INCRBYFLOAT(key string, value float64, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	times := int(expires / time.Second)

	if times > 0 {
		_, err := conn.Do("INCRBYFLOAT", key, value)
		if err != nil {
			return err
		}
		_, err = conn.Do("EXPIRE", key, times)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("INCRBYFLOAT", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}
