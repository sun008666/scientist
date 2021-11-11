package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/models"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/lib/pq"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
)

// OnPayListAddReq 充值列表-添加
type OnPayListAddReq struct {
	GameID              int    `json:"game_id" binding:"required"`
	GameAreaID          string `json:"game_area_id"`
	CardNum             int    `json:"card_num"`
	Status              int    `json:"status"`
	Price               int    `json:"price"`
	AddTime             string `json:"add_time"`
	Category            string `json:"category"`
	PresentRoomCard     int    `json:"present_room_card"`
	EndTime             string `json:"end_time"`
	Represent           string `json:"represent"`
	AppleKey            string `json:"apple_key"`
	IsWhite             bool   `json:"is_white"`
	IsBlack             bool   `json:"is_black"`
	WhiteList           []int  `json:"white_list"`
	BlackList           []int  `json:"black_list"`
	CardPackID          int32  `json:"card_pack_id"`
	IsSeniorAgent       int32  `json:"is_senior_agent"`
	PurchaseLimitNumber int32  `json:"purchase_limit_number"`
}

// OnPayListAddRequest 充值列表-添加
func OnPayListAddRequest(c *gin.Context) {
	req := &OnPayListAddReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPayListAddRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPayListAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	if req.GameAreaID == "" {
		req.GameAreaID = "-1" //普通房卡套餐
		if req.CardPackID > 0 {
			req.GameAreaID = "-2" //卡包套餐
		}
	}
	whiteList := removeDuplicate(req.WhiteList)
	blackList := removeDuplicate(req.BlackList)
	if err := ValidateBlackWhiteList(whiteList, blackList); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": fmt.Sprintf("%v", err)})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	obj := &models.RoomCardProductAddRequest{
		GameID:     req.GameID,
		GameAreaID: req.GameAreaID,
		CardNum:    req.CardNum,
		Status:     req.Status,
		Price:      req.Price,
		// AddTime:         time.Unix(req.AddTime, 0).Format("2006-01-02 15:04:05"),
		AddTime:         req.AddTime,
		Category:        req.Category,
		PresentRoomCard: req.PresentRoomCard,
		// EndTime:         time.Unix(req.EndTime, 0).Format("2006-01-02 15:04:05"),
		EndTime:             req.EndTime,
		Represent:           req.Represent,
		AppleKey:            req.AppleKey,
		IsWhite:             req.IsWhite,
		IsBlack:             req.IsBlack,
		WhiteList:           changeToPqArray(whiteList),
		BlackList:           changeToPqArray(blackList),
		CardPackID:          req.CardPackID,
		IsSeniorAgent:       req.IsSeniorAgent,
		PurchaseLimitNumber: req.PurchaseLimitNumber,
	}
	if err := db.Debug().Table("room_card_product_list").Create(obj).Error; err != nil {
		log.Errorf("OnPayListAddRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	afterContent, err1 := models.GetLastProductConfig(db, req.GameID, req.GameAreaID, req.Category)
	if err1 == nil {
		config := c.MustGet("config").(*config.Config)
		args := make(map[string]interface{})
		args["game_id"] = req.GameID
		args["module_id"] = 1
		before, _ := json.Marshal(models.RoomCardProductList{})
		args["before_content"] = string(before)
		args["operate_id"] = userID
		args["op_type"] = 1
		args["after_content"] = afterContent
		err := utils.AddConfigLog(config.ConfigUpdateLogAPI, args)
		if err != nil {
			log.Errorf("OnPayListAddRequest-AddConfigLog: err: %v\n", err)
		}
	} else {
		log.Errorf("OnPayListAddRequest-GetLastProductConfig: err: %v\n", err1)
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type OnPayListListReq struct {
	GameID         int    `json:"game_id" binding:"required"`
	GameAreaID     int    `json:"game_area_id"`
	SourceType     string `json:"source_type"`
	IsAppleProduct bool   `json:"is_apple_product"`
}

// OnPayListListRequest 充值列表-列表
func OnPayListListRequest(c *gin.Context) {
	req := &OnPayListListReq{}
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPayListListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)

	//data := &[]models.RoomCardProductList{}
	// if err := db.Table("room_card_product_list").Where("game_id = ?", req.GameID).Order("status desc, category, price").Find(data).Error; err != nil {
	// 	log.Errorf("OnPayListListRequest: err: %v\n", err.Error())
	// 	c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
	// 	return
	// }

	data, err := models.GetProductList(db, req.GameID, req.GameAreaID, req.SourceType, req.IsAppleProduct)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnPayListListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": data})
}

// OnPayListUpdateReq 充值列表-修改
type OnPayListUpdateReq struct {
	ID                  int    `json:"id" binding:"required"`
	GameID              int    `json:"game_id" binding:"required"`
	GameAreaID          string `json:"game_area_id"`
	CardNum             int    `json:"card_num"`
	Status              int    `json:"status"`
	Price               int    `json:"price"`
	AddTime             string `json:"add_time"`
	Category            string `json:"category"`
	PresentRoomCard     int    `json:"present_room_card"`
	EndTime             string `json:"end_time"`
	Represent           string `json:"represent"`
	AppleKey            string `json:"apple_key"`
	IsWhite             bool   `json:"is_white"`
	IsBlack             bool   `json:"is_black"`
	WhiteList           []int  `json:"white_list"`
	BlackList           []int  `json:"black_list"`
	CardPackID          int32  `json:"card_pack_id"`
	IsSeniorAgent       int32  `json:"is_senior_agent"`
	PurchaseLimitNumber int32  `json:"purchase_limit_number"`
}

// OnPayListUpdateRequest 充值列表-修改
func OnPayListUpdateRequest(c *gin.Context) {
	req := &OnPayListUpdateReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPayListUpdateRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPayListUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	beforeContent, err1 := models.GetProductConfigByID(db, req.ID)
	if err1 != nil {
		log.Errorf("OnPayListUpdateRequest-GetProductConfigByID: err: %v\n", err1)
	}
	if req.GameAreaID == "" {
		req.GameAreaID = "-1" //普通房卡套餐
		if req.CardPackID > 0 {
			req.GameAreaID = "-2" //卡包套餐
		}
	}
	if req.CardNum <= 0 && req.PresentRoomCard <= 0 {
		log.Errorf("OnPayListUpdateRequest: req: %+v\n", req)
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "购买房卡数和赠送房卡数至少要有一个大于0"})
		return
	}
	whiteList := removeDuplicate(req.WhiteList)
	blackList := removeDuplicate(req.BlackList)
	if err := ValidateBlackWhiteList(whiteList, blackList); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": "A0410", "errmsg": fmt.Sprintf("%v", err)})
	}
	//obj := &models.RoomCardProductAddRequest{
	//	GameID:     req.GameID,
	//	GameAreaID: req.GameAreaID,
	//	//CardNum:    req.CardNum,
	//	//Status:     req.Status,
	//	Price: req.Price,
	//	// AddTime:         time.Unix(req.AddTime, 0).Format("2006-01-02 15:04:05"),
	//	AddTime:  req.AddTime,
	//	Category: req.Category,
	//	//PresentRoomCard: req.PresentRoomCard,
	//	// EndTime:         time.Unix(req.EndTime, 0).Format("2006-01-02 15:04:05"),
	//	EndTime:   req.EndTime,
	//	Represent: req.Represent,
	//	AppleKey:  req.AppleKey,
	//	//CardPackID: req.CardPackID,
	//}
	update := make(map[string]interface{})
	update["is_white"] = req.IsWhite
	update["is_black"] = req.IsBlack
	update["white_list"] = changeToPqArray(whiteList)
	update["black_list"] = changeToPqArray(blackList)
	update["card_num"] = req.CardNum
	update["present_room_card"] = req.PresentRoomCard
	update["game_area_id"] = req.GameAreaID
	update["card_pack_id"] = req.CardPackID
	update["status"] = req.Status
	update["is_senior_agent"] = req.IsSeniorAgent
	update["game_id"] = req.GameID
	update["price"] = req.Price
	update["add_time"] = req.AddTime
	update["category"] = req.Category
	update["end_time"] = req.EndTime
	update["represent"] = req.Represent
	update["apple_key"] = req.AppleKey
	update["purchase_limit_number"] = req.PurchaseLimitNumber

	if err := db.Debug().Table("room_card_product_list").Where("id = ?", req.ID).Updates(update).Error; err != nil {
		log.Errorf("OnPayListUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	afterContent, err2 := models.GetProductConfigByID(db, req.ID)
	if err1 == nil && err2 == nil {
		config := c.MustGet("config").(*config.Config)
		args := make(map[string]interface{})
		before, _ := json.Marshal(beforeContent)
		args["before_content"] = string(before)
		args["game_id"] = beforeContent.GameID
		args["module_id"] = 1
		after, _ := json.Marshal(afterContent)
		args["after_content"] = string(after)
		args["operate_id"] = userID
		args["op_type"] = 2
		err := utils.AddConfigLog(config.ConfigUpdateLogAPI, args)
		if err != nil {
			log.Errorf("OnPayListDeleteRequest-AddConfigLog: err: %v\n", err)
		}
	} else {
		log.Errorf("OnPayListUpdateRequest-GetProductConfigByID: err1: %v,err2:%v\n", err1, err2)
	}

	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type OnPayListDeleteReq struct {
	ID int `json:"id" binding:"required"`
}

// OnPayListDeleteRequest 充值列表-删除
func OnPayListDeleteRequest(c *gin.Context) {
	req := &OnPayListDeleteReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPayListDeleteRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPayListDeleteRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	beforeContent, err1 := models.GetProductConfigByID(db, req.ID)
	if err1 != nil {
		log.Errorf("OnPayListDeleteRequest-GetProductConfigByID: err: %v\n", err1)
	}
	if err := db.Table("room_card_product_list").Where("id = ?", req.ID).Delete(req).Error; err != nil {
		log.Errorf("OnPayListDeleteRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	if err1 == nil {
		config := c.MustGet("config").(*config.Config)
		args := make(map[string]interface{})
		before, _ := json.Marshal(beforeContent)
		args["before_content"] = string(before)
		after, _ := json.Marshal(models.RoomCardProductList{})
		args["after_content"] = string(after)
		args["game_id"] = beforeContent.GameID
		args["module_id"] = 1
		args["operate_id"] = userID
		args["op_type"] = 3
		err := utils.AddConfigLog(config.ConfigUpdateLogAPI, args)
		if err != nil {
			log.Errorf("OnPayListDeleteRequest-AddConfigLog: err: %v\n", err)
		}
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type Req struct {
	GameID     int `json:"game_id" binding:"required"`
	GameAreaID int `json:"game_area_id"`
}

func GameAreaID(c *gin.Context) {
	req := &Req{}
	if err := c.Bind(req); err != nil {
		log.Errorf("err:%v", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	log.Printf("game_id:%v", req.GameID)
	log.Printf("game_area_id:%v", req.GameAreaID)

	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type PayListCopyReq struct {
	IDList   []int `json:"id_list" binding:"required"`
	ToGameID int   `json:"to_game_id" binding:"required"`
}

func OnPayListCopyRequest(c *gin.Context) {
	req := &PayListCopyReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnPayListCopyRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnPayListCopyRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "参数不匹配" + err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	//删除目标平台所有配置
	err := models.DeleteRoomcardProductListByToGameID(db, req.ToGameID)
	if err != nil {
		log.Errorf("OnPayListCopyRequest-err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "删除目标平台配置出错" + err.Error()})
		return
	}
	//找出这些ID的套餐信息
	data, err := models.GetProductConfigByIDList(db, req.IDList)
	if err != nil {
		log.Errorf("OnPayListCopyRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	newProductList := make([]models.RoomCardProductList, 0)

	for _, v := range data {
		//去除掉卡包和活动房卡套餐
		//if v.GameAreaID != "-1" {
		//	continue
		//}
		//修改列表的平台 添加时间
		temp := v
		temp.GameID = int32(req.ToGameID)
		addTime, _ := v.AddTime.MarshalJSON()
		temp.AddTimeStr = string(addTime)
		temp.ID = 0
		endTime, _ := v.EndTime.MarshalJSON()
		temp.EndTimeStr = string(endTime)
		newProductList = append(newProductList, temp)
	}
	//批量插入套餐列表
	err = models.InsertIntoRoomcardProductList(db, newProductList)
	if err != nil {
		log.Errorf("OnPayListCopyRequest-err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "复制失败," + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "复制成功"})
	return
}

// ValidateBlackWhiteList 校验黑白名单是否存在相同的用户
func ValidateBlackWhiteList(whiteList []int, blackList []int) error {
	nums := make([]int, 0)
	nums = append(append(nums, whiteList...), blackList...)
	temp := map[int]struct{}{}
	for _, item := range nums {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
		} else {
			return errors.New("黑白名单中存在相同的用户")
		}
	}
	return nil
}

// removeDuplicate 数组去重
func removeDuplicate(nums []int) []int {
	res := make([]int, 0)
	temp := map[int]struct{}{}
	for _, item := range nums {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			res = append(res, item)
		}
	}
	return res
}

// changeToPqArray
func changeToPqArray(nums []int) pq.Int64Array {
	var res = make(pq.Int64Array, 0)
	for _, item := range nums {
		res = append(res, int64(item))
	}
	return res
}
