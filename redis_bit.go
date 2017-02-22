package redis_full

import (
	"errors"
	"strings"
	"time"
)

func (c RedisCache) BITCOUNT(key string, start, end int) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	ok, err := c.Exists(key)
	if err != nil {
		return -1, err
	}
	if !ok {
		return -1, ErrCacheMiss
	}

	res, err := conn.Do("BITCOUNT", key, start, end)
	if err != nil {
		return -1, err
	}

	if int64(res.(int64)) == 0 {
		return -1, ErrCacheMiss
	} else {
		return int64(res.(int64)), nil
	}
}

func (c RedisCache) BITOP(opt string, key1, key2 string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	if !IsBody([]interface{}{"AND", "OR", "XOR", "NOT"}, strings.ToUpper(opt)) {
		return "", errors.New("Value only allowed AND,OR,XOR,NOT.")
	}

	ok1, err := c.Exists(key1)
	if err != nil {
		return "", err
	}
	if !ok1 {
		return "", errors.New("Not found key: " + key1)
	}

	ok2, err := c.Exists(key2)
	if err != nil {
		return "", err
	}
	if !ok2 {
		return "", errors.New("Not found key: " + key2)
	}

	value := "bItOp_rEcIve_dATa"
	ok3, err := c.Exists(value)
	if err != nil {
		return "", err
	}
	if ok3 {
		value = value + "_suan_ni_xiao_zi_niu_bi_wo_fu_le_from_yangfei"
	}

	_, err = conn.Do("BITOP", strings.ToUpper(opt), value, key1, key2)
	if err != nil {
		return "", err
	}

	var str string
	if err := c.GetJSON(value, &str); err != nil {
		return "", err
	} else {
		//This is not a must delete key.
		c.Del(value)
		return str, nil
	}
}

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

	if !IsBody([]interface{}{0, 1}, value) {
		return errors.New("Value only allowed 0 or 1.")
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
