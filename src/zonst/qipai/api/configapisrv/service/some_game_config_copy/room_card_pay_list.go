package some_game_config_copy

import (
	"zonst/logging"
	"zonst/qipai/api/configapisrv/models"

	"github.com/go-xweb/log"
)

type RoomCardPayList struct{}

var _ Interface = RoomCardPayList{}

// Copy 复制房卡套餐
func (s RoomCardPayList) Copy(param Param) error {
	//删除目标平台所有配置
	//err := models.DeleteRoomcardProductListByToGameID(param.QipaiDBOrm, param.DestGameID)
	err := models.BatchDeleteRoomcardProductListByToGameID(param.QipaiDBOrm, param.DestGameIDs)
	if err != nil {
		log.Errorf("Copy-ProductConfig: %v\n", err.Error())
		return err
	}
	//找出所有套餐信息
	data, err := models.GetProductConfig(param.QipaiDBOrm, int(param.SrcGameID))
	if err != nil {
		logging.Errorf("Copy-ProductConfig: err: %v\n", err.Error())
		return err
	}
	newProductList := make([]models.RoomCardProductList, 0)
	for _, gameId := range param.DestGameIDs {
		for _, v := range data {
			//去除掉卡包和活动房卡套餐
			//if v.GameAreaID != "-1" {
			//	continue
			//}
			//修改列表的平台 添加时间
			temp := v
			temp.GameID = gameId
			addTime, _ := v.AddTime.MarshalJSON()
			temp.AddTimeStr = string(addTime)
			temp.ID = 0
			endTime, _ := v.EndTime.MarshalJSON()
			temp.EndTimeStr = string(endTime)
			newProductList = append(newProductList, temp)
		}
	}

	//批量插入套餐列表
	err = models.InsertIntoRoomcardProductList(param.QipaiDBOrm, newProductList)
	if err != nil {
		log.Errorf("Copy-ProductConfig: err: %v\n", err.Error())
		return err
	}

	return nil
}
