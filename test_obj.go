package redis_full

import (
	"testing"
	"time"
)

type RedisFactory func(*testing.T, time.Duration) RedisCache

func CheckExists(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)
	if err := redisDB.Set("name", "tom", 2*time.Hour); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}
	ok, err := redisDB.Exists("name")
	if err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}
	if !ok {
		t.Error("The key is not exist!")
		return
	}
}

func CheckGETBIT(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)
	if err := redisDB.Set("name", "tom", 2*time.Hour); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}
	res, err := redisDB.GETBIT("name", 2)
	if err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}
	if res != 1 {
		t.Errorf("result is error! right result is 1, but get %v", res)
		return
	}
}

func CheckGetBitSetBit(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)
	if err := redisDB.Set("bit", "abc", 2*time.Hour); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}
	value := 1
	key := "bit"
	offset := 2
	if err := redisDB.SETBIT(key, offset, value, redisDB.defaultExpiration); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}
	res, err := redisDB.GETBIT(key, offset)
	if err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}
	if res != int64(value) {
		t.Errorf("GETBIT found an error result: %v", res)
		return
	}
}
