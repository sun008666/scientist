package some_game_config_copy

import (
	"encoding/json"
	"fmt"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/utils"
)

type SubsidyConfig struct {
}

var _ Interface = SubsidyConfig{}

// Copy 复制累计签到
func (s SubsidyConfig) Copy(param Param) error {
	subsidyConfigBatchCopyURL := fmt.Sprintf("%v%v", param.API, "/web/v1/subsidy/config/batch/copy")
	tempMap := make(map[string]interface{})
	tempMap["src_game_id"] = param.SrcGameID
	tempMap["dest_game_ids"] = param.DestGameIDs
	jsonBodyBytes, err := json.Marshal(tempMap)
	if err != nil {
		logging.Errorf("json marshal error")
		return err
	}
	if _, err := utils.SomeGameConfigCopyHttpRequest(subsidyConfigBatchCopyURL, string(jsonBodyBytes), param.JwtToken); err != nil {
		logging.Errorf("Copy-SubsidyConfig-HttpRequest: body:%v,err:%v\n", string(jsonBodyBytes), err)
		return err
	}
	return nil
}
