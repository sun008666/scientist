package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/models"
)

// GanzhouGameConfigurationRequest 赣州游戏配置保存设置接口
func GanzhouGameConfiguratioRequest(configDB *gorm.DB, pool *redis.Pool, gameID int, gameAreaID int, id int, subGameID int) error {
	conn := pool.Get()
	a := models.GanzhouGameConfiguration{GameID: gameID, GameAreaID: gameAreaID, ID: id, SubGameID: subGameID}
	if id != 0 {
		err := a.UpdateGameArea(configDB)
		if err != nil {
			return err
		}
		e := cache.DeleteCache(conn)
		if e != nil {
			return e
		}
	} else {
		err := a.InsertGameArea(configDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// GanzhouGameConfigurationFindRequest 赣州游戏配置查询配置是否存在接口
func GanzhouGameConfigurationFindRequest(configDB *gorm.DB, gameID int, gameAreaID int, subGameID int) bool {
	err := models.FindByGameIDAndGameAreaID(configDB, gameID, gameAreaID, subGameID)
	if err == nil {
		return true
	}
	return false
}
