package redis_full

import (
	"testing"
	"time"
)

type redisFactory func(*testing.T, time.Duration) RedisCache

//TEST EXIST METHOD
func CheckExists(t *testing.T, newredis redisFactory) {
	redisDB := newredis(t, time.Hour)
	if err := redisDB.Set("name", "tom", 2*time.Hour); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
	}
	ok, err := redisDB.Exists("name")
	if err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}
	if !ok {
		t.Error("The key is not exist!")
	}
}
