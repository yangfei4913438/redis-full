package redis_full

import (
	"github.com/revel/config"
	"github.com/revel/revel"
	"net"
	"testing"
	"time"
)

// These tests require redis server running on localhost:6379 (the default)
const (
	TestServer  = "10.0.0.253:6379"
	password    = ""
	database    = 0
	MaxIdle     = 1000
	MaxActive   = 1000
	IdleTimeout = 30 * time.Minute
)

var newRedisCache = func(t *testing.T, defaultExpiration time.Duration) RedisCache {
	revel.Config = config.NewContext()

	c, err := net.Dial("tcp", TestServer)
	if err == nil {
		c.Write([]byte("flush_all\r\n"))
		c.Close()
		redisCache := NewRedisCache(TestServer, password, database, MaxIdle, MaxActive, IdleTimeout, 24*time.Hour)
		redisCache.FlushDB()
		return redisCache
	}
	t.Errorf("couldn't connect to redis on %s", TestServer)
	t.FailNow()
	panic("")
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

// Please insert test method before the Test_END method!!!

//CHECK END!
func Test_END(t *testing.T) {
	CheckEND(t, newRedisCache)
}
