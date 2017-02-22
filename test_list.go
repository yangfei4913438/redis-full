package redis_full

import (
	"testing"
	"time"
)

func CheckLINDEX(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)

	if err := redisDB.LPUSH("list_test_index", 2*time.Hour, "1", "2", "3"); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}

	var value interface{}
	if err := redisDB.LINDEX("list_test_index", 1, &value); err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}

	if string(value.(string)) != "2" {
		t.Errorf("result is error! right result is 2, but get %v", value)
		return
	}
}
