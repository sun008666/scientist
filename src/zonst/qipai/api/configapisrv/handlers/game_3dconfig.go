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

// On3DGameConfigListReq 游戏配置列表
type On3DGameConfigListReq struct {
	GameID int `json:"game_id" form:"game_id" binding:"required"`
}

// On3DGameConfigListRequest 3D大厅游戏展示位
func On3DGameConfigListRequest(c *gin.Context) {
	req := &On3DGameConfigListReq{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": "参数绑定错误"})
		return
	}
	defer LogStat(`On3DGameConfigListRequest`, c, req, time.Now())
	configdb := c.MustGet("configdb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	cacheID, err := service.On3DConfigListConfigID(configdb, pool, req.GameID)
	if err != nil {
		logging.Errorf("On3DConfigListConfigID: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "C0300", "errmsg": "查询数据失败,"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errno":  "0",
		"errmsg": "",
		"cache":  cacheID,
	})
}
