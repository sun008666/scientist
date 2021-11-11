package handlers

import (
	"github.com/fwhezfwhez/errorx"
	"net/http"
	"time"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/dependency/errs"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/models"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
)

type AddWhiteListReq struct {
	GameID int32  `json:"game_id" binding:"required"`
	UserID int32  `json:"user_id" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}

func OnPopupWhileListAddRequest(c *gin.Context) {
	req := &AddWhiteListReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPopupWhileListAddRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupWhileListAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	qipaidb := c.MustGet("qipaidb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	// 判断白名单是否存在
	total, _, err := models.GetWhiteList(qipaidb, req.GameID, req.UserID, 1, 3)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnPopupWhileListAddRequest: 获取弹窗白名单错误,err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	if total > 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "用户已在白名单中"})
		return
	}
	// 插入数据库
	add := models.WhiteList{
		GameID:  req.GameID,
		UserID:  req.UserID,
		Remark:  req.Remark,
		LogDate: time.Now().Format("2006-01-02"),
		LogTime: int32(time.Now().Unix()),
	}
	if err := qipaidb.Debug().Table("popup_white_list").Create(&add).Error; err != nil {
		log.Errorf("OnPopupWhileListAddRequest: 获取弹窗白名单错误,err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "添加白名单失败," + err.Error()})
		return
	}

	// 插入缓存
	err = models.SetPopupUserWhiteList(pool, req.UserID)
	if err != nil {
		log.Errorf("OnPopupWhileListAddRequest: 获取弹窗白名单错误,err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "添加白名单缓存失败," + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "添加成功"})
}

type GetWhiteListReq struct {
	GameID    int32 `json:"game_id" binding:"required"`
	PageID    int32 `json:"page_id" binding:"required"`
	PageCount int32 `json:"page_count" binding:"required"`
}

func OnGetPopupWhileListRequest(c *gin.Context) {
	req := &GetWhiteListReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPopupWhileListAddRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGetPopupWhileListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	qipaidb := c.MustGet("qipaidb").(*gorm.DB)
	total, data, err := models.GetWhiteList(qipaidb, req.GameID, 0, req.PageID, req.PageCount)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnGetPopupWhileListRequest: 获取弹窗白名单错误,err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "", "total": total, "data": data})
	return
}

type DeleteWhiteListReq struct {
	GameID int32 `json:"game_id" binding:"required"`
	UserID int32 `json:"user_id" binding:"required"`
}

func OnDeletePopupWhileListRequest(c *gin.Context) {
	req := &DeleteWhiteListReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnDeletePopupWhileListRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnDeletePopupWhileListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	pool := c.MustGet("whiteListCache").(*redis.Pool)

	// 删数据库
	err := models.DeleteWhiteList(db, req.GameID, req.UserID)
	if err != nil {
		log.Errorf("OnDeletePopupWhileListRequest:删除白名单失败, err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	// 删缓存
	err = models.RemovePopupUserWhitelist(pool, req.UserID)
	if err != nil {
		log.Errorf("OnDeletePopupWhileListRequest:删除缓存白名单失败, err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "删除成功"})
	return
}

type IsWhiteListReq struct {
	GameID int32 `json:"game_id" form:"game_id" binding:"required"`
	UserID int32 `json:"user_id" form:"user_id" binding:"required"`
}

func OnIsPopuWhiteListRequest(c *gin.Context) {
	req := &IsWhiteListReq{}
	//userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnIsPopuWhiteListRequest", c, req.UserID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnIsPopuWhiteListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	pool := c.MustGet("whiteListCache").(*redis.Pool)
	config := c.MustGet("config").(*config.Config)
	userIdentity, err := utils.GetUserIdentityInfo(config.UserInfo, c.ClientIP(), req.GameID, req.UserID)
	if err != nil {
		logging.Errorf("OnSendRedBagToUserRequest:获取用户身份信息错误:%v,err:%v\n", c.Request.RequestURI, err)
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "获取用户身份信息错误," + err.Error()})
		return
	}

	type UserInfoBack struct {
		Data         utils.UserInfo `json:"data"`
		Address      utils.Address  `json:"address"`
		IsWhiteKList bool           `json:"is_white_list"`
	}

	flag, err := models.IsWhiteUser(pool, req.UserID)
	if err != nil {
		errs.SaveError(errorx.Wrap(err), map[string]interface{}{
			"req": req,
		})
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	data := UserInfoBack{}
	data.IsWhiteKList = flag
	data.Address = userIdentity.Address
	data.Data = userIdentity.Data

	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "", "data": data})
	return
}
