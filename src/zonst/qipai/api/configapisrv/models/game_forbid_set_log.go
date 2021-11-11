package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// 开启/关闭游戏平台的日志

type GameForbidSetLog struct {
	GameID       int32  `json:"game_id" db:"game_id"`
	UserID       int32  `json:"user_id" db:"user_id"`
	Username     string `json:"username" db:"username"`
	BeforeStatus int32  `json:"before_forbid" db:"before_forbid"`
	AfterStatus  int32  `json:"after_forbid" db:"after_forbid"`
	CreateTime   int64  `json:"create_time" db:"create_time"`
}

// 插入
func InsertGameForbidSetLog(db *sqlx.DB, gameID, userID, beforeStatus, afterStatus int32, username string) error {
	now := time.Now()
	tpl := "INSERT INTO game_forbid_set_log(game_id,user_id,username,before_forbid,after_forbid,create_date,create_time) VALUES($1,$2,$3,$4,$5,$6,$7);"
	_, err := db.Unsafe().Exec(tpl, gameID, userID, username, beforeStatus, afterStatus, now.Format("2006-01-02"), now.Unix())
	return err
}

// // 查询
// func FindGameForbidSetLog(db *sqlx.DB, gameID, page, pageSize int32) ([]*GameForbidSetLog, error) {
// 	var s []*GameForbidSetLog
// 	tpl := "SELECT game_id,user_id,username,before_status,after_status,create_time FROM game_forbid_set_log WHERE game_id=$1 ORDER BY id DESC LIMIT $2 OFFSET $3;"
// 	err := db.Unsafe().Select(&s, tpl, gameID, pageSize, pageSize*page)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return s, nil
// }
//
// // 查询所有日志
// func FindAllGameForbidSetLog(db *sqlx.DB, page, pageSize int32) ([]*GameForbidSetLog, error) {
// 	var s []*GameForbidSetLog
// 	tpl := "SELECT game_id,user_id,username,before_status,after_status,create_time FROM game_forbid_set_log ORDER BY id DESC LIMIT $1 OFFSET $2;"
// 	err := db.Unsafe().Select(&s, tpl, pageSize, pageSize*page)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return s, nil
// }
