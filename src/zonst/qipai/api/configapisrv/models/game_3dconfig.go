package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Game3DConfigUpdate struct {
	GameID     int           `json:"game_id" gorm:"column:game_id"`
	GameAreaID pq.Int64Array `json:"game_area_id" gorm:"column:game_area_id"`
}

func (*Game3DConfigUpdate) FindByAreaGameID(configDB *gorm.DB, gameID int) (data Game3DConfigUpdate, err error) {
	sql := configDB.Debug().Table("u3d_default_game_area").Select("game_area_id,game_id")
	if gameID != 0 {
		sql = sql.Where("game_id=?", gameID)
	}
	if err := sql.Scan(&data).Error; err != nil {
		return Game3DConfigUpdate{}, err
	}
	return data, nil
}
func (a *Game3DConfigUpdate) InsertGameArea(configDB *gorm.DB) error {
	err := configDB.Debug().Exec("insert into u3d_default_game_area(game_id,game_area_id) values(?, ?)", a.GameID, a.GameAreaID).Error
	if err != nil {
		return err
	}
	return nil
}
func (a *Game3DConfigUpdate) UpdateGameArea(configDB *gorm.DB) error {
	err := configDB.Debug().Exec("update u3d_default_game_area set game_area_id=? where game_id=?", a.GameAreaID, a.GameID).Error
	if err != nil {
		return err
	}
	return nil
}
