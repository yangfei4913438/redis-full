package redis_full

import (
	"testing"
	"time"
)

type redisFactory func(*testing.T, time.Duration) RedisCache

// Test typical cache interactions
func typicalGetSet(t *testing.T, newredis redisFactory) {
	var err error
	redisDB := newredis(t, time.Hour)

	value := "foo"
	if err = redisDB.Set("value", value, 2*time.Hour); err != nil {
		t.Errorf("Error setting a value: %s", err)
	}

	value = ""
	err = redisDB.Get("value", &value)
	if err != nil {
		t.Errorf("Error getting a value: %s", err)
	}
	if value != "foo" {
		t.Errorf("Expected to get foo back, got %s", value)
	}
}
