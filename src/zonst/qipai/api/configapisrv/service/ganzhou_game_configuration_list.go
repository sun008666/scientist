package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/models"
)

func GetGanzhouGameConfigurationList(configDB *gorm.DB, pool *redis.Pool) ([]cache.GanzhouGameConfiguration, error) {
	conn := pool.Get()
	var ganzhouCache cache.GanzhouGameConfiguration
	ganzhouCacheData, exists, _ := ganzhouCache.FindGanConfiguration(conn)
	if exists {
		return ganzhouCacheData, nil
	}
	//从数据库里面查询数据
	userInfoModel, err := models.FindConfig(configDB)
	if err != nil && err != gorm.ErrRecordNotFound {
		return []cache.GanzhouGameConfiguration{}, err
	}
	cacheList := make([]cache.GanzhouGameConfiguration, 0)
	for _, data := range userInfoModel {
		cacheList = append(cacheList, cache.GanzhouGameConfiguration{ID: data.ID, GameID: data.GameID, GameAreaID: data.GameAreaID, SubGameID: data.SubGameID})
	}
	//无论数据库里面存在都需要放入缓存里面
	e := ganzhouCache.CreateList(conn, cacheList)
	if e != nil {
		return cacheList, e
	}
	return cacheList, nil
}
