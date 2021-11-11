package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/models"
)

// DeleteCache 删除缓存
func DeleteCache(pool *redis.Pool, gameID int, gameAreaID pq.Int64Array) (err error) {
	conn := pool.Get()
	c := cache.ThreeDLobbyGameUpdate{GameAreaID: gameAreaID}
	err = c.DeleteCache(conn, gameID)
	if err != nil {
		return err
	}
	return nil
}

// On3DGameConfigUpdateRequest 将前端获取的数据存入数据库里面
func On3DGameConfigUpdateRequest(configDB *gorm.DB, pool *redis.Pool, gameID int, gameAreaID pq.Int64Array) error {
	a := models.Game3DConfigUpdate{GameID: gameID, GameAreaID: gameAreaID}
	data, err := a.FindByAreaGameID(configDB, gameID)
	if data.GameID != 0 {
		err = a.UpdateGameArea(configDB)
		if err != nil {
			return err
		}
	} else {
		err = a.InsertGameArea(configDB)
		if err != nil {
			return err
		}
	}
	err = DeleteCache(pool, gameID, gameAreaID)
	if err != nil {
		return err
	}
	return nil
}
