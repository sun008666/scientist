package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"zonst/logging"
)

func SomeGameConfigCopyHttpRequest(url string, body string, token string) ([]byte, error) {
	t1 := time.Now()
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	content := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, content)
	req.Header.Add("x-xq5-jwt", token)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	strBody := strings.Replace(string(rspBody), "\n", "", -1)
	logging.Debugf("SomeGameConfigCopyHttpRequest: url:%v body:%v Response: status:%v body:%v cost:%v\n",
		url, body, rsp.Status, strBody, time.Since(t1))
	return rspBody, nil
}
