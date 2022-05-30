package redispool

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"github.com/jishulangcom/go-config"
	"net"
	"time"
)

var DB *redis.Client

func NewDB(redisCnf config.RedisCnfDto, redisPoolCnf *config.RedisPoolCnfDto) *redis.Client {
	if redisPoolCnf == nil {
		redisPoolCnf = &config.RedisPoolCnf
	}

	/*
		poolCnt := 4 * runtime.NumCPU()
		if poolCnt < 4 {
			poolCnt = 4
		}
	*/

	addr := fmt.Sprintf("%s:%d", redisCnf.Host, redisCnf.Port)
	DB = redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp",        //网络类型，tcp or unix，默认tcp
		Addr:     addr,         //主机名+冒号+端口，默认localhost:6379
		Password: redisCnf.Pwd, //密码
		DB:       redisCnf.DB,  // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     redisPoolCnf.PoolSize,     // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: redisPoolCnf.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  time.Duration(redisPoolCnf.DialTimeout) * time.Second,  //连接建立超时时间，默认5秒。
		ReadTimeout:  time.Duration(redisPoolCnf.ReadTimeout) * time.Second,  //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: time.Duration(redisPoolCnf.WriteTimeout) * time.Second, //写超时，默认等于读超时
		PoolTimeout:  time.Duration(redisPoolCnf.PoolTimeout) * time.Second,  //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: time.Duration(redisPoolCnf.IdleCheckFrequency) * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        time.Duration(redisPoolCnf.IdleTimeout) * time.Minute,        //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         time.Duration(redisPoolCnf.MaxConnAge) * time.Second,         //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      redisPoolCnf.MaxRetries,                                        // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: time.Duration(redisPoolCnf.MinRetryBackoff) * time.Millisecond, //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: time.Duration(redisPoolCnf.MaxRetryBackoff) * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
		//可自定义连接函数
		Dialer: func(con context.Context, network string, addr string) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   time.Duration(redisPoolCnf.Timeout) * time.Second,
				KeepAlive: time.Duration(redisPoolCnf.KeepAlive) * time.Minute,
			}
			return netDialer.Dial(network, addr)
		},

		//钩子函数
		OnConnect: func(ctx context.Context, conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
			//logger.Debug("conn=%v\n", conn)
			return nil
		},
	})

	cmd := DB.Ping(context.Background())
	err := cmd.Err()
	if err != nil {
		panic(err)
	}

	return DB
}

func CloseDB() {
	DB.Close()
}
