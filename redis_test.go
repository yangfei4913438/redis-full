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
	redisTestServer = "10.0.0.7:6379"
	redispassword   = ""
	MaxIdle         = 1000
	MaxActive       = 1000
	IdleTimeout     = 30 * time.Minute
)

var newRedisCache = func(t *testing.T, defaultExpiration time.Duration) RedisCache {
	revel.Config = config.NewContext()

	c, err := net.Dial("tcp", redisTestServer)
	if err == nil {
		c.Write([]byte("flush_all\r\n"))
		c.Close()
		redisCache := NewRedisCache(redisTestServer, redispassword, MaxIdle, MaxActive, IdleTimeout, 24*time.Hour)
		redisCache.Flush()
		return redisCache
	}
	t.Errorf("couldn't connect to redis on %s", redisTestServer)
	t.FailNow()
	panic("")
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
