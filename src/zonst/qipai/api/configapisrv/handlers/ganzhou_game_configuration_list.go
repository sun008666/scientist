package handlers

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/service"
)

// OnGanzhouGameConfigurationList 赣州游戏配置展示
func OnGanzhouGameConfigurationList(c *gin.Context) {
	configdb := c.MustGet("configdb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	data, err := service.GetGanzhouGameConfigurationList(configdb, pool)
	if err != nil {
		logging.Errorf("GetGanzhouGameConfigurationList: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "C0300", "errmsg": "查询数据失败,"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errno":  "0",
		"errmsg": "",
		"data":   data,
	})
}
