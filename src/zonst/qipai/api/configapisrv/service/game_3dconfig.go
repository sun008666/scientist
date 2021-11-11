package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/models"
)

// On3DConfigListConfigID 获得子游戏一二的ID
func On3DConfigListConfigID(configDB *gorm.DB, pool *redis.Pool, gameID int) (AreaIDCache cache.ThreeDLobbyGameUpdate, err error) {
	conn := pool.Get()
	var c cache.ThreeDLobbyGameUpdate
	areaCache, exists, _ := c.FindByGameID(conn, gameID)
	if exists {
		return areaCache, nil
	}
	//从数据库里面查询数据
	var data models.Game3DConfigUpdate
	userInfoModel, err := data.FindByAreaGameID(configDB, gameID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return cache.ThreeDLobbyGameUpdate{}, err
	}
	//无论数据库是否存在都需要往缓存里面存
	c = cache.ThreeDLobbyGameUpdate{GameAreaID: userInfoModel.GameAreaID}
	err = c.CreateList(conn, gameID)
	if err != nil {
		return c, err
	}
	if userInfoModel.GameAreaID != nil {
		return cache.ThreeDLobbyGameUpdate{userInfoModel.GameAreaID}, err
	}
	return cache.ThreeDLobbyGameUpdate{pq.Int64Array{-2, -2}}, nil
}
