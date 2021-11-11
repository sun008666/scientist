package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/models"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
)

type ChangeDistinctConfig struct {
	ID           int32    `json:"id" form:"id"`
	GameID       int32    `json:"game_id" form:"game_id" binding:"required"`
	Province     string   `json:"province" form:"province" binding:"required"`
	City         string   `json:"city" form:"city" binding:"required"`
	District     []string `json:"district" form:"district"`
	WXName       string   `json:"wx_name" form:"wx_name" binding:"required"`
	WxAccount    string   `json:"wx_account" form:"wx_account" binding:"required"`
	TitleContent string   `json:"title_content" form:"title_content" binding:"required"`
}

func OnDistinctGameConfigChangeRequest(c *gin.Context) {
	req := &ChangeDistinctConfig{}
	defer LogStat("OnDistinctGameConfigChangeRequest", c, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnDistinctGameConfigChangeRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	qipaidb := c.MustGet("qipaidb").(*gorm.DB)
	cache := c.MustGet("distrinctconfig").(*cache.Store)
	// 修改
	if req.ID != 0 {
		data, err := models.GetDistrictConfigByID(qipaidb, req.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Errorf("OnDistinctGameConfigChangeRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "查询数据失败," + err.Error()})
			return
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "要修改的数据不存在"})
			return
		}
		update := make(map[string]interface{})
		if data.Province != req.Province {
			update["province"] = req.Province
		}
		if data.City != req.City {
			update["city"] = req.City
		}

		if data.Val(req.District) {
			config, err := cache.DistrictTable.Get(req.GameID)
			if err != nil && err.Error() != "key not found" {
				c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "错误:" + err.Error()})
				return
			}
			if len(req.District) == 0 {
				for _, v := range config {
					if req.Province == v.Province && req.City == v.City && len(v.DistrictArray) == 0 {

						log.Errorf("OnDistinctGameConfigChangeRequest:要修改的区/(县)游戏配置已存在\n")
						c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "要修改的区/(县)游戏配置已存在"})
						return

					}

				}
			} else {
				for _, va := range req.District {
					for _, v := range config {
						if req.Province == v.Province && req.City == v.City && v.ID != req.ID {
							if len(v.DistrictArray) == 0 {
								log.Errorf("OnDistinctGameConfigChangeRequest:要修改的区/(县)游戏配置已存在\n")
								c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "要修改的区/(县)游戏配置已存在"})
								return
							}
							for _, dis := range v.DistrictArray {
								if va == dis {
									log.Errorf("OnDistinctGameConfigChangeRequest:要修改的区/(县)游戏配置已存在\n")
									c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "要修改的区/(县)游戏配置已存在"})
									return
								}

							}
						}

					}

				}
			}

			jsonByte, err := json.Marshal(req.District)
			if err != nil {
				log.Errorf("OnDistinctGameConfigChangeRequest:err:%v\n", err)
				c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "序列化错误," + err.Error()})
				return
			}
			update["district"] = string(jsonByte)

		}
		if data.TitleContent != req.TitleContent {
			update["title_content"] = req.TitleContent
		}
		if data.WXName != req.WXName {
			update["wx_name"] = req.WXName
		}
		if data.WxAccount != req.WxAccount {
			update["wx_account"] = req.WxAccount
		}

		if err := qipaidb.Debug().Table("game_district_config").Where("id=?", req.ID).Updates(update).Error; err != nil {
			log.Errorf("OnDistinctGameConfigChangeRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "修改失败," + err.Error()})
			return
		}
		//刷新缓存
		updateData, err := models.GetDistrictConfigList(qipaidb, req.GameID, "", "")
		if err != nil {
			log.Errorf("OnDistinctGameConfigChangeRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "获取数据失败"})
			return
		}
		cache.DistrictTable.Set(req.GameID, updateData)
		c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "修改成功"})
		return
	}
	// 增加
	config, err := cache.DistrictTable.Get(req.GameID)
	if err != nil && err.Error() != "key not found" {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "错误:" + err.Error()})
		return
	}
	if len(req.District) == 0 {
		for _, v := range config {
			if req.Province == v.Province && req.City == v.City {
				log.Errorf("OnDistinctGameConfigChangeRequest:要添加的区/(县)游戏配置已存在\n")
				c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "要添加的区/(县)游戏配置已存在"})
				return
			}

		}

	} else {
		for _, va := range req.District {
			for _, v := range config {
				if req.Province == v.Province && req.City == v.City {
					if len(v.DistrictArray) == 0 {
						log.Errorf("OnDistinctGameConfigChangeRequest:要添加的区/(县)游戏配置已存在\n")
						c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "要修改的区/(县)游戏配置已存在"})
						return
					}
					for _, dis := range v.DistrictArray {
						if va == dis {
							log.Errorf("OnDistinctGameConfigChangeRequest:要添加的区/(县)游戏配置已存在\n")
							c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "要修改的区/(县)游戏配置已存在"})
							return
						}

					}
				}

			}

		}
	}
	jsonByte, err := json.Marshal(req.District)
	if err != nil {
		log.Errorf("OnDistinctGameConfigChangeRequest:err:%v\n", err)
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "序列化错误," + err.Error()})
		return
	}
	add := models.AddGameDistrictConfig{
		GameID:       req.GameID,
		Province:     req.Province,
		City:         req.City,
		District:     string(jsonByte),
		WXName:       req.WXName,
		WxAccount:    req.WxAccount,
		TitleContent: req.TitleContent,
	}
	if err := qipaidb.Debug().Table("game_district_config").Create(&add).Error; err != nil {
		log.Errorf("OnDistinctGameConfigChangeRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "修改失败," + err.Error()})
		return
	}
	//刷新缓存
	addData, err := models.GetDistrictConfigList(qipaidb, req.GameID, "", "")
	if err != nil {
		log.Errorf("OnDistinctGameConfigChangeRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "获取数据失败"})
		return
	}
	cache.DistrictTable.Set(req.GameID, addData)
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "添加成功"})
	return

}

