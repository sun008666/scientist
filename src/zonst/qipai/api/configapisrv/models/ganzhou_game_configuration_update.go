package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type GanzhouGameConfiguration struct {
	ID         int `json:"id" gorm:"column:id"`
	GameID     int `json:"game_id" gorm:"column:game_id"`
	GameAreaID int `json:"game_area_id" gorm:"column:game_area_id"`
	SubGameID  int `json:"sub_game_id" gorm:"column:sub_game_id"`
}

// FindConfig 查询赣州游戏配置
func FindConfig(configDB *gorm.DB) (data []GanzhouGameConfiguration, err error) {
	sql := configDB.Debug().Table("sub_game_game_area").Select("game_id,game_area_id,id,sub_game_id").Order("game_id,game_area_id").Where("sub_game_id  = 17")
	if err := sql.Scan(&data).Error; err != nil {
		return []GanzhouGameConfiguration{}, err
	}
	return data, nil
}

// FindByGameIDAndGameAreaID 根据平台ID和子游戏ID查询赣州游戏配置是否存在
func FindByGameIDAndGameAreaID(configDB *gorm.DB, gameID int, gameAreaID int, subGameID int) (err error) {
	var data GanzhouGameConfiguration
	sql := configDB.Debug().Table("sub_game_game_area").Select("game_area_id,game_id,id").
		Where("game_id=? and game_area_id=? and sub_game_id = ?", gameID, gameAreaID, subGameID)
	err = sql.Find(&data).Error
	if err == gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

// InsertGameArea 插入赣州游戏配置
func (a *GanzhouGameConfiguration) InsertGameArea(configDB *gorm.DB) error {
	err := configDB.Debug().Exec("insert into sub_game_game_area(game_id,game_area_id,sub_game_id) values(?, ?, ?)", a.GameID, a.GameAreaID, a.SubGameID).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateGameArea 修改赣州游戏配置
func (a *GanzhouGameConfiguration) UpdateGameArea(configDB *gorm.DB) error {
	UpdateDate := time.Now()
	UpdateTime := time.Now().Unix()
	err := configDB.Debug().Exec("update sub_game_game_area set game_area_id=?,game_id=?,update_date=?,update_time=? where id=? and sub_game_id = ?", a.GameAreaID, a.GameID, UpdateDate, UpdateTime, a.ID, a.SubGameID).Error
	if err != nil {
		return err
	}
	return nil
}
