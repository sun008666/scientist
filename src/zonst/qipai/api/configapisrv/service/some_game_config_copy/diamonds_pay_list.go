package some_game_config_copy

import (
	"encoding/json"
	"fmt"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/utils"
)

type DiamondsPayList struct{}

var _ Interface = DiamondsPayList{}

// Copy 复制钻石套餐
func (s DiamondsPayList) Copy(param Param) error {
	signAwardURL := fmt.Sprintf("%v%v", param.API, "/web/v1/diamond/package/config/batch/copy")
	tempMap := make(map[string]interface{})
	tempMap["game_id1"] = param.SrcGameID
	tempMap["game_id2"] = param.DestGameIDs
	jsonBodyBytes, err := json.Marshal(tempMap)
	if err != nil {
		logging.Errorf("json marshal error")
		return err
	}
	if _, err := utils.SomeGameConfigCopyHttpRequest(signAwardURL, string(jsonBodyBytes), param.JwtToken); err != nil {
		logging.Errorf("Copy-DiamondsPay-HttpRequest: body:%v,err:%v\n", string(jsonBodyBytes), err)
		return err
	}
	return nil
}
