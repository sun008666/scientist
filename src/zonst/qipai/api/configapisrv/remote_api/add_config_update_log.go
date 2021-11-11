package remote_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"net/url"
	"time"
	"zonst/qipai/gin-middlewares/secretauth"
)

const (
	addConfigUpdateLogAPI = "/v1/config/log/add"           // 添加配置更新日志
	downHtmlUpdateLogAPI  = "/v1/download/page/url/update" // 更新下载页面

	ModuleIDRoomCardPackage   = 1
	ModuleIDYuanBaoPackage    = 2
	ModuleIDCoinRoomConfig    = 3
	ModuleIDCoinPackageConfig = 4
	ModuleIDGameConfig        = 5
	ModuleIDHomePageConfig    = 6

	OpTypeAdd    = 1
	OpTypeUpdate = 2
	OpTypeDel    = 3
)

type addConfigUpdateLogResp struct {
	Errno  string `json:"errno" form:"errno"`
	Errmsg string `json:"errmsg" form:"errmsg"`
}

// AddConfigUpdateLog 添加配置更新日志
func AddConfigUpdateLog(gameID, moduleID, operateID, opType int32, beforeContent, afterContent, host, appID, appKey string) error {
	values := url.Values{}
	values.Add("game_id", fmt.Sprintf("%v", gameID))
	values.Add("module_id", fmt.Sprintf("%v", moduleID))
	values.Add("operate_id", fmt.Sprintf("%v", operateID))
	values.Add("before_content", beforeContent)
	values.Add("after_content", afterContent)
	values.Add("op_type", fmt.Sprintf("%v", opType))
	values.Add("timestamp", fmt.Sprintf("%v", time.Now().Unix()))
	values.Add("appid", appID)

	var valuem = map[string]string{
		"game_id":        fmt.Sprintf("%v", gameID),
		"module_id":      fmt.Sprintf("%v", moduleID),
		"operate_id":     fmt.Sprintf("%v", operateID),
		"before_content": beforeContent,
		"after_content":  afterContent,
		"op_type":        fmt.Sprintf("%v", opType),
		"timestamp":      fmt.Sprintf("%v", time.Now().Unix()),
		"appid":          appID,
	}

	// 包装sign
	sign := secretauth.CalculcateSign(valuem, appKey)
	values.Add("sign", sign)

	// 发送请求
	api := host + addConfigUpdateLogAPI

	body, e := c.PostForm(api, values)
	if e != nil {
		return errorx.Wrap(e)
	}

	if body == nil {
		return errorx.NewServiceError("rsp nil", 1)
	}
	if body.Body == nil {
		return errorx.NewServiceError("body nil", 2)
	}

	// 读取响应
	resp := addConfigUpdateLogResp{}
	if e = json.NewDecoder(body.Body).Decode(&resp); e != nil {
		return errorx.Wrap(e)
	}

	// 判断响应是否成功
	if resp.Errno != "0" {
		e := errors.New(resp.Errmsg)
		return errorx.Wrap(e)
	}

	return nil
}

// UpdateDownloadHtmlURL 更新下载地址
func UpdateDownloadHtmlURL(id, gameID int, gameName, enName, androidUrl, iosUrl, adURL, token string, IsAdShow, IsAndroidRedirect, IsIosRedirect bool) (err error) {
	values := url.Values{}
	values.Add("id", fmt.Sprintf("%v", id))
	values.Add("game_id", fmt.Sprintf("%v", gameID))
	values.Add("game_name", gameName)
	values.Add("en_name", enName)
	values.Add("android_url", androidUrl)
	values.Add("ios_url", iosUrl)
	values.Add("ad_url", adURL)
	values.Add("is_ad_show", fmt.Sprintf("%v", IsAdShow))
	values.Add("is_android_redirect", fmt.Sprintf("%v", IsAndroidRedirect))
	values.Add("is_ios_redirect", fmt.Sprintf("%v", IsIosRedirect))

	api := "http://127.0.0.1:9901" + downHtmlUpdateLogAPI
	resp, err := HTTPTimeoutPost(api, values.Encode(), token)
	if err != nil {
		return err
	}

	// 从返回的json数据中解析
	r := &addConfigUpdateLogResp{}
	if err = json.Unmarshal(resp, r); err != nil {
		return err
	}

	if r.Errno != "0" {
		return errors.New(r.Errmsg)
	}
	return nil
}
