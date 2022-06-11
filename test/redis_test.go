package test

import (
	"fmt"
	"github.com/jishulangcom/go-def"
	"github.com/jishulangcom/go-redis"
	"testing"
)

func Test(t *testing.T) {
	// redisCnf := &config.RedisCnfDto{} // 这里可以填写自己的配置
	// redisPoolCnf := &config.RedisPoolCnfDto{} // 这里可以填写自己的配置
	redis.NewDB(nil, nil) // nil时用默认值
	defer redis.CloseDB()

	//
	redis.DB.Set(def.Ctx, "site", "技术狼|jishulang.com", 0)
	val := redis.DB.Get(def.Ctx, "site").Val()
	fmt.Println(val)
}
