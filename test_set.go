package redis_full

import (
	"testing"
	"time"
)

func CheckSMEMBERS(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)

	if err := redisDB.SADD("list_test_set", 2*time.Hour, "1", "2", "3"); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}

	res, err := redisDB.SMEMBERS("list_test_set")
	if err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}

	for _, v := range res {
		ok, err := redisDB.SISMEMBER("list_test_set", string(v.(string)))
		if err != nil {
			t.Errorf("An unexpected error occurred: %s", err)
			return
		}
		if !ok {
			t.Errorf("All the results should be true, but the results appear false! result: %v", string(v.(string)))
		}
	}

}

func CheckSCARD(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)

	if err := redisDB.SADD("list_test_set", 2*time.Hour, "1", "2", "3"); err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}

	res, err := redisDB.SCARD("list_test_set")
	if err != nil {
		t.Errorf("An unexpected error occurred: %v", err)
		return
	}

	if res != int64(3) {
		t.Errorf("The result should be 3, but the result is not 3! result: %v", res)
	}
}
