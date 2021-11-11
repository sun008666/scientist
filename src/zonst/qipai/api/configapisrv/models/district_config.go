package models

import (
	"encoding/json"
	"zonst/logging"

	"github.com/jinzhu/gorm"
)

type GameDistrictConfig struct {
	ID            int32    `json:"id" gorm:"column:id"`
	GameID        int32    `json:"game_id" gorm:"column:game_id"`
	Province      string   `json:"province" gorm:"column:province"`
	City          string   `json:"city" gorm:"column:city"`
	District      string   `json:"-" gorm:"column:district"`
	DistrictArray []string `json:"district_array"`
	WXName        string   `json:"wx_name" gorm:"column:wx_name"`
	WxAccount     string   `json:"wx_account" gorm:"column:wx_account"`
	TitleContent  string   `json:"title_content" gorm:"column:title_content"`
}

type AddGameDistrictConfig struct {
	ID           int32  `json:"id" gorm:"column:id"`
	GameID       int32  `json:"game_id" gorm:"column:game_id"`
	Province     string `json:"province" gorm:"column:province"`
	City         string `json:"city" gorm:"column:city"`
	District     string `json:"district" gorm:"column:district"`
	WXName       string `json:"wx_name" gorm:"column:wx_name"`
	WxAccount    string `json:"wx_account" gorm:"column:wx_account"`
	TitleContent string `json:"title_content" gorm:"column:title_content"`
}

func GetDistrictConfigList(qipaiDB *gorm.DB, gameID int32, province, city string) (data []*GameDistrictConfig, err error) {
	sql := qipaiDB.Debug().Table("game_district_config").Order("game_id asc")
	if gameID != 0 {
		sql = sql.Where("game_id=?", gameID)
	}
	if province != "" {
		sql = sql.Where("province=?", province)
	}
	if city != "" {
		sql = sql.Where("city=?", city)
	}
	if err := sql.Find(&data).Error; err != nil {
		return nil, err
	}
	for _, v := range data {
		err := json.Unmarshal([]byte(v.District), &v.DistrictArray)
		if err != nil {
			logging.Errorf("GetDistrictConfigList-Unmarshal:%v", err)
			return nil, err
		}

	}
	return data, nil
}
func GetDistrictConfigByID(qipaiDB *gorm.DB, id int32) (data GameDistrictConfig, err error) {
	if err := qipaiDB.Debug().Table("game_district_config").Where("id=?", id).First(&data).Error; err != nil {
		return data, err
	}
	err = json.Unmarshal([]byte(data.District), &data.DistrictArray)
	if err != nil {
		logging.Errorf("GetDistrictConfigByID-Unmarshal:%v", err)
		return data, err
	}

	return data, nil
}
func (c *GameDistrictConfig) Val(req []string) bool {
	if len(c.DistrictArray) != len(req) {
		return true
	}
	for _, va := range req {
		flag := false
		for _, v := range c.DistrictArray {
			if va == v {
				flag = true
				break
			}
		}
		if !flag {
			return true
		}

	}

	return false
}
func DeleteGameDistrictConfig(qipaiDB *gorm.DB, id int32) error {
	del := AddGameDistrictConfig{}
	if err := qipaiDB.Debug().Table("game_district_config").Where("id=?", id).Delete(&del).Error; err != nil {
		return err
	}
	return nil
}
