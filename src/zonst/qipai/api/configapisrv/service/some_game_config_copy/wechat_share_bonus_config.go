package some_game_config_copy

import (
	"encoding/json"
	"fmt"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/utils"
)

type WechatShareBonusConfig struct{}

var _ Interface = WechatShareBonusConfig{}

// Copy 复制道具系统
func (s WechatShareBonusConfig) Copy(param Param) error {
	wechatShareBonusURL := fmt.Sprintf("%v%v", param.API, "/web/v1/config/batch/copy")
	tempMap := make(map[string]interface{})
	tempMap["src_game_id"] = param.SrcGameID
	tempMap["dest_game_ids"] = param.DestGameIDs
	jsonBodyBytes, err := json.Marshal(tempMap)
	if err != nil {
		logging.Errorf("json marshal error")
		return err
	}
	if _, err := utils.SomeGameConfigCopyHttpRequest(wechatShareBonusURL, string(jsonBodyBytes), param.JwtToken); err != nil {
		logging.Errorf("Copy-WechatShareBonusConfig-HttpRequest: body:%v,err:%v\n", string(jsonBodyBytes), err)
		return err
	}
	return nil
}
