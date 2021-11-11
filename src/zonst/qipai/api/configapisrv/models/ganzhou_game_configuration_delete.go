package models

import "github.com/jinzhu/gorm"

type GanzhouGameConfigurationDelete struct {
	ID        int `json:"id" gorm:"column:id"`
	SubGameID int `json:"sub_game_id" gorm:"column:sub_game_id"`
}

// DeleteConfig 查询赣州游戏配置
func (g *GanzhouGameConfigurationDelete) DeleteConfig(configDB *gorm.DB) error {
	return configDB.Debug().Table("sub_game_game_area").
		Where("id=? and sub_game_id = ?", g.ID, g.SubGameID).Delete(g).Error
}
