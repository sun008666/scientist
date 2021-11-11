package redistool

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// newPool New redis pool.
func NewPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialDatabase(db), redis.DialPassword(password))
			if err != nil {
				return nil, err
			}
			// if password != "" {
			// 	if _, err := c.Do("AUTH", password); err != nil {
			// 		c.Close()
			// 		log.Errorf("occur error at newPool: %v\n", err)
			// 		return nil, err
			// 	}
			// }
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Once(conn redis.Conn, key string, seconds int) bool {
	rs, e := redis.String(conn.Do("set", key, "done", "ex", seconds, "nx"))
	if e == redis.ErrNil {
		return false
	}
	if rs == "OK" {
		return true
	}
	return false
}
