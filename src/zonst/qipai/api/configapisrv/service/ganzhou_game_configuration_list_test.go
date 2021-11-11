package service

import (
	"github.com/helloteemo/ashe"
	"regexp"
	"testing"
	"zonst/qipai/api/configapisrv/models"
)

func TestFindConfig(t *testing.T) {
	cat := ashe.New(t).Use(ashe.DB)
	cat.UnitTesting(`TestFindConfig`, func() {
		db := cat.GetGormConn()
		rows := cat.NewRows([]string{"game_id", "game_area_id", "id"}).AddRow(7, 30, 1)
		cat.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(rows)
		userInfoModel, err := models.FindConfig(db)
		cat.Nil(err)
		cat.Equal(userInfoModel, []models.GanzhouGameConfiguration{{GameID: 7, GameAreaID: 30, ID: 1}})
	})
}
