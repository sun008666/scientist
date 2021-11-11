package utils

import (
	"bytes"
	"encoding/json"
	"github.com/fwhezfwhez/errorx"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/dependency/errs"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
)

type cdnResponse struct {
	Code int `json:"code"`
}

// RefreshCDNURL 刷新cdn url
//func RefreshCDNURL(c *gin.Context, uri string) int {
//	/**get SecretKey & SecretId from https://console.qcloud.com/capi**/
//
//	cdnSecretID := c.MustGet("config").(*config.Config).CdnSecretID
//	cdnSecretKEY := c.MustGet("config").(*config.Config).CdnSecretKEY
//
//	var Requesturl string = "cdn.api.qcloud.com/v2/index.php"
//	var SecretKey string = cdnSecretKEY
//	var Method string = "POST"
//
//	/**params to signature**/
//	params := make(map[string]interface{})
//	params["SecretId"] = cdnSecretID
//	params["Action"] = "RefreshCdnUrl"
//	params["urls.0"] = uri
//
//	/*use qcloudcdn_api.Signature to obtain signature and params with correct signature**/
//	signature, request_params := Signature(SecretKey, params, Method, Requesturl)
//	log.Println("signature : ", signature)
//
//	/*use qcloudcdn_api.SendRequest to send request**/
//	response := SendRequest(Requesturl, request_params, Method)
//	log.Printf("cdn刷新response: %v\n", response)
//
//	obj := &cdnResponse{}
//	if err := json.Unmarshal([]byte(response), obj); err != nil {
//		log.Errorf("RefreshCDNURL: err: %v\n", err.Error())
//		return -1
//	}
//	return obj.Code
//}

// RefreshCDNDir 刷新cdn dir
func RefreshCDNDir(c *gin.Context, dir string) int {
	/**get SecretKey & SecretId from https://console.qcloud.com/capi**/

	cdnSecretID := c.MustGet("config").(*config.Config).CdnSecretID
	cdnSecretKEY := c.MustGet("config").(*config.Config).CdnSecretKEY

	var Requesturl string = "cdn.api.qcloud.com/v2/index.php"
	var SecretKey string = cdnSecretKEY
	var Method string = "POST"

	/**params to signature**/
	params := make(map[string]interface{})
	params["SecretId"] = cdnSecretID
	params["Action"] = "RefreshCdnDir"
	params["dirs.0"] = dir

	/*use qcloudcdn_api.Signature to obtain signature and params with correct signature**/
	signature, request_params := Signature(SecretKey, params, Method, Requesturl)
	log.Println("signature : ", signature)

	/*use qcloudcdn_api.SendRequest to send request**/
	response := SendRequest(Requesturl, request_params, Method)
	log.Printf("cdn刷新response: %v\n", response)

	obj := &cdnResponse{}
	if err := json.Unmarshal([]byte(response), obj); err != nil {
		log.Errorf("RefreshCDNDir: err: %v\n", err.Error())
		return -1
	}
	return obj.Code
}

type CDNapiResp struct {
	Errno         string `json:"errno"`
	Errmsg        string `json:"errmsg"`
	RefreshResult int    `json:"refreshResult"`
}

func RefreshCDNURL(uri string) int {
	values := url.Values{}
	environmentFlag := config.Cfg.EnvironmentFlag
	var refreshenvironment string
	switch environmentFlag {
	case "tencent-cos":
		refreshenvironment = "3"
	case "huawei-obs":
		refreshenvironment = "2"
	default:
		log.Errorf("RefresCDN: err:环境标志出错\n")
		return -1
	}

	values.Add("refreshenvironment", refreshenvironment)
	values.Add("refreshurl", uri)
	cdnAPISrv := config.Cfg.CDNAPISrv
	resp, e := HTTPTimeoutPost(cdnAPISrv, values.Encode())
	if e != nil {

		errs.SaveError(errorx.Wrap(e), map[string]interface{}{
			"uri": uri,
		})
		log.Errorf("RefreshCDN: err:%v\n", e.Error())
		return -1
	}

	var jsonUnmarshaledResp CDNapiResp
	e = json.Unmarshal(resp, &jsonUnmarshaledResp)
	if e != nil {
		log.Errorf("RefreshCDN: err:%v\n", e.Error())

		errs.SaveError(errorx.Wrap(e), map[string]interface{}{
			"uri": uri,
		})
		return -1
	}

	if jsonUnmarshaledResp.Errno != "0" || jsonUnmarshaledResp.RefreshResult == 0 || jsonUnmarshaledResp.RefreshResult == 2 {
		errs.SaveError(errorx.NewServiceError("未返回正确信息", 2), map[string]interface{}{
			"uri": uri,
			"rs":  jsonUnmarshaledResp,
		})

		log.Errorf("RefreshCDN: 未成功，返回信息：%+v", jsonUnmarshaledResp)
		return -1
	}

	return 0
}

func HTTPTimeoutPost(url string, content string) ([]byte, error) {
	t1 := time.Now()
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	contentBuffer := bytes.NewBufferString(content)
	req, err := http.NewRequest("POST", url, contentBuffer)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	rsp, err := client.Do(req)
	if err != nil {
		logging.Errorf("HTTPTimeoutPost: request_url:%v, err:%v\n", url, err)
		return []byte(""), err
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)

	// 记录日志
	logging.Debugf("HttpTimeoutPost: request_url:%v, request_body:%v, response_body:%v, cost: %v",
		url, content, string(body), time.Since(t1))
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}
