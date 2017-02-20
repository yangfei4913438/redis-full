package redis_full

import (
	"testing"
	"time"
)

type redisFactory func(*testing.T, time.Duration) RedisCache

//TEST EXIST METHOD
func CheckExists(t *testing.T, newredis redisFactory) {
	redisDB := newredis(t, time.Hour)
	if err := redisDB.Set("value", "1234", 2*time.Hour); err != nil {
		t.Errorf("Error setting a value: %s", err)
	}
	ok, err := redisDB.Exists("value")
	if err != nil {
		t.Errorf("found an accident %v", err)
		return
	}
	if !ok {
		t.Error("the key value is not exist!")
	}
}
