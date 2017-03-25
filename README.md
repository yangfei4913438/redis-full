# redis-full [![wercker status](https://app.wercker.com/status/5151a44054cbd71b158674b5b6093a6b/s/master "wercker status")](https://app.wercker.com/project/byKey/5151a44054cbd71b158674b5b6093a6b) [![Image of license](https://camo.githubusercontent.com/890acbdcb87868b382af9a4b1fac507b9659d9bf/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6c6963656e73652d4d49542d626c75652e737667)](https://github.com/yangfei4913438/redis-full/blob/master/LICENSE)         
all of the redis method

## DEV Version
- Redis version: 3.2.5
- Go version: 1.8

## How to install
Use `go get` to install or upgrade (`-u`) the `redis-full` package.

    go get -u github.com/yangfei4913438/redis-full

## Usage
Like on the command line using redis to use it! 

#### - use in revel

1）add a init file

```golang

package app

import (
	"xxxxx/app/models"
	"github.com/revel/revel"
	"time"
)

func InitRedis() {
	hosts, _ := revel.Config.String("cache.redis.hosts")
	password := revel.Config.StringDefault("cache.redis.password", "")
	database := revel.Config.IntDefault("cache.redis.db", 0)
	MaxIdle := revel.Config.IntDefault("cache.redis.maxidle", 100)
	MaxActive := revel.Config.IntDefault("cache.redis.maxactive", 1000)
	IdleTimeout := time.Duration(revel.Config.IntDefault("cache.redis.idletimeout", 600)) * time.Second

	if err := models.NewRD(hosts, password, database, MaxIdle, MaxActive, IdleTimeout); err != nil {
		revel.ERROR.Println("Redis Server:" + hosts + " Connect failed: " + err.Error() + "!")
	} else {
		revel.INFO.Println("Redis Server:" + hosts + " Connected!")
	}
}

```

2) regist to init.go 

```golang

func init(){
    revel.OnAppStart(InitRedis)
}

```

3) add redis model file

```golang
package models

import (
	redis "github.com/yangfei4913438/redis-full"
	"time"
)

var RedisDB redis.RedisCache

func NewRD(hosts, password string, database, MaxIdle, MaxActive int, IdleTimeout time.Duration) error {

	RedisDB = redis.NewRedisCache(hosts, password, database, MaxIdle, MaxActive, IdleTimeout, 24*time.Hour)

	return RedisDB.CheckRedis()
}


```


4) use it! so easy!

For Example, a model file:

```golang

func LoginOut(name, token string) (bool, error) {

	ok, err := CheckLogin(name, token)
	if err != nil {
		return false, err
	}

	if ok {
		if err := RedisDB.Del(strings.ToLower(name)); err != nil {
			return false, err
		} else {
			return true, nil
		}
	} else {
		return false, errors.New("非法操作，不是当前用户!")
	}

}
```

## Be careful!
    - The GETJSON method and SETJSON method is depend on each other!
    - Before the Objects are stored to redis, it will first serialized using JSON.
    - Objects were taken out from the redis, before using, it will first deserialization using JSON.

## More documentation, please be patient!
