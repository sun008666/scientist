package models

import (
	"encoding/json"
	"strconv"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// RoomCardProductAddRequest 充值列表
type RoomCardProductAddRequest struct {
	GameID              int           `gorm:"column:game_id" json:"game_id"`
	GameAreaID          string        `gorm:"column:game_area_id" json:"game_area_id"`
	CardNum             int           `gorm:"column:card_num" json:"card_num"`
	Status              int           `gorm:"column:status" json:"status"`
	Price               int           `gorm:"column:price" json:"price"`
	AddTime             string        `gorm:"column:add_time" json:"add_time"`
	Category            string        `gorm:"column:category" json:"category"`
	PresentRoomCard     int           `gorm:"column:present_room_card" json:"present_room_card"`
	EndTime             string        `gorm:"column:end_time" json:"end_time"`
	Represent           string        `gorm:"column:represent" json:"represent"`
	AppleKey            string        `gorm:"column:apple_key" json:"apple_key"`
	IsWhite             bool          `gorm:"column:is_white" json:"is_white"`
	IsBlack             bool          `gorm:"column:is_black" json:"is_black"`
	WhiteList           pq.Int64Array `gorm:"column:white_list" json:"white_list"`
	BlackList           pq.Int64Array `gorm:"column:black_list" json:"black_list"`
	CardPackID          int32         `gorm:"column:card_pack_id" json:"card_pack_id"` //卡包ID
	IsSeniorAgent       int32         `gorm:"column:is_senior_agent" json:"is_senior_agent"`
	PurchaseLimitNumber int32         `gorm:"column:purchase_limit_number" json:"purchase_limit_number"`
}

// RoomCardProductList 充值列表
type RoomCardProductList struct {
	ID                  int             `gorm:"column:id" json:"id"`
	GameID              int32           `gorm:"column:game_id" json:"game_id"`
	GameAreaID          string          `gorm:"column:game_area_id" json:"game_area_id"`
	CardNum             int             `gorm:"column:card_num" json:"card_num"`
	Status              int             `gorm:"column:status" json:"status"`
	Price               int             `gorm:"column:price" json:"price"`
	AddTime             utils.LocalTime `gorm:"column:add_time" json:"add_time"`
	Category            string          `gorm:"column:category" json:"category"`
	PresentRoomCard     int             `gorm:"column:present_room_card" json:"present_room_card"`
	EndTime             utils.LocalTime `gorm:"column:end_time" json:"end_time"`
	Represent           string          `gorm:"column:represent" json:"represent"`
	AppleKey            string          `gorm:"column:apple_key" json:"apple_key"`
	IsWhite             bool            `gorm:"column:is_white" json:"is_white"`
	IsBlack             bool            `gorm:"column:is_black" json:"is_black"`
	WhiteList           pq.Int64Array   `gorm:"column:white_list" json:"white_list"`
	BlackList           pq.Int64Array   `gorm:"column:black_list" json:"black_list"`
	CardPackID          int32           `gorm:"card_pack_id" json:"card_pack_id"`
	IsSeniorAgent       int32           `gorm:"is_senior_agent" json:"is_senior_agent"`
	AddTimeStr          string          `gorm:"-" json:"-"`
	EndTimeStr          string          `gorm:"-" json:"-"`
	PurchaseLimitNumber int32           `gorm:"column:purchase_limit_number" json:"purchase_limit_number"`
}

// TableName 表名
func (a *RoomCardProductList) TableName() string {
	return "room_card_product_list"
}

// RoomCardProductUpdateRequest 充值列表
type RoomCardProductUpdateRequest struct {
	ID              int    `gorm:"column:id" json:"id"`
	GameID          int    `gorm:"column:game_id" json:"game_id"`
	CardNum         int    `gorm:"column:card_num" json:"card_num"`
	Status          int    `gorm:"column:status" json:"status"`
	Price           int    `gorm:"column:price" json:"price"`
	AddTime         string `gorm:"column:add_time" json:"add_time"`
	Category        string `gorm:"column:category" json:"category"`
	PresentRoomCard int    `gorm:"column:present_room_card" json:"present_room_card"`
	EndTime         string `gorm:"column:end_time" json:"end_time"`
	Represent       string `gorm:"column:represent" json:"represent"`
	AppleKey        string `gorm:"column:apple_key" json:"apple_key"`
	IsSeniorAgent   int32  `gorm:"column:is_senior_agent" json:"is_senior_agent"`
}

func GetProductList(qipaiDB *gorm.DB, gameID, gameAreaID int, sourceType string, isAppleProduct bool) (data []*RoomCardProductList, err error) {
	sql := qipaiDB.Debug().Table("room_card_product_list").Where("game_id = ?", gameID)
	// 活动房卡有子游戏选择
	if gameAreaID != -2 {
		sql = sql.Where("game_area_id=?", strconv.Itoa(int(gameAreaID)))
	}
	// 套餐类型
	if sourceType != "" {
		sql = sql.Where("category=?", sourceType)
	}
	// 是否是苹果内购套餐
	if isAppleProduct != false {
		sql = sql.Where("apple_key!=?", "")
	}
	if err := sql.Order("status desc, category,game_area_id,card_pack_id, price").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func GetLastProductConfig(qipaiDB *gorm.DB, gameID int, gameAreaID, sourceType string) (string, error) {
	data := RoomCardProductList{}
	where := map[string]interface{}{"game_id": gameID, "game_area_id": gameAreaID, "category": sourceType}
	if err := qipaiDB.Debug().Last(&data, where).Error; err != nil {
		return "", err
	}
	result, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func GetProductConfigByID(qipaiDB *gorm.DB, productID int) (RoomCardProductList, error) {
	data := RoomCardProductList{}
	if err := qipaiDB.Debug().First(&data, "id=?", productID).Error; err != nil {
		return data, err
	}
	return data, nil
}

func GetProductConfigByIDList(qipaiDB *gorm.DB, productID []int) ([]RoomCardProductList, error) {
	data := make([]RoomCardProductList, 0)
	if err := qipaiDB.Debug().Table("room_card_product_list").Where("id in(?)", productID).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
func DeleteRoomcardProductListByToGameID(qipaiDB *gorm.DB, toGameID int) error {
	tx := qipaiDB.Begin()
	sql := "delete from room_card_product_list where game_id=?"
	if err := tx.Exec(sql, toGameID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func BatchDeleteRoomcardProductListByToGameID(qipaiDB *gorm.DB, toGameIDs []int32) error {
	tx := qipaiDB.Begin()
	sql := "delete from room_card_product_list where game_id in (?)"
	if err := tx.Exec(sql, toGameIDs).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func InsertIntoRoomcardProductList(qipaiDB *gorm.DB, data []RoomCardProductList) error {
	tx := qipaiDB.Begin()
	insertIntoSql := "insert into room_card_product_list(game_id,game_area_id,card_num,status,price,add_time,category,present_room_card,end_time,represent,apple_key,is_white," +
		"is_black,white_list,black_list,card_pack_id,is_senior_agent,purchase_limit_number) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	for _, v := range data {
		if err := tx.Exec(insertIntoSql, v.GameID, v.GameAreaID, v.CardNum, v.Status, v.Price, v.AddTimeStr, v.Category, v.PresentRoomCard,
			v.EndTimeStr, v.Represent, v.AppleKey, v.IsWhite, v.IsBlack, v.WhiteList, v.BlackList, v.CardPackID, v.IsSeniorAgent, v.PurchaseLimitNumber).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
