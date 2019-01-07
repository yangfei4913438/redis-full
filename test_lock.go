package redis_full

import (
	"testing"
	"time"
)

// 相同两次设置来判断是否加锁成功
func CheckLock(t *testing.T, newRedis RedisFactory) {
	redisDB := newRedis(t, time.Hour)

	// 第一次加锁，应该成功
	ok, err := redisDB.Lock("name", time.Minute*10)
	if !ok || err != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}

	// 第二次加锁，应该是失败
	ok2, err2 := redisDB.Lock("name", time.Minute*10)
	if err2 != nil {
		t.Errorf("An unexpected error occurred: %s", err)
		return
	}

	if ok2 {
		t.Errorf("已经存在一个key, 正确的结果应该是false, 但是获取的值为 %v", ok2)
		return
	}

	// 解锁
	err3 := redisDB.Unlock("name")
	if err3 != nil {
		t.Errorf("An unexpected error occurred: %s", err3)
		return
	}

	// 再次加锁，应该成功
	ok3, err4 := redisDB.Lock("name", time.Minute*10)
	if !ok3 || err4 != nil {
		t.Errorf("An unexpected error occurred: %s", err4)
		return
	}
}
