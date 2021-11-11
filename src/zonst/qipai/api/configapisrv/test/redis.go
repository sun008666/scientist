package test

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

func (cfg RedisConfig) NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (c redis.Conn, err error) {
			if cfg.Password == "" {
				c, err = redis.Dial("tcp", cfg.Addr, redis.DialDatabase(cfg.DB))
			} else {
				c, err = redis.Dial("tcp", cfg.Addr, redis.DialDatabase(cfg.DB),
					redis.DialPassword(cfg.Password))
			}
			return
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
