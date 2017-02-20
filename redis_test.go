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
	redispassword = ""
	MaxIdle = 1000
	MaxActive = 1000
	IdleTimeout = 30 * time.Minute
)

var newRedisCache = func(t *testing.T, defaultExpiration time.Duration) RedisCache {
	revel.Config = config.NewContext()

	c, err := net.Dial("tcp", redisTestServer)
	if err == nil {
		c.Write([]byte("flush_all\r\n"))
		c.Close()
		redisCache := NewRedisCache(redisTestServer,redispassword, MaxIdle, MaxActive, IdleTimeout, 24*time.Hour)
		redisCache.Flush()
		return redisCache
	}
	t.Errorf("couldn't connect to redis on %s", redisTestServer)
	t.FailNow()
	panic("")
}

//start test ...

func TestRedisCache_TypicalGetSet(t *testing.T) {
	typicalGetSet(t, newRedisCache)
}