package middlewares

import (
	"time"

	"zonst/logging"
	"zonst/qipai-sports/api/configapisrv/config"

	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

func Redis(cacheName string, tomlConfig *config.Config) gin.HandlerFunc {
	cacheConfig, ok := tomlConfig.RedisServerConf(cacheName)
	if !ok {
		panic(fmt.Sprintf("%v not set.", cacheName))
	}

	// 链接数据库
	pool := newPool(cacheConfig.Addr, cacheConfig.Password, cacheConfig.DB)

	return func(c *gin.Context) {
		c.Set(cacheName, pool)
		c.Next()
	}

}

// newPool New redis pool.
func newPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					logging.Errorf("occur error at newPool: %v\n", err)
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
