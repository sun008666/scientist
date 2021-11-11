package remote_api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

var c = http.Client{
	Timeout: 15 * time.Second,
}

// HTTPTimeoutPost 发起POST请求
func HTTPTimeoutPost(url string, content, jwtToken string) ([]byte, error) {
	timeout := time.Duration(time.Second * 10)
	client := http.Client{
		Timeout: timeout,
	}

	contentBuffer := bytes.NewBufferString(content)
	req, err := http.NewRequest("POST", url, contentBuffer)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Add("x-xq5-jwt", jwtToken)
	rsp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return []byte(""), err
	}
	return body, nil
}
