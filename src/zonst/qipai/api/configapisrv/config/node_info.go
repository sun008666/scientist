package config

import (
	"fmt"
	"runtime"
)

var Cfg *Config

var Node = &ServiceInfo{}

type ServiceInfo struct {
	AppName  string `json:"app_name"`
	Mode     string `json:"mode"`
	MainPath string `json:"main_path"`
}

func RegisterNode(appName string) {
	Node.AppName = appName

	_, f, l, _ := runtime.Caller(1)
	Node.MainPath = fmt.Sprintf("%s:%d", f, l)
}
