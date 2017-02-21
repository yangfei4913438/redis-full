package redis_full

import (
	"testing"
	"time"
)

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

func CheckBitOP(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)
	data := map[string]interface{}{
		"a": "bar",
		"b": "aar",
	}
	if err := redisDB.MSet(data, redisDB.defaultExpiration); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}
	res, err := redisDB.BITOP("OR", "a", "b")
	if err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}
	if res != "car" {
		t.Errorf("result is error! right result is car, but get %v", res)
	}
}
