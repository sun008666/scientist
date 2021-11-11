package errs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"
	"zonst/qipai/api/configapisrv/config"
)

var c = http.Client{
	Timeout: 15 * time.Second,
}
var errLimitMap = sync.Map{}

var l = rate.NewLimiter(1, 50)

func SaveError(e error, context ...map[string]interface{}) string {

	r := l.Reserve()
	delay := r.Delay()

	if delay > 3*time.Second {
		return "too frequent"
	}

	if len(context) != 1 {
		context = []map[string]interface{}{
			{
				"srv": config.Node.AppName,
			},
		}
	}

	context[0]["host"], _ = os.Hostname()

	context[0]["srv_name"] = config.Node.AppName

	context[0]["main_path"] = config.Node.MainPath

	context[0]["mode"] = config.Node.Mode

	// 如果keyword="",则keyword等于label
	if GetString(context[0], "keyword") == "" {
		context[0]["keyword"] = GetString(context[0], "label")
	}

	// 模块
	elem := StringDefault(GetString(context[0], "elem"), config.Node.AppName)

	// 标签
	label := StringDefault(GetString(context[0], "label"), fmt.Sprintf("%s-common", config.Node.AppName))

	return SaveE(e, elem, label, context...)
}

// context特殊值:
// tip: 中文描述，会上报
// elem: 模块
// label: 标签
func SaveE(e error, elem string, label string, context ...map[string]interface{}) string {
	if len(context) != 1 {
		context = []map[string]interface{}{
			{
				"srv": config.Node.AppName,
			},
		}
	}
L:
	switch v := e.(type) {
	case errorx.Error:
		break L
	case error:
		return SaveError(errorx.NewFromString(string(fmt.Sprintf("err '%s' \n %s", v.Error(), debug.Stack()))), context...)
	}

	report(
		elem,                         //模块
		label,                        // 标签
		e.Error(),                    // 错误堆栈
		GetString(context[0], "tip"), // 中文描述
		context[0],                   // 键值附加信息
	)
	return ""
}

func report(business string, label string, stack string, msg string, kvs map[string]interface{}) {

	var url string

	switch config.Node.Mode {
	case "pro":
		url = "https://log.zzect.com/api/v1/backend/report"
	default:
		url = "http://129.211.113.36:9097/api/v1/backend/report"

	}

	type Request struct {
		Level    string                 `json:"level"`  // 级别
		Caller   string                 `json:"caller"` // 堆栈
		Srv      string                 `json:"srv"`    //
		Host     string                 `json:"host"`
		Business string                 `json:"business"` // 模块
		Label    string                 `json:"label"`    // 报警标签
		Msg      string                 `json:"msg"`      // 打印信息
		Kvs      map[string]interface{} `json:"kvs"`      // 附加map
	}

	var req Request

	req.Level = "ERROR"
	req.Srv = config.Node.AppName

	req.Business = business
	req.Label = label

	req.Caller = stack
	req.Msg = msg
	req.Kvs = kvs

	host, e := os.Hostname()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	req.Host = host

	b, _ := json.Marshal(req)

	reqst, e := http.NewRequest("POST", url, bytes.NewReader(b))
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	reqst.Header.Set("Content-Type", "application/json")

	rsp, e := c.Do(reqst)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}

	rspb, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	fmt.Println(string(rspb))

	return
}
