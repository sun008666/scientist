package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jmoiron/sqlx/types"
)

// OnGameConfigListReq 游戏配置列表
type OnGameConfigListReq struct {
	GameID int `json:"game_id" binding:"required"`
}

// GameConfig 游戏配置
type GameConfig struct {
	Copyright                   string         `json:"copyright"`
	GameNotice                  string         `json:"game_notice"`
	GameShareTitle              string         `json:"game_share_title"`
	GameShareURL                string         `json:"game_share_url"`
	GameShareImage              string         `json:"game_share_image"`
	GameShowRank                string         `json:"game_show_rank"`
	IosCheck                    string         `json:"ios_check"`
	Version                     string         `json:"version"`
	LoginIP                     string         `json:"login_ip"`
	WxName                      string         `json:"wx_name"`
	WxName1                     string         `json:"wx_name1"`
	WxValue                     string         `json:"wx_value"`
	WxValue1                    string         `json:"wx_value1"`
	IsNeedPhone                 bool           `json:"is_need_phone"`
	GameAreaShareTitle          types.JSONText `json:"game_area_share_title"`
	DistrictWxList              types.JSONText `json:"district_wx_list"`
	GamePublishNumber           string         `json:"game_publish_number"`           // 网络游戏出版物号
	ApprovalNumber              string         `json:"approval_number"`               // 网络游戏出版物号
	CopyrightRegistrationNumber string         `json:"copyright_registration_number"` // 著作权登记号
	OnlineGamePrepareWord       string         `json:"online_game_prepare_word"`      // 文网游备字
	CopyrightOwner              string         `json:"copyright_owner"`               // 著作权人
	Publisher                   string         `json:"publisher"`                     // 出版单位
}

// HTTPGetResponce HTTPGetResponce
type HTTPGetResponce struct {
	Errno    string     `json:"errno"`
	Errmsg   string     `json:"errmsg"`
	DataStr  string     `json:"data"`
	DataJSON GameConfig `json:"-"`
}

// OnGameConfigListRequest 游戏配置列表
func OnGameConfigListRequest(c *gin.Context) {
	req := &OnGameConfigListReq{}
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGameConfigListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	gameConfigAddr := c.MustGet("config").(*config.Config).GameConfigAddr
	url := gameConfigAddr + "/v1/config/get?game_id=" + strconv.Itoa(req.GameID)

	data, err := httpGet(url)
	if err != nil && err.Error() != "unexpected end of JSON input" {
		log.Errorf("OnGameConfigListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": data.DataJSON})

}

func httpGet(url string) (data HTTPGetResponce, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return data, err
	}

	if err = json.Unmarshal(body, &data); err != nil {
		log.Error(err)
		return data, err
	}

	if err = json.Unmarshal([]byte(data.DataStr), &data.DataJSON); err != nil {
		log.Error(err)
		return data, err
	}
	return data, err
}

// OnGameConfigUpdateReq OnGameConfigUpdateReq
type OnGameConfigUpdateReq struct {
	GameID                      int            `json:"game_id"`
	Copyright                   string         `json:"copyright"`
	GameNotice                  string         `json:"game_notice"`
	GameShareTitle              string         `json:"game_share_title"`
	GameShareURL                string         `json:"game_share_url"`
	GameShareImage              string         `json:"game_share_image"`
	GameShowRank                string         `json:"game_show_rank"`
	IosCheck                    string         `json:"ios_check"`
	Version                     string         `json:"version"`
	LoginIP                     string         `json:"login_ip"`
	WxName                      string         `json:"wx_name"`
	WxName1                     string         `json:"wx_name1"`
	WxValue                     string         `json:"wx_value"`
	WxValue1                    string         `json:"wx_value1"`
	IsNeedPhone                 bool           `json:"is_need_phone"`
	GameAreaShareTitle          types.JSONText `json:"game_area_share_title"`
	DistrictWxList              types.JSONText `json:"district_wx_list"`
	GamePublishNumber           string         `json:"game_publish_number"`           // 网络游戏出版物号
	ApprovalNumber              string         `json:"approval_number"`               // 网络游戏出版物号
	CopyrightRegistrationNumber string         `json:"copyright_registration_number"` // 著作权登记号
	OnlineGamePrepareWord       string         `json:"online_game_prepare_word"`      // 文网游备字
	CopyrightOwner              string         `json:"copyright_owner"`               // 著作权人
	Publisher                   string         `json:"publisher"`                     // 出版单位
}

// 游戏配置修改
func OnGameConfigUpdateRequest(c *gin.Context) {
	req := &OnGameConfigUpdateReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnGameConfigUpdateRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGameConfigUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	data, err := httpPostForm(c, req)
	if err != nil {
		log.Errorf("OnGameConfigUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	if data.Errno != "0" {
		log.Errorf("OnGameConfigUpdateRequest: err: %v\n", data.Errmsg)
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": data.Errmsg})
		return
	}
	config := c.MustGet("config").(*config.Config)
	uri := config.ConfigDomainName + "/" + strconv.Itoa(req.GameID) + ".json"

	responce := utils.RefreshCDNURL(uri)
	if responce != 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

// HTTPPostFormResponce HTTPPostFormResponce
type HTTPPostFormResponce struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

func httpPostForm(c *gin.Context, req *OnGameConfigUpdateReq) (data HTTPPostFormResponce, err error) {
	gameConfigAddr := c.MustGet("config").(*config.Config).GameConfigAddr
	uri := gameConfigAddr + "/v1/config/update"

	//这里添加post的body内容
	formData := make(url.Values)
	formData["game_id"] = []string{strconv.Itoa(req.GameID)}
	formData["copyright"] = []string{req.Copyright}
	formData["game_notice"] = []string{req.GameNotice}
	formData["game_share_title"] = []string{req.GameShareTitle}
	formData["game_share_url"] = []string{req.GameShareURL}
	formData["game_share_image"] = []string{req.GameShareImage}
	formData["game_show_rank"] = []string{req.GameShowRank}
	formData["ios_check"] = []string{req.IosCheck}
	formData["version"] = []string{req.Version}
	formData["login_ip"] = []string{req.LoginIP}
	formData["wx_name"] = []string{req.WxName}
	formData["wx_name1"] = []string{req.WxName1}
	formData["wx_value"] = []string{req.WxValue}
	formData["wx_value1"] = []string{req.WxValue1}
	formData["is_need_phone"] = []string{strconv.FormatBool(req.IsNeedPhone)}
	formData["game_area_share_title"] = []string{req.GameAreaShareTitle.String()}
	formData["district_wx_list"] = []string{req.DistrictWxList.String()}
	formData["game_publish_number"] = []string{req.GamePublishNumber}
	formData["approval_number"] = []string{req.ApprovalNumber}
	formData["copyright_registration_number"] = []string{req.CopyrightRegistrationNumber}
	formData["online_game_prepare_word"] = []string{req.OnlineGamePrepareWord}
	formData["copyright_owner"] = []string{req.CopyrightOwner}
	formData["publisher"] = []string{req.Publisher}

	resp, err := http.PostForm(uri, formData)
	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return data, err
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return data, err
	}
	return data, err
}
