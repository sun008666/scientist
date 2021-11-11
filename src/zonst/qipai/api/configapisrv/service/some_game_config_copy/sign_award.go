package some_game_config_copy

import (
	"encoding/json"
	"fmt"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/utils"
)

type SignAward struct {
}

var _ Interface = SignAward{}

// Copy 复制累计签到
func (s SignAward) Copy(param Param) error {
	signAwardURL := fmt.Sprintf("%v%v", param.API, "/v1/sign/award/config/batch/copy")
	tempMap := make(map[string]interface{})
	tempMap["src_game_id"] = param.SrcGameID
	tempMap["dest_game_id"] = param.DestGameIDs
	jsonBodyBytes, err := json.Marshal(tempMap)
	if err != nil {
		logging.Errorf("json marshal error")
		return err
	}
	if _, err := utils.SomeGameConfigCopyHttpRequest(signAwardURL, string(jsonBodyBytes), param.JwtToken); err != nil {
		logging.Errorf("Copy-SignAward-HttpRequest: body:%v,err:%v\n", string(jsonBodyBytes), err)
		return err
	}
	return nil
}
