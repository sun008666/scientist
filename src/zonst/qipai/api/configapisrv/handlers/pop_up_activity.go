package handlers

import (
	"errors"
	"net/http"
	"time"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/models"

	"github.com/lib/pq"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
)

// PopupActivity 弹窗活动
type PopupActivity struct {
	ID             int            `json:"id"`
	GameID         int            `json:"game_id" binding:"required"`
	Title          string         `json:"title" binding:"required"`      //标题
	ImgURL         string         `json:"img_url"`                       //图片链接
	StartTime      int            `json:"start_time" binding:"required"` //开始时间
	EndTime        int            `json:"end_time" binding:"required"`   //结束时间
	LinkTo         int            `json:"link_to" binding:"required"`    //跳转至， 任务类型
	LinkURL        string         `json:"link_url"`                      //外部链接URL
	Show           int            `json:"show"`                          //显示
	ShowWho        int            `json:"show_who"`                      //显示属性，新用户，所有人，亲友圈创始人
	Position       int            `json:"position"`
	RedDot         int            `json:"red_dot" binding:"required"`
	PlayerType     pq.Int64Array  `json:"player_type"`
	PlayerIdentity pq.Int64Array  `json:"player_identity"`
	ShowArea       pq.StringArray `json:"show_area"`
	ShowGameAreaID pq.Int64Array  `json:"show_game_area_id"`
	IsWhiteList    bool           `json:"is_white_list"` // 是否是白名单    true false
	ShowSystem     string         `json:"show_system"`   // 显示的操作系统 all android ios
	GameList       pq.Int64Array  `json:"game_list"`
}

func (req *PopupActivity) Validate() error {
	if req.StartTime > req.EndTime {
		return errors.New(`开始时间大于结束时间`)
	}
	return nil
}

// OnPopupActivityAddRequest 弹窗活动添加
func OnPopupActivityAddRequest(c *gin.Context) {
	req := &PopupActivity{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPopupActivityAddRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupActivityAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		log.Errorf("OnPopupActivityAddRequest: 参数验证出错 err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "2", "errmsg": "参数验证出错，" + err.Error()})
		return
	}
	log.Debugf("req:%+v\n", req)

	obj := &models.PopUpActivityTable{
		GameID:         req.GameID,
		Title:          req.Title,
		ImgURL:         req.ImgURL,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		LinkTo:         req.LinkTo,
		LinkURL:        req.LinkURL,
		IsShow:         req.Show,
		SendToClient:   req.Show,
		ShowWho:        req.ShowWho,
		RedDot:         req.RedDot,
		PlayerType:     req.PlayerType,
		PlayerIdentity: req.PlayerIdentity,
		ShowArea:       req.ShowArea,
		ShowGameAreaID: req.ShowGameAreaID,
		IsWhiteList:    req.IsWhiteList,
		ShowSystem:     req.ShowSystem,
		GameList:       req.GameList,
	}
	err := models.AddPopUpActivity(c, obj)
	if err != nil {
		log.Errorf("OnPopupActivityAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "添加配置出错," + err.Error()})
		return
	}
	cacheName := "pop_up_activity"
	if !models.RefleshCache(c, cacheName) {
		log.Errorf("AddPopUpActivity: err: %v\n", "刷新缓存出错")
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新缓存出错"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
	return
}

// OnPopupActivityListReq 弹窗活动列表
type OnPopupActivityListReq struct {
	GameID int `json:"game_id" binding:"required"`
}

// OnPopupActivityListRequest 弹窗活动列表
func OnPopupActivityListRequest(c *gin.Context) {
	req := &OnPopupActivityListReq{}
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupActivityListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	log.Debugf("req:%+v\n", req)
	db := c.MustGet("qipaidb").(*gorm.DB)
	obj := []models.PopUpActivityTable{}
	if err := db.Debug().Table("pop_up_activity").Where("game_id = ? or (?=any(game_list))", req.GameID, req.GameID).Order("is_show desc, position desc").Find(&obj).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnPopupActivityListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	data := []PopupActivity{}
	for _, v := range obj {
		tmp := PopupActivity{
			ID:             v.ID,
			GameID:         v.GameID,
			Title:          v.Title,
			ImgURL:         v.ImgURL,
			StartTime:      v.StartTime,
			EndTime:        v.EndTime,
			LinkTo:         v.LinkTo,
			LinkURL:        v.LinkURL,
			Show:           v.IsShow,
			ShowWho:        v.ShowWho,
			Position:       v.Position,
			RedDot:         v.RedDot,
			PlayerType:     v.PlayerType,
			PlayerIdentity: v.PlayerIdentity,
			ShowArea:       v.ShowArea,
			ShowGameAreaID: v.ShowGameAreaID,
			IsWhiteList:    v.IsWhiteList,
			ShowSystem:     v.ShowSystem,
			GameList:       v.GameList,
		}
		data = append(data, tmp)
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": data})
}

// OnPopupActivityUpdateRequest 弹窗活动修改
func OnPopupActivityUpdateRequest(c *gin.Context) {
	req := &PopupActivity{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPopupActivityUpdateRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupActivityUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	log.Debugf("req:%+v\n", req)
	obj := &models.PopUpActivityTable{
		GameID:         req.GameID,
		Title:          req.Title,
		ImgURL:         req.ImgURL,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		LinkTo:         req.LinkTo,
		LinkURL:        req.LinkURL,
		IsShow:         req.Show,
		SendToClient:   req.Show,
		ShowWho:        req.ShowWho,
		Position:       req.Position,
		RedDot:         req.RedDot,
		PlayerType:     req.PlayerType,
		PlayerIdentity: req.PlayerIdentity,
		ShowArea:       req.ShowArea,
		ShowGameAreaID: req.ShowGameAreaID,
		IsWhiteList:    req.IsWhiteList,
		ShowSystem:     req.ShowSystem,
		GameList:       req.GameList,
	}
	err := models.UpdatePopActivity(c, req.ID, obj)
	if err != nil {
		log.Errorf("OnPopupActivityUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	cacheName := "pop_up_activity"
	if !models.RefleshCache(c, cacheName) {
		log.Errorf("AddPopUpActivity: err: %v\n", "刷新缓存出错")
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新缓存出错"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
	return
}

type OnPopupActivityDeleteReq struct {
	ID int `json:"id" binding:"required"`
}

func OnPopupActivityDeleteRequest(c *gin.Context) {
	req := &OnPopupActivityDeleteReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPopupActivityDeleteRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupActivityDeleteRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	log.Debugf("req:%+v\n", req)
	db := c.MustGet("qipaidb").(*gorm.DB)
	id := []int{0}
	if err := db.Debug().Table("pop_up_activity").Delete(req).Pluck("id", &id).Error; err != nil {
		log.Errorf("OnPopupActivityDeleteRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	cacheName := "pop_up_activity"
	if !models.RefleshCache(c, cacheName) {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新缓存出错"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}
