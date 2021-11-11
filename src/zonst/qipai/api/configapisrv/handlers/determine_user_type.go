package handlers

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/service"
)

type DeterMineUerTypeReq struct {
	UserID int `json:"user_id" form:"user_id" binding:"required"`
}

// DeterMineUserType 查询用户类别
func DeterMineUserType(c *gin.Context) {
	req := &DeterMineUerTypeReq{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": err.Error()})
		return
	}
	defer LogStat(`DeterMineUserType`, c, req, time.Now())
	pointdb := c.MustGet("pointdb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	data, e := service.DeterMineUserType(pointdb, pool, req.UserID)
	if e != nil {
		logging.Errorf("DeterMineUserType: err: %v\n", e.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "C0300", "errmsg": "查询数据失败,"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":    "0",
		"errmsg":   "",
		"usertype": data,
	})
}
