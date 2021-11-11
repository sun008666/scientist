package service

import (
	"github.com/helloteemo/ashe"
	"regexp"
	"testing"
	"zonst/qipai/api/configapisrv/models"
)

func TestFindAndroidOrIosLoadUrl(t *testing.T) {
	cat := ashe.New(t).Use(ashe.DB)
	cat.UnitTesting(`FindAndroidOrIosLoadUrl`, func() {
		db := cat.GetGormConn()

		rows := cat.NewRows([]string{"game_id", "android_url", "ios_url"}).AddRow(7, "huahuadashuaige", "huahuadashuaige")
		cat.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(1).WillReturnRows(rows)

		user, err := FindAndroidOrIosLoadUrl(db, 7)
		cat.Nil(err)
		cat.Equal(user, models.AndroidOrIosLoadUrlRep{GameID: 7, AndroidUrl: "huahuadashuaige", IosUrl: "huahuadashuaige"})
	})
}
