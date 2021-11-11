package service

import (
	"github.com/jinzhu/gorm"
	"zonst/qipai/api/configapisrv/models"
)

func FindAndroidOrIosLoadUrl(qipaiDB *gorm.DB, gameID int) (models.AndroidOrIosLoadUrlRep, error) {
	var androidOrIosLoadUrl models.AndroidOrIosLoadUrlRep
	data, err := androidOrIosLoadUrl.FindAndroidOrIosLoadUrl(qipaiDB, gameID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return models.AndroidOrIosLoadUrlRep{}, err
	}
	return data, nil
}
