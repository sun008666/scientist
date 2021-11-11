package handlers

import (
	"net/http"
	"time"
	"zonst/qipai/api/configapisrv/cache"

	"zonst/logging"

	"zonst/qipai/api/configapisrv/models"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jmoiron/sqlx"
)

// 游戏平台屏蔽

type GetGameForbidListRequest struct {
	GameIDs []int32 `json:"game_id_list" binding:"required"`
}

// OnGetGameForbidListRequest 获取游戏平台开启/关闭状态列表
func OnGetGameForbidListRequest(c *gin.Context) {
	req := &GetGameForbidListRequest{}

	if err := c.Bind(req); err != nil {
		log.Errorf("OnGetGameForbidListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	defer LogStat("OnGetGameForbidListRequest", c, req, time.Now())

	userphoneCache := c.MustGet("userphone").(*redis.Pool)
	m, err := cache.BatchGetGameForbidStatus(userphoneCache, req.GameIDs)
	if err != nil {
		logging.Errorln("fail forbid:", err)
	}

	type Data struct {
		GameID int32 `json:"game_id"`
		Forbid int32 `json:"forbid"`
	}

	s := make([]*Data, 0, len(m)+1)
	for k, v := range m {
		s = append(s, &Data{GameID: k, Forbid: v})
	}

	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": s})
}

type UpdateGameForbidRequest struct {
	GameID   int32  `json:"game_id" binding:"required"`
	Forbid   int32  `json:"forbid" binding:""`
	UserID   int32  `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// OnUpdateGameForbidRequest 设置游戏平台开启/关闭状态
func OnUpdateGameForbidRequest(c *gin.Context) {
	req := &UpdateGameForbidRequest{}

	if err := c.Bind(req); err != nil {
		log.Errorf("OnUpdateGameForbidRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	defer LogStat("OnUpdateGameForbidRequest", c, req, time.Now())

	userphoneCache := c.MustGet("userphone").(*redis.Pool)
	statDB := c.MustGet("logdb").(*sqlx.DB)

	m, err := cache.BatchGetGameForbidStatus(userphoneCache, []int32{req.GameID})
	if err != nil {
		logging.Errorln(err)
	}
	if m == nil {
		m = make(map[int32]int32, 1)
	}

	err = cache.UpdateGameForbidStatus(userphoneCache, req.GameID, req.Forbid)
	if err != nil {
		logging.Errorln("fail update forbid:", err)
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "更新失败"})
		return
	}

	// 落日志，由于日志量不大，直接插入数据库
	err = models.InsertGameForbidSetLog(statDB, req.GameID, req.UserID, m[req.GameID], req.Forbid, req.Username)
	if err != nil {
		logging.Errorln(err)
	}

	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}
