package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"zonst/logging"
	"zonst/qipai-sports/api/configapisrv/constants"
)


type Response struct {
	Code int `json:"err_code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type Game struct {
	Status        int    `gorm:"column:status" json:"status"`
	GameID            int    `gorm:"column:game_id" json:"game_id"`
	GameName       string    `gorm:"column:game_name" json:"game_name"`
}

type GameListResponse struct {

}
// OnGameListRequest 获取游戏平台列表
func OnGameListRequest(c *gin.Context) {
	db := c.MustGet(constants.Configdb).(*gorm.DB)
	gameList := make([]*Game, 0)

	if err := db.Table("game_game").Order("game_id").Find(&gameList).Error; err != nil{
		logging.Errorf("OnGameListRequest: query err:%v\n", err)
		resp := &Response{Code:1, Msg:"获取游戏平台列表失败"}
		c.JSON(http.StatusOK, resp)
		return
	}

	if data,err := json.Marshal(gameList); err != nil{
		logging.Errorf("OnGameListRequest: Marshal err:%v\n", err)
		resp := &Response{Code:1, Msg:"解析json数据失败"}
		c.JSON(http.StatusOK, resp)
		return
	}else{
		resp := &Response{Code:1, Msg:"获取游戏平台列表成功", Data:string(data)}
		c.JSON(http.StatusOK, resp)
		return
	}
}
