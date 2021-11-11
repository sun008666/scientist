package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/service"
)

// AndroidOrIosLoadUrlReq 游戏配置列表
type AndroidOrIosLoadUrlReq struct {
	GameID int `json:"game_id" form:"game_id" binding:"required"`
}

// AndroidOrIosDownLoadUrl 3D大厅游戏展示位
func AndroidOrIosDownLoadUrl(c *gin.Context) {
	req := &AndroidOrIosLoadUrlReq{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": "参数绑定错误"})
		return
	}
	defer LogStat(`AndroidOrIosLoadUrlReq`, c, req, time.Now())
	qipaidb := c.MustGet("qipaidb").(*gorm.DB)
	data, err := service.FindAndroidOrIosLoadUrl(qipaidb, req.GameID)
	if err != nil {
		logging.Errorf("AndroidOrIosLoadUrl:Failed to query the database err: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"errno": "C0300", "errmsg": "查询数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errno":  "0",
		"errmsg": "",
		"data":   data,
	})
}
