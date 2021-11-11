package main

import (
	"flag"
	"runtime"
	"zonst/qipai-sports/api/configapisrv/constants"

	"zonst/logging"
	"zonst/qipai-sports/api/configapisrv/config"
	"zonst/qipai-sports/api/configapisrv/handlers"
	"zonst/qipai-sports/api/configapisrv/middlewares"

	"github.com/gin-gonic/gin"
)

var (
	tomlFile = flag.String("config", "docs/test.toml", "config file")
)

// init 初始化配置
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	flag.Parse()

	// 解析配置文件
	tomlConfig, err := config.UnmarshalConfig(*tomlFile)
	if err != nil {
		logging.Errorf("UnmarshalConfig: err:%v\n", err)
		return
	}

	// 绑定路由，及公共的tomlConfig
	// router := gin.Default()
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(middlewares.Config(tomlConfig))
	router.Use(middlewares.Gorm(constants.Configdb, tomlConfig))

	// 路由配置
	router.GET("/", handlers.Index) // 首页

	v1 := router.Group("/v1")                               // v1
	{
		v1.GET("/gamelist", handlers.OnGameListRequest)			//获取游戏平台列表
	}



	// 启动服务
	logging.Debugf("run srvapisrv at %v\n", tomlConfig.GetListenAddr())
	router.Run(tomlConfig.GetListenAddr())
}
