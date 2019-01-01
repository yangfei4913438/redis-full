# redis-full [![wercker status](https://app.wercker.com/status/5151a44054cbd71b158674b5b6093a6b/s/master "wercker status")](https://app.wercker.com/project/byKey/5151a44054cbd71b158674b5b6093a6b) [![Image of license](https://camo.githubusercontent.com/890acbdcb87868b382af9a4b1fac507b9659d9bf/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6c6963656e73652d4d49542d626c75652e737667)](https://github.com/yangfei4913438/redis-full/blob/master/LICENSE)         
all of the redis method

## DEV Version
- Redis version: 3.2.12
- Go version: 1.11

## How to install
Use `go get` to install or upgrade (`-u`) the `redis-full` package.

    go get -u github.com/yangfei4913438/redis-full

## Usage
Like on the command line using redis to use it! 

#### - use in beego

1）add a init file

```golang

package dbs

import (
	"github.com/astaxie/beego"
	redis "github.com/yangfei4913438/redis-full"
	"strings"
	"time"
)

//redis对外接口
var RedisDB redis.RedisCache

func initRedis() {
	hosts := beego.AppConfig.String("redis.host")
	password := beego.AppConfig.DefaultString("redis.password", "")
	database := beego.AppConfig.DefaultInt("redis.db", 0)
	MaxIdle := beego.AppConfig.DefaultInt("redis.maxidle", 100)
	MaxActive := beego.AppConfig.DefaultInt("redis.maxactive", 1000)
	IdleTimeout := time.Duration(beego.AppConfig.DefaultInt("redis.idletimeout", 600)) * time.Second

	//通过赋值对外接口来使用
	RedisDB = redis.NewRedisCache(hosts, password, database, MaxIdle, MaxActive, IdleTimeout, 24*time.Hour)

	if err := RedisDB.CheckRedis(); err != nil {
		panic("Redis Server:" + hosts + " Connect failed: " + err.Error() + "!")
	} else {
		beego.Info("Connect Redis Server(" + hosts + ") to successful!")
	}
}

```

2) register to init.go 

```golang

package dbs

func init() {
	initMysql()
	initRedis()
}

```

4) use it! so easy!

For Example, a model file:

```golang

package models

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"testapi/dbs"
	"testapi/tools"
)

// 用户表结构体,用于接收数据库查询出来的对象，数据类型和数据库尽量保持一致
type User struct {
	Id          int64 `json:"id" db:"id"`
	ReceiveUser
}

// 添加用户时，接收用户传值的对象
type ReceiveUser struct {
	Name  string `json:"name" db:"name"`
	Age   int64  `json:"age" db:"age"`
	Email string `json:"email" db:"email"`
}

// 查询用户
func SelectUser(id int64) (resObj *User, resErr error) {
	// 定义redis的key, id转string类型
	redisKey := "test:user_" + strconv.FormatInt(id, 10)

	// 定义接收数据的对象
	var user User

	// 先从缓存查询，没有再从数据库查
	if err := dbs.RedisDB.GetJSON(redisKey, &user); err != nil {
		if strings.Contains(err.Error(), "key not found") {
			// key不存在，就重新查询一次

			// 预处理SQL语句
			selectSql := "select * from users where id=? limit 1"

			// 打印日志
			beego.Debug("[sql]: "+selectSql, id)
			err := dbs.MysqlDB.Get(&user, selectSql, id)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					beego.Trace("查询结果为空值!")

					// 将空值添加到缓存, 有效期1小时
					if err1 := dbs.RedisDB.SetJSON(redisKey, nil, tools.OneHour); err1 != nil {
						beego.Error(err1)
						return nil, err1
					}
					// 返回空值
					return nil, nil
				} else {
					// 打印错误日志
					beego.Error(err)
					// 返回错误信息
					return nil, err
				}
			}

			// 将结果添加到缓存
			if err2 := dbs.RedisDB.SetJSON(redisKey, &user, tools.OneDay); err2 != nil {
				// 打印错误日志
				beego.Error(err2)
				// 返回错误信息
				return nil, err2
			}

			// 返回结果给用户
			return &user, nil
		}
	}

	// 空值返回用户信息
	if user.Id == 0 {
		// 因为正常情况下，ID是从1开始的。0就表示读取出来的值是空值
		return nil, nil
	} else {
		return &user, nil
	}
}
```

## Be careful!
    - The GETJSON method and SETJSON method is depend on each other!
    - Before the Objects are stored to redis, it will first serialized using JSON.
    - Objects were taken out from the redis, before using, it will first deserialization using JSON.

## More documentation, please be patient!
