package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"zonst/qipai/api/configapisrv/cache"
	"zonst/qipai/api/configapisrv/models"
)

func DeterMineUserType(pointdb *gorm.DB, pool *redis.Pool, userID int) (cache.DeterMineUserTypeFind, error) {
	conn := pool.Get()
	var d cache.DeterMineUserTypeFind
	// 从缓存里面去取数据
	userType, exist, _ := d.FindUserTypeByUserID(conn, userID)
	if exist {
		return userType, nil
	}
	var data models.AreaUserType
	userTypeModel, e := data.FindUserType(pointdb, userID)
	if e != nil && e != gorm.ErrRecordNotFound {
		return userType, e
	}
	// 无论数据库是否存在都要去放入缓存
	c := cache.DeterMineUserTypeFind{BindGameID: userTypeModel.BindGameID}
	e = c.CreatList(conn, userID)
	if e != nil {
		return userType, e
	}
	return cache.DeterMineUserTypeFind{BindGameID: userTypeModel.BindGameID}, nil
}