type DeleteDataFromReq struct {
	ID     int32 `json:"id" form:"id" binding:"required"`
	GameID int32 `json:"game_id" form:"game_id" binding:"required"`
}

func OnDistinctGameConfigDelRequest(c *gin.Context) {
	req := &DeleteDataFromReq{}
	defer LogStat("OnDistinctGameConfigDelRequest", c, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnDistinctGameConfigDelRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	qipaidb := c.MustGet("qipaidb").(*gorm.DB)
	err := models.DeleteGameDistrictConfig(qipaidb, req.ID)
	if err != nil {
		log.Errorf("OnDistinctGameConfigDelRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "删除数据失败," + err.Error()})
		return

	}
	deleteData, err := models.GetDistrictConfigList(qipaidb, req.GameID, "", "")
	if err != nil {
		log.Errorf("OnDistinctGameConfigDelRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "获取数据失败"})
		return
	}
	cache := c.MustGet("distrinctconfig").(*cache.Store)
	cache.DistrictTable.Set(req.GameID, deleteData)

	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "删除成功"})
	return

}

type DistinctGameConfigGetReq struct {
	GameID   int32  `json:"game_id" form:"game_id" binding:"required"`
	Province string `json:"province" form:"province"`
	City     string `json:"city" form:"city"`
}

func OnDistinctGameConfigGetRequest(c *gin.Context) {
	req := &DistinctGameConfigGetReq{}
	defer LogStat("OnPopupWhileListAddRequest", c, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPopupWhileListAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	qipaidb := c.MustGet("qipaidb").(*gorm.DB)
	data, err := models.GetDistrictConfigList(qipaidb, req.GameID, req.Province, req.City)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnDistinctGameConfigDelRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "获取数据失败," + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "", "data": data})
	return

}

type GetGameDistrictConfigReq struct {
	GameID   int32  `json:"game_id" form:"game_id" binding:"required"`
	Province string `json:"province" form:"province" binding:"required"`
	City     string `json:"city" form:"city" binding:"required"`
	District string `json:"district" form:"district" binding:"required"`
}

func OnGetGameDistrictConfigRequest(c *gin.Context) {
	req := &GetGameDistrictConfigReq{}
	defer LogStat("OnGetGameDistrictConfigRequest", c, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGetGameDistrictConfigRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "参数不匹配," + err.Error()})
		return
	}
	cache := c.MustGet("distrinctconfig").(*cache.Store)

	config, err := cache.DistrictTable.GetDistrictConfig(req.Province, req.City, req.District, req.GameID)
	if err != nil && err.Error() != "key not found" {
		log.Errorf("OnGetGameDistrictConfigRequest:%v\n", err)
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "获取数据出错," + err.Error()})
		return
	}
	log.Debugf("game_id:%v,province:%v,city:%v,district:%v,err:%v\n", req.GameID, req.Province, req.City, req.District, err)
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "", "data": config})
	return
}
