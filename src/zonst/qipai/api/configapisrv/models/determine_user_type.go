package models

import (
	"github.com/jinzhu/gorm"
)

type AreaUserType struct {
	BindGameID int `gorm:"column:bind_game_id" jsonx:"bind_game_id"`
}

// FindUserType 查询用户所属地区
func (a *AreaUserType) FindUserType(pointdb *gorm.DB, userID int) (AreaUserType, error) {
	var areaUserType AreaUserType
	if e := pointdb.Debug().Table("area_user_type").Select("bind_game_id").Where("user_id = ?", userID).Find(&areaUserType).Error; e != nil {
		return areaUserType, e
	}
	return areaUserType, nil
}
