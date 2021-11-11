package some_game_config_copy

import (
	"errors"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/models"
)

type Task struct{}

var _ Interface = Task{}

// Copy 复制任务系统
func (t Task) Copy(param Param) error {
	result, err := models.GetCopyTaskConfigList(param.TaskDBOrm, int32(param.SrcGameID))
	if err != nil {
		logging.Errorf("Copy-Task: 获取任务配置列表错误，err:%+v\n", err)
		return err
	}
	if len(result) == 0 {
		logging.Errorf("Copy-TaskConfig: err:%+v\n", errors.New("未查询到相关任务配置"))
		return errors.New("未查询到相关任务配置")
	}
	//删除目标平台所有配置
	if err = models.BatchDeleteCopyTaskConfigByDestGameID(param.TaskDB, param.DestGameIDs); err != nil {
		logging.Errorf("CopyTaskConfigListToNewGameID: err:%v\n", err)
		return errors.New("删除目标平台配置出错")
	}
	// 复制任务配置到新游戏平台
	if err := models.BatchCopyTaskConfigToGameIdList(param.TaskDB, param.DestGameIDs, result); err != nil {
		logging.Errorf("Copy-TaskConfig: 复制失败, err:%v\n", err)
		return err
	}
	return nil
}
