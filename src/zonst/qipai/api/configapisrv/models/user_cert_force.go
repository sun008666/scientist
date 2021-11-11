package models

import (
	"fmt"
	"time"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

// UserCertForceStatus 用户强制实名认证状态
type UserCertForceStatus struct {
	GameID int                       `json:"game_id"`
	Status utils.UserCertForceStatus `json:"status"`
	Scope  utils.UserCertForceScope  `json:"scope"`
}

// TableName name
func (u UserCertForceStatus) TableName() string {
	return "user_cert_force_status"
}

// UpdateUserCertForceStatus 更新用户强制实名认证状态
func UpdateUserCertForceStatus(qipaidb *gorm.DB, gameID int, status utils.UserCertForceStatus, scope utils.UserCertForceScope) error {
	now := time.Now()
	err := qipaidb.Debug().Set(`gorm:insert_option`,
		fmt.Sprintf("ON CONFLICT(game_id) do update set status=%d,scope=%d,create_date=now(),create_time=%d",
			status, scope, now.Unix())).Create(UserCertForceStatus{
		GameID: gameID,
		Status: status,
		Scope:  scope,
	}).Error
	return err
}

// FindUserCertForceStatus 查找用户强制实名认证状态
func FindUserCertForceStatus(qipaidb *gorm.DB, gameIDS []int) (status []UserCertForceStatus, err error) {
	var db = qipaidb.Debug()
	if len(gameIDS) != 0 {
		db = db.Where(`game_id in (?)`, gameIDS)
	}
	err = db.Order(`game_id`).Find(&status).Error
	return
}

// FindUserCertForceStatusByGameID 查找用户强制实名认证状态
func FindUserCertForceStatusByGameID(qipaidb *gorm.DB, gameID int) (status UserCertForceStatus, err error) {
	err = qipaidb.Debug().Where(`game_id =  ?`, gameID).Find(&status).Error
	return
}

// UserCertForceStatusUpdateLog log
type UserCertForceStatusUpdateLog struct {
	GameID       int
	OpUserId     int
	OpUsername   string
	BeforeStatus utils.UserCertForceStatus `gorm:"column:before_statu"`
	AfterStatus  utils.UserCertForceStatus
	BeforeScope  utils.UserCertForceScope
	AfterScope   utils.UserCertForceScope
}

// TableName name
func (log UserCertForceStatusUpdateLog) TableName() string {
	return "user_cert_force_status_update_log"
}

// CreateUserCertForceStatusUpdateLog create log
func CreateUserCertForceStatusUpdateLog(logdb *sqlx.DB, log UserCertForceStatusUpdateLog) error {
	//err := logdb.Debug().Create(&log).Error
	_, err := logdb.Exec("insert into user_cert_force_status_update_log"+
		"(game_id,op_user_id,op_username,before_statu,after_status,before_scope,after_scope)"+
		"values($1,$2,$3,$4,$5,$6,$7)", log.GameID, log.OpUserId, log.OpUsername,
		log.BeforeStatus, log.AfterStatus, log.BeforeScope, log.AfterScope)
	return err
}
