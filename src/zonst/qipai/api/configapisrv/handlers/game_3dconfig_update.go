package handlers

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"net/http"
	"time"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/service"
)

// On3DGameConfigUpdateReq 游戏配置列表
type On3DGameConfigUpdateReq struct {
	GameID     int           `json:"game_id" form:"game_id" binding:"required"`
	GameAreaID pq.Int64Array `json:"game_area_id" form:"game_area_id" binding:"required"`
}

// On3DGameConfigUpdateRequest 3D大厅游戏子游戏修改
func On3DGameConfigUpdateRequest(c *gin.Context) {
	req := &On3DGameConfigUpdateReq{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": err.Error()})
		return
	}
	defer LogStat(`On3DGameConfigUpdateRequest`, c, req, time.Now())
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0400", "errmsg": "参数错误," + err.Error()})
		return
	}
	configdb := c.MustGet("configdb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	err := service.On3DGameConfigUpdateRequest(configdb, pool, req.GameID, req.GameAreaID)
	if err != nil {
		logging.Errorf("On3DGameConfigUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "C0300", "errmsg": "添加数据库失败失败," + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":  "0",
		"errmsg": "更新成功",
	})
}
func (c *On3DGameConfigUpdateReq) Validate() error {
	if len(c.GameAreaID) != 2 {
		return errors.New("必须传入两个子游戏名")
	}
	return nil
}
