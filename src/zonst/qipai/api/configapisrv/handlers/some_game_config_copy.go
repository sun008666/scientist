package handlers

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/service/some_game_config_copy"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SomeGameConfigCopyReq struct {
	CopyTypeList []int32 `json:"copy_type_list" form:"copy_type_list" binding:"required"`
	SrcGameID    int32   `json:"src_game_id" form:"src_game_id" binding:"required"`
	DestGameIDs  []int32 `json:"dest_game_id" form:"dest_game_id" binding:"required"`
}

// OnSomeGameConfigCopyRequest 平台复制模块
func OnSomeGameConfigCopyRequest(c *gin.Context) {
	req := SomeGameConfigCopyReq{}
	if err := c.ShouldBind(&req); err != nil {
		logging.Errorf("OnSomeGameConfigCopyRequest 参数绑定异常 err:%v\n", err)
		c.JSON(http.StatusOK, gin.H{"errno": "A0403", "errmsg": "参数不匹配，请重试"})
		return
	}
	if len(req.CopyTypeList) != 0 {
		jwtToken := c.Request.Header.Get("x-xq5-jwt")
		cfg := c.MustGet("config").(*config.Config)
		for key := range req.CopyTypeList {
			CopyType := some_game_config_copy.CopyType(req.CopyTypeList[key])
			copy_type, err := some_game_config_copy.InterfaceFactory(CopyType)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"errno": "A0403", "errmsg": err})
				return
			}
			var param some_game_config_copy.Param
			param.API = getGameConfigCopyAPI(CopyType, cfg)
			if param.API == "" {
				c.JSON(http.StatusOK, gin.H{"errno": "A0403", "errmsg": "参数不匹配，请重试"})
				return
			}
			param.JwtToken = jwtToken
			param.SrcGameID = req.SrcGameID
			param.DestGameIDs = req.DestGameIDs
			param.TaskDBOrm = c.MustGet(utils.TaskDBOrm).(*gorm.DB)
			param.TaskDB = c.MustGet(utils.TaskDB).(*sqlx.DB)
			param.QipaiDBOrm = c.MustGet(utils.QipaiDBOrm).(*gorm.DB)

			if err := copy_type.Copy(param); err != nil {
				logging.Errorf("OnSomeGameConfigCopyRequest: 复制失败，req:%v，err:%v\n", req, err)
				c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "复制失败"})
				return
			}

		}
		c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "复制成功"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "A0403", "errmsg": "平台配置复制列表为空"})
	return

}

func getGameConfigCopyAPI(copyType some_game_config_copy.CopyType, cfg *config.Config) string {
	switch copyType {
	case 1:
		return cfg.SignAwardAPI
	case 2:
		return cfg.ShopConfigAPi
	case 3:
		return cfg.TaskAPI
	case 4:
		return cfg.RoomCardPayAPi
	case 5:
		return cfg.DiamondsPayApi
	case 6:
		return cfg.SubsidyConfigAPI
	case 7:
		return cfg.WechatShareBonusConfigAPI
	default:
		return ""
	}
}
