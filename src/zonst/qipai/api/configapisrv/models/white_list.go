package models

import (
	"github.com/fwhezfwhez/errorx"
	"github.com/garyburd/redigo/redis"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
)

type WhiteList struct {
	GameID  int32  `json:"game_id" gorm:"column:game_id"`
	UserID  int32  `json:"user_id" gorm:"column:user_id"`
	Remark  string `json:"remark" gorm:"column:remark"`
	LogDate string `json:"log_date" gorm:"column:log_date"`
	LogTime int32  `json:"log_time" gorm:"column:log_time"`
}

func GetWhiteList(qipaidb *gorm.DB, gameID, userID, pageID, pageCount int32) (int32, []WhiteList, error) {
	offset := (pageID - 1) * pageCount
	if offset < 0 {
		offset = 0
	}
	var total int32
	data := make([]WhiteList, 0)
	sql := qipaidb.Debug().Table("popup_white_list").Where("game_id=?", gameID)
	if userID != 0 {
		sql = sql.Where("user_id=?", userID)
	}
	if err := sql.Count(&total).Order("log_time desc").Limit(pageCount).Offset(offset).Find(&data).Error; err != nil {
		return 0, nil, err
	}

	return total, data, nil

}
func DeleteWhiteList(qipaidb *gorm.DB, gameID, userID int32) error {
	w := WhiteList{}
	if err := qipaidb.Debug().Table("popup_white_list").Where("game_id=? and user_id=?", gameID, userID).Unscoped().Delete(&w).Error; err != nil {
		return err
	}

	return nil
}

// SetPopupUserWhiteList 设置弹窗用户白名单
func SetPopupUserWhiteList(pool *redis.Pool, userID int32) error {
	conn := pool.Get()
	defer conn.Close()
	tblName := "popup:white:list"
	_, err := conn.Do("SADD", tblName, userID)
	if err != nil {
		log.Errorf("SetUserForbidRedis，添加用户白名单用户:%v", err)
		return err
	}
	return nil

}

// RemovePopupUserWhitelist 移除弹窗用户白名单
func RemovePopupUserWhitelist(pool *redis.Pool, userID int32) error {
	//pool := c.MustGet("whiteListCache").(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()
	tblName := "popup:white:list"
	_, err := conn.Do("SREM", tblName, userID)
	if err != nil {
		log.Errorf("RemoveUserRedis,移除用户白名单:%v", err)
		return err
	}
	return nil

}

// IsWhiteUser 判断是否是白名单
func IsWhiteUser(pool *redis.Pool, userID int32) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	tblName := "popup:white:list"
	flag := false
	flag, err := redis.Bool(conn.Do("SISMEMBER", tblName, userID))
	if err != nil && err != redis.ErrNil {
		return false, errorx.Wrap(err)
	}
	return flag, nil
}
