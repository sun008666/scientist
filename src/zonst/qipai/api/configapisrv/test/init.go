package test

import (
	"flag"
	"log"
	"os"
	"runtime"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/handlers"
	"zonst/qipai/api/configapisrv/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	// router 全局路由
	router   *gin.Engine
	tomlFile string
	mode     string
	clubdb   *sqlx.DB
	//clubcache *redis.Pool
)

func init() {
	log.SetFlags(log.Ltime | log.Llongfile | log.Ldate)
	flag.StringVar(&tomlFile, "config", "../docs/test.toml", "服务配置文件")
	flag.StringVar(&mode, "mode", "DEBUG", "模型-DEBUG还是RELEASE")
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 解析配置文件
	tomlConfig, err := config.UnmarshalConfig(tomlFile)
	if err != nil {
		log.Fatalf("UnmarshalConfig: err:%v\n", err)
		return
	}

	// 绑定路由，及公共的tomlConfig
	router = gin.Default()
	router.Use(gin.Recovery())

	router.Use(middlewares.Config(tomlConfig))
	router.Use(middlewares.Gorm("qipaidb", tomlConfig)) // Postgresql数据库服务

	// 路由配置
	router.GET("/", handlers.Index) // 首页
	router.HEAD("/", handlers.Index)
	{
		v1 := router.Group("/v1")
		v1.Use(middlewares.ValidateToken())
		{
			v1.POST("/pay/list/add", handlers.OnPayListAddRequest)       //充值列表-添加
			v1.POST("/pay/list/list", handlers.OnPayListListRequest)     //充值列表-列表
			v1.POST("/pay/list/update", handlers.OnPayListUpdateRequest) //充值列表-修改
			v1.POST("/pay/list/delete", handlers.OnPayListDeleteRequest) //充值列表-删除
			v1.POST("/pay/list/copy", handlers.OnPayListCopyRequest)     //充值列表-删除
		}
	}
	token, _ := middlewares.MakeToken(393)
	unitTest.AddHeader("x-xq5-jwt", token)
	unitTest.SetRouter(router)
	myLog := log.New(os.Stdout, "", log.Lshortfile|log.Ltime)
	unitTest.SetLog(myLog)

	clubdb = ClubdbConfigTest.Connect()
}
