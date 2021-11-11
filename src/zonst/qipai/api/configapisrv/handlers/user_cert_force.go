package handlers

import (
	"errors"
	"net/http"
	"time"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/models"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

// UserCertForceStatusReq req
type UserCertForceStatusReq struct {
	GameID int                       `json:"game_id" form:"game_id" binding:"required"` // 平台ID
	Status utils.UserCertForceStatus `json:"status" form:"status"`                      // 开启状态
	Scope  utils.UserCertForceScope  `json:"scope" form:"scope"`                        // 执行范围
}

// Validate validate
func (req *UserCertForceStatusReq) Validate() error {
	if req.Status != utils.OpenUserCertForce && req.Status != utils.CloseUserCertForce {
		return errors.New(`错误的开关类型`)
	}
	if req.Scope != utils.InsideJiangXi &&
		req.Scope != utils.OutsideJiangXi &&
		req.Scope != utils.WholeCountry {
		return errors.New(`错误的执行范围`)
	}
	return nil
}

// UpdateUserCertForce 更新强制实名认证
func UpdateUserCertForce(c *gin.Context) {
	var req UserCertForceStatusReq

	if err := c.Bind(&req); err != nil {
		log.Errorf("UpdateUserCertForce: 参数绑定出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "1", "errmsg": "参数绑定出错"})
		return
	}
	if err := req.Validate(); err != nil {
		log.Errorf("UpdateUserCertForce: 参数验证出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "2", "errmsg": "参数验证出错，" + err.Error()})
		return
	}
	defer LogStat(`UpdateUserCertForce`, c, req, time.Now())

	qipaidb := c.MustGet("usercertforcedb").(*gorm.DB)

	oldInfo, err := models.FindUserCertForceStatusByGameID(qipaidb, req.GameID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("UpdateUserCertForce: 获取用户强制实名信息出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "3", "errmsg": "获取用户强制实名信息出错"})
		return
	}
	// scope为过渡字段现在都强制设置为utils.WholeCountry
	if err := models.UpdateUserCertForceStatus(qipaidb, req.GameID, req.Status, utils.WholeCountry); err != nil {
		log.Errorf("UpdateUserCertForce: 更新用户强制实名信息出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "4", "errmsg": "更新用户强制实名信息出错"})
		return
	}

	pool := c.MustGet(`usercertforce`).(*redis.Pool)
	if err := cache.DeleteUserCertForceScope(pool, req.GameID); err != nil {
		log.Errorf("UpdateUserCertForce: 删除用户强制实名信息范围缓存出错 err: %v\n", err.Error())
	}
	if err := cache.DeleteUserCertForceStatus(pool, req.GameID); err != nil {
		log.Errorf("UpdateUserCertForce: 删除用户强制实名信息状态缓存出错 err: %v\n", err.Error())
	}
	var (
		opUserID   int
		opUserName string
	)
	token, exist := c.Get("token")
	if exist {
		info := token.(middlewares.Claims)
		opUserID = int(info.UserID)
		opUserName = info.UserName
	}
	updateLog := models.UserCertForceStatusUpdateLog{
		GameID:       req.GameID,
		OpUserId:     opUserID,
		OpUsername:   opUserName,
		BeforeStatus: oldInfo.Status,
		AfterStatus:  req.Status,
		BeforeScope:  oldInfo.Scope,
		AfterScope:   utils.WholeCountry, // scope为过渡字段现在都强制设置为utils.WholeCountry
	}
	logdb := c.MustGet("logdb").(*sqlx.DB)

	if err := models.CreateUserCertForceStatusUpdateLog(logdb, updateLog); err != nil {
		log.Errorf("UpdateUserCertForce: 保存操作日志出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "5", "errmsg": "保存操作日志出错"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "OK"})
}

// FindUserCertForceReq req
type FindUserCertForceReq struct {
	GameIDS []int `json:"game_ids" form:"game_ids[]"`
}

// FindUserCertForce 查询
func FindUserCertForce(c *gin.Context) {
	var req FindUserCertForceReq
	if err := c.Bind(&req); err != nil {
		log.Errorf("FindUserCertForce: 参数绑定出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "1", "errmsg": "参数绑定出错"})
		return
	}
	defer LogStat(`FindUserCertForce`, c, req, time.Now())
	qipaidb := c.MustGet("usercertforcedb").(*gorm.DB)

	status, err := models.FindUserCertForceStatus(qipaidb, req.GameIDS)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("FindUserCertForce: 查询用户强制实名信息出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "1", "errmsg": "查询用户强制实名信息出错"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "OK", "data": status})
}
