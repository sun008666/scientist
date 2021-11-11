package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
)

// LogStat 记录日志
func LogStat(logName string, c *gin.Context, reqBody interface{}, t1 time.Time) {
	r := c.Request
	log.Printf("%v: url:%v client_ip:%v body:%+v cost:%v\n",
		logName, r.RequestURI, c.ClientIP(), reqBody, time.Since(t1))
}

// LogStatUserID 记录日志内部工号
func LogStatUserID(logName string, c *gin.Context, userID int32, reqBody interface{}, t1 time.Time) {
	r := c.Request
	log.Printf("%v: url:%v client_ip:%v user_id:%v body:%+v cost:%v\n",
		logName, r.RequestURI, c.ClientIP(), userID, reqBody, time.Since(t1))
}
