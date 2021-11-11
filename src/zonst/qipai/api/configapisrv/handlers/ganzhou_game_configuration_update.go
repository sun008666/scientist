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

// GanzhouGameConfigurationUpdateReq 赣州游戏配置修改列表
type GanzhouGameConfigurationUpdateReq struct {
	ID         int `json:"id" form:"id"`
	GameID     int `json:"game_id" form:"game_id" binding:"required"`
	GameAreaID int `json:"game_area_id" form:"game_area_id" `
	SubGameID  int `json:"sub_game_id" form:"sub_game_id" binding:"required"`
}

// GanzhouGameConfigurationUpdateRequest 赣州游戏配置列表修改配置
func GanzhouGameConfigurationUpdateRequest(c *gin.Context) {
	req := &GanzhouGameConfigurationUpdateReq{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": err.Error()})
		return
	}
	defer LogStat(`GanzhouGameConfigurationUpdateReq`, c, req, time.Now())
	configdb := c.MustGet("configdb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	b := service.GanzhouGameConfigurationFindRequest(configdb, req.GameID, req.GameAreaID, req.SubGameID)
	if b == true {
		c.JSON(http.StatusOK, gin.H{"errno": "200", "errmsg": "配置已存在"})
		return
	}
	err := service.GanzhouGameConfiguratioRequest(configdb, pool, req.GameID, req.GameAreaID, req.ID, req.SubGameID)
	if err != nil {
		logging.Errorf("GanzhouGameConfiguratioRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "C0300", "errmsg": "修改失败," + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errno":  "0",
		"errmsg": "更新成功",
	})
}
