package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/models"
)

// GanzhouGameConfigurationDeleteRequest 赣州游戏配置删除设置接口
func GanzhouGameConfiguratioDeleteRequest(configDB *gorm.DB, pool *redis.Pool, id int, subGameID int) error {
	conn := pool.Get()
	g := models.GanzhouGameConfigurationDelete{ID: id, SubGameID: subGameID}
	err := g.DeleteConfig(configDB)
	if err != nil {
		return err
	}
	e := cache.DeleteCache(conn)
	if e != nil {
		return e
	}
	return nil
}
