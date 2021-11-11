package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"zonst/qipai/api/configapisrv/config"

	"github.com/gin-gonic/gin"

	"github.com/go-xweb/log"
)

// PopUpControl 弹窗控制表
type PopUpControl struct {
	ID       int `gorm:"id" json:"id"`
	GameID   int `gorm:"game_id" json:"game_id"`
	Strategy int `gorm:"strategy" json:"strategy"`
}

type RefleshCacheResponse struct {
	ErrNo  int    `json:"errno"`
	ErrMsg string `json:"errmsg"`
}

func RefleshCache(c *gin.Context, cacheName string) bool {
	reqURL := fmt.Sprintf("%s/reload_cache/", c.MustGet("config").(*config.Config).GameClientApi)

	postValues := url.Values{}
	postValues.Set("cache_name", cacheName)

	res, err := http.PostForm(reqURL, postValues)
	if err != nil {
		log.Errorf("RefleshCache: err: %v\n", err.Error())
		return false
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("RefleshCache: err: %v\n", err.Error())
		return false
	}

	rp := &RefleshCacheResponse{}
	if err := json.Unmarshal(body, rp); err != nil {
		log.Errorf("RefleshCache: err: %v\n", err.Error())
		return false
	}
	if rp.ErrNo != 0 {
		log.Errorf("RefleshCache: err: %v\n", rp.ErrMsg)
		return false
	}
	return true
}
