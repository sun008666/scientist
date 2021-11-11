package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func InitRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial("tcp", "123.206.176.76:6379")
			if e != nil {
				c.Close()
				return nil, e
			}
			_, e = c.Do("AUTH", "crs-obchkyho:qipai0918")
			if e != nil {
				c.Close()
				return nil, e
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, e := c.Do("whiteListCache")
			if e != nil {
				fmt.Printf("fail to ping redis: %v\n", e)
				return e
			}
			return nil
		},
	}
}

func TestRedis(t *testing.T) {
	conn := InitRedisPool().Get()
	defer conn.Close()

	_, e := conn.Do("KEYS *")
	fmt.Println(e.Error())
	assert.Nil(t, e)
}

func TestDeterMineUserType(t *testing.T) {
	pointdb, e := gorm.Open("postgres", "host=123.206.176.76 port=5432 user=qipai dbname=pointdb password=qipai#xq5 sslmode=disable")
	assert.Nil(t, e)
	defer pointdb.Close()

	//flag.Parse()
	pool := InitRedisPool()
	var userID int
	userID = 1234
	bindGameID, e := DeterMineUserType(pointdb, pool, userID)
	assert.Nil(t, e)
	fmt.Println(bindGameID)
}
