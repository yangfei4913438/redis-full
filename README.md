# redis-full
all of the redis method

## How to install
Use `go get` to install or upgrade (`-u`) the `redis-full` package.

    go get -u github.com/yangfei4913438/redis-full

## Usage
Like on the command line using redis to use it!
For Example:

```golang

func (c App) SET() revel.Result {
	value1 := "hello"

	if err := app.RedisDB.Set("student", value1, 12*time.Hour); err != nil {
		data := map[string]interface{}{
			"status": false,
			"result": "Set the value of the key to redis failed!" + err.Error(),
		}
		return c.RenderJson(data)
	}

	data := map[string]interface{}{
		"status": true,
		"result": "Set the value of the key to redis success!",
	}
	return c.RenderJson(data)
}

func (c App) GET() revel.Result {
	var res string

	if err := app.RedisDB.Get("student", &res); err != nil {
		data := map[string]interface{}{
			"status": false,
			"result": "Failed to get the value of the key! " + err.Error(),
		}
		return c.RenderJson(data)
	}

	data := map[string]interface{}{
		"status": true,
		"result": res,
	}
	return c.RenderJson(data)

}
```

## Be careful!
    - The GET method and SET method is depend on each other!
    - Objects are stored in the redis before will first serialized using JSON.
    - Objects were taken out from the redis, before using, will first deserialization using JSON.

# Document
`Please be patient...`