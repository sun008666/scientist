package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"zonst/logging"
)

type CommRespone struct {
	ErrNo  string `json:"errno"`
	ErrMsg string `json:"errmsg"`
}

// AddConfigLog 添加配置日志
func AddConfigLog(URL string, args map[string]interface{}) error {
	t1 := time.Now()
	bytesData, err := json.Marshal(args)
	if err != nil {
		logging.Errorf("AddConfigLog:序列化化出错了:err:%v\n", err)
		return err
	}
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	logging.Debugf("请求参数列表%v\n", args)
	req, err := http.NewRequest("POST", URL+"/v1/config/log/add", bytes.NewReader(bytesData))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	rsp, _ := client.Do(req)

	defer rsp.Body.Close()
	obj := &CommRespone{}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	strBody := strings.Replace(string(body), "\n", "", -1)
	logging.Debugf("AddConfigLog: Request: url:%v Response: status:%v body:%v cost:%v\n",
		URL, rsp.Status, strBody, time.Since(t1))
	if err := json.Unmarshal(body, obj); err != nil {
		return err
	}
	if obj.ErrNo != "0" {
		return errors.New(obj.ErrMsg)
	}
	return nil

}
