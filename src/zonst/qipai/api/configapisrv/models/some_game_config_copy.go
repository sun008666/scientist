package models

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//任务配置列表
type TaskConfig struct {
	GameID           int32          `json:"game_id" db:"game_id"`
	TaskID           int32          `json:"task_id" db:"task_id"`
	UserGroup        pq.Int64Array  `json:"user_group" db:"user_group"`
	TaskType         int32          `json:"task_type" db:"task_type"`
	TaskName         string         `json:"task_name" db:"task_name"`
	TaskModule       pq.Int64Array  `json:"task_module" db:"task_module"`
	TaskGameAreaID   pq.Int64Array  `json:"task_game_area_id" db:"task_game_area_id"`
	TaskGameAreaName pq.StringArray `json:"task_game_area_name" db:"task_game_area_name"`
	StartTime        int32          `json:"start_time" db:"start_time"`
	EndTime          int32          `json:"end_time" db:"end_time"`
	TaskTarget       string         `json:"task_target" db:"task_target"`
	TaskFinished     int32          `json:"task_finished" db:"task_finished"`
	PreTaskID        int32          `json:"pre_task_id" db:"pre_task_id"`
	AwardType        int32          `json:"award_type" db:"award_type"`
	AwardMax         int32          `json:"award_max" db:"award_max"`
	AwardMin         int32          `json:"award_min" db:"award_min"`
	HaoPaiIDs        pq.Int64Array  `json:"hao_pai_ids" db:"hao_pai_ids"`
	Strategy         int32          `json:"strategy" db:"strategy"`
	HaoPaiName       pq.StringArray `json:"hao_pai_name" db:"hao_pai_name"`
	BeilvStartTime   int32          `json:"beilv_start_time" db:"beilv_start_time"`
	BeilvEndTime     int32          `json:"beilv_end_time" db:"beilv_end_time"`
	Beilv            int32          `json:"beilv" db:"beilv"`
}

// GetCopyTaskConfigList 获取需要被复制任务配置
func GetCopyTaskConfigList(taskdb *gorm.DB, gameID int32) (result []TaskConfig, err error) {
	psql := taskdb.Table("task_config").Where("game_id = ?", gameID).Debug()
	err = psql.Order("task_id").Find(&result).Error
	return
}

func DeleteCopyTaskConfigByDestGameID(taskdb *sqlx.DB, gameID int32) (err error) {
	var tx *sql.Tx
	if tx, err = taskdb.Begin(); err != nil {
		return
	}
	sql := "delete from task_config where game_id=$1"
	if _, err = tx.Exec(sql, gameID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func BatchDeleteCopyTaskConfigByDestGameID(taskdb *sqlx.DB, gameIDs []int32) (err error) {
	var tx *sql.Tx
	if tx, err = taskdb.Begin(); err != nil {
		return
	}
	query, args, err := sqlx.In("delete from task_config where game_id in (?)", gameIDs)
	if err != nil {
		return err
	}
	query = taskdb.Rebind(query)
	if _, err = tx.Exec(query, args...); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// BashCopyTaskConfig 批量复制任务配置到指定的游戏平台
func BashCopyTaskConfig(taskdb *sqlx.DB, gameID int32, taskConfigList []TaskConfig) (err error) {
	var tx *sql.Tx
	if tx, err = taskdb.Begin(); err != nil {
		return
	}
	tpl := `insert into task_config(game_id,task_type,task_name,task_module,task_game_area_id,start_time,end_time,task_target,task_finished,pre_task_id,award_type,award_max,award_min,task_game_area_name,user_group,hao_pai_ids,strategy,hao_pai_name,beilv_start_time,beilv_end_time,beilv) 
			values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21);`
	var taskConfigStmt *sql.Stmt
	if taskConfigStmt, err = tx.Prepare(tpl); err != nil {
		tx.Rollback()
		return
	}
	for _, val := range taskConfigList {
		if _, err = taskConfigStmt.Exec(gameID, val.TaskType, val.TaskName, val.TaskModule, val.TaskGameAreaID, val.StartTime, val.EndTime, val.TaskTarget, val.TaskFinished, 0, val.AwardType, val.AwardMax, val.AwardMin, val.TaskGameAreaName, val.UserGroup, val.HaoPaiIDs, val.Strategy, val.HaoPaiName, val.BeilvStartTime, val.EndTime, val.Beilv); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = taskConfigStmt.Close(); err != nil {
		tx.Rollback()
		return fmt.Errorf("关掉stmt出错, err:%v", err)
	}
	return tx.Commit()
}

func BatchCopyTaskConfigToGameIdList(taskdb *sqlx.DB, gameIDs []int32, taskConfigList []TaskConfig) (err error) {
	var tx *sql.Tx
	if tx, err = taskdb.Begin(); err != nil {
		return
	}
	tpl := `insert into task_config(game_id,task_type,task_name,task_module,task_game_area_id,start_time,end_time,task_target,task_finished,pre_task_id,award_type,award_max,award_min,task_game_area_name,user_group,hao_pai_ids,strategy,hao_pai_name,beilv_start_time,beilv_end_time,beilv) 
			values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21);`
	var taskConfigStmt *sql.Stmt
	if taskConfigStmt, err = tx.Prepare(tpl); err != nil {
		tx.Rollback()
		return
	}
	for _, val := range taskConfigList {
		for _, gameID := range gameIDs {
			if _, err = taskConfigStmt.Exec(gameID, val.TaskType, val.TaskName, val.TaskModule, val.TaskGameAreaID, val.StartTime, val.EndTime, val.TaskTarget, val.TaskFinished, 0, val.AwardType, val.AwardMax, val.AwardMin, val.TaskGameAreaName, val.UserGroup, val.HaoPaiIDs, val.Strategy, val.HaoPaiName, val.BeilvStartTime, val.EndTime, val.Beilv); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if err = taskConfigStmt.Close(); err != nil {
		tx.Rollback()
		return fmt.Errorf("关掉stmt出错, err:%v", err)
	}
	return tx.Commit()
}

func GetProductConfig(userDB *gorm.DB, srcGameID int) ([]RoomCardProductList, error) {
	data := make([]RoomCardProductList, 0)
	if err := userDB.Debug().Table("room_card_product_list").Where("game_id=?", srcGameID).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
