package handlers

import (
	"net/http"
	"time"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/models"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
)

type OnPopupControlAddReq struct {
	GameID   int `json:"game_id" binding:"required"`
	Strategy int `json:"strategy"`
}

func OnPopupControlAddRequest(c *gin.Context) {
	req := &OnPopupControlAddReq{}
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupControlAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)

	if err := db.Table("pop_up_control").Create(req).Error; err != nil {
		log.Errorf("OnPopupControlAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type OnPopupControlListReq struct {
	GameID int `json:"game_id" binding:"required"`
}

// OnPopupControlListRequest 弹窗控制状态
func OnPopupControlListRequest(c *gin.Context) {
	req := &OnPopupControlListReq{}
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupControlListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	data := &models.PopUpControl{}
	if err := db.Table("pop_up_control").Where("game_id = ?", req.GameID).First(data).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnPopupControlListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": data})

}

type OnPopupControlUpdateReq struct {
	GameID   int `json:"game_id" binding:"required"`
	Strategy int `json:"strategy"`
}

func OnPopupControlUpdateRequest(c *gin.Context) {
	req := &OnPopupControlUpdateReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPopupControlUpdateRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupControlUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)

	//判断弹窗控制是否存在
	var data []models.PopUpControl
	if err := db.Table("pop_up_control").Where("game_id = ?", req.GameID).Find(&data).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnPopupControlListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	if len(data) > 0 {
		//存在修改
		if err := db.Table("pop_up_control").Where("game_id = ?", req.GameID).Update("strategy", req.Strategy).Error; err != nil {
			log.Errorf("OnPopupControlUpdateRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}

	} else {
		//不存在添加
		if err := db.Table("pop_up_control").Create(req).Error; err != nil {
			log.Errorf("OnPopupControlAddRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
	}
	cacheName := "pop_up_strategy"
	if !models.RefleshCache(c, cacheName) {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新缓存出错"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}
