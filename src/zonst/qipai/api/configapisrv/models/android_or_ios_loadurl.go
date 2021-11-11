package models

import "github.com/jinzhu/gorm"

// AndroidOrIosLoadUrlRep 游戏配置列表
type AndroidOrIosLoadUrlRep struct {
	GameID     int    `json:"game_id" gorm:"column:game_id"`
	AndroidUrl string `json:"android_url" gorm:"column:android_url"`
	IosUrl     string `json:"ios_url" gorm:"column:ios_url"`
}

// FindAndroidOrIosLoadUrl 从数据库中查询查询平台安卓包或ios包下载地址
func (*AndroidOrIosLoadUrlRep) FindAndroidOrIosLoadUrl(qipaiDB *gorm.DB, gameID int) (data AndroidOrIosLoadUrlRep, err error) {
	sql := qipaiDB.Debug().Table("html_page").
		Select("android_url,ios_url,game_id").Where("game_id=?", gameID)
	if err := sql.Find(&data).Error; err != nil {
		return AndroidOrIosLoadUrlRep{}, err
	}
	return data, nil
}
