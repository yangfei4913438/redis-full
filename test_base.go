package redis_full

import (
	"testing"
	"time"
)

type RedisFactory func(*testing.T, time.Duration) RedisCache

func CheckGET(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)

	if err := redisDB.SetJSON("name", "LiLei", 2*time.Hour); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}

	var name string
	if err := redisDB.GetJSON("name", &name); err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}

	if name != "LiLei" {
		t.Errorf("Failed to get the value of the key: %v", name)
	}
}

func CheckMGET(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)

	data := map[string]interface{}{
		"name": "LiLei",
		"age":  "18",
	}

	if err := redisDB.MSetJSON(data, 2*time.Hour); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}

	res, err := redisDB.MGetJSON("name", "age")
	if err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}

	name := res["name"]
	age := res["age"]

	if name != "LiLei" || age != "18" {
		t.Errorf("Failed to get the value of the key: %v", res)
	}
}
