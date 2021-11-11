package some_game_config_copy

import (
	"encoding/json"
	"fmt"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/utils"
)

type ShopConfig struct{}

var _ Interface = ShopConfig{}

// Copy 复制道具系统
func (s ShopConfig) Copy(param Param) error {
	signAwardURL := fmt.Sprintf("%v%v", param.API, "/v1/shop/config/batch/copy")
	tempMap := make(map[string]interface{})
	tempMap["from_game_id"] = param.SrcGameID
	tempMap["to_game_id"] = param.DestGameIDs
	jsonBodyBytes, err := json.Marshal(tempMap)
	if err != nil {
		logging.Errorf("json marshal error")
		return err
	}
	if _, err := utils.SomeGameConfigCopyHttpRequest(signAwardURL, string(jsonBodyBytes), param.JwtToken); err != nil {
		logging.Errorf("Copy-ShopConfig-HttpRequest: body:%v,err:%v\n", string(jsonBodyBytes), err)
		return err
	}
	return nil
}
