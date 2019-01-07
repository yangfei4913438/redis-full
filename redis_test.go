package redis_full

import (
	"testing"
	"time"
)

// These tests require redis server running on localhost:6379 (the default)
const (
	TestServer  = "127.0.0.1:6379"
	password    = ""
	database    = 15
	MaxIdle     = 1000
	MaxActive   = 1000
	IdleTimeout = 30 * time.Minute
)

var newRedisCache = func(t *testing.T, defaultExpiration time.Duration) RedisCache {
	redisCache := NewRedisCache(TestServer, password, database, MaxIdle, MaxActive, IdleTimeout, 24*time.Hour)
	if err := redisCache.CheckRedis(); err != nil {
		t.Error("Redis Server:" + TestServer + " Connect failed: " + err.Error() + "!")
		t.FailNow()
		panic("")
	}
	return redisCache
}

//THE END METHOD
func CheckEND(t *testing.T, newredis RedisFactory) {
	redisDB := newredis(t, time.Hour)
	redisDB.FlushDB()
}

//TaskName must be start with Test_ prefix。

func Test_CheckGET(t *testing.T) {
	CheckGET(t, newRedisCache)
}

func Test_CheckMGET(t *testing.T) {
	CheckMGET(t, newRedisCache)
}

func Test_CheckGETBIT(t *testing.T) {
	CheckGETBIT(t, newRedisCache)
}

func Test_CheckGetBitSetBit(t *testing.T) {
	CheckGetBitSetBit(t, newRedisCache)
}

func Test_CheckBitOP(t *testing.T) {
	CheckBitOP(t, newRedisCache)
}

func Test_CheckLINDEX(t *testing.T) {
	CheckLINDEX(t, newRedisCache)
}

func Test_CheckSMEMBERS(t *testing.T) {
	CheckSMEMBERS(t, newRedisCache)
}

func Test_CheckSCARD(t *testing.T) {
	CheckSCARD(t, newRedisCache)
}

func Test_CheckLock(t *testing.T) {
	CheckLock(t, newRedisCache)
}

// Please insert test method before the Test_END method!!!

//CHECK END!
func Test_END(t *testing.T) {
	CheckEND(t, newRedisCache)
}
