# go-redis


```go
package main

import (
	"context"
	"fmt"
	//"github.com/jishulangcom/go-config"
	"github.com/jishulangcom/go-redis"
)

func main() {
	// redisCnf := &config.RedisCnfDto{} // 这里可以填写自己的配置
	// redisPoolCnf := &config.RedisPoolCnfDto{} // 这里可以填写自己的配置
	redis.NewDB(nil, nil) // nil时用默认值
	defer redis.CloseDB()

	//
	redis.DB.Set(context.Background(), "site", "技术狼|jishulang.com", 0)
	val := redis.DB.Get(context.Background(), "site").Val()
	fmt.Println(val)
}
```

