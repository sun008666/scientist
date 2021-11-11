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

// GanzhouGameConfigurationDeleteReq 赣州游戏配置删除列表
type GanzhouGameConfigurationDeleteReq struct {
	ID        int `json:"id" form:"id" binding:"required"`
	SubGameID int `json:"sub_game_id" form:"sub_game_id" binding:"required"`
}

// GanzhouGameConfigurationDeleteRequest 赣州游戏配置列表修改配置
func GanzhouGameConfigurationDeleteRequest(c *gin.Context) {
	req := &GanzhouGameConfigurationDeleteReq{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": err.Error()})
		return
	}
	defer LogStat(`GanzhouGameConfigurationUpdateReq`, c, req, time.Now())
	configdb := c.MustGet("configdb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	err := service.GanzhouGameConfiguratioDeleteRequest(configdb, pool, req.ID, req.SubGameID)
	if err != nil {
		logging.Errorf("GanzhouGameConfiguratioDeleteRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "C0300", "errmsg": "删除失败," + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errno":  "0",
		"errmsg": "删除成功",
	})
}
