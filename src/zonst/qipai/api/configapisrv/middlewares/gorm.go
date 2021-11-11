package middlewares

import (
	"fmt"
	"github.com/jinzhu/gorm"

	"zonst/qipai-sports/api/configapisrv/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Gorm Gorm
func Gorm(dbName string, tomlConfig *config.Config) gin.HandlerFunc {
	// 读取配置
	dbConfig, ok := tomlConfig.DBServerConf(dbName)
	if !ok {
		panic(fmt.Sprintf("Postgres: %v no set.", dbName))
	}

	db, err := gorm.Open("postgres", dbConfig.ConnectString())
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("gorm.Open: err:%v", err))
	}

	// 设置最大链接数
	db.DB().SetMaxOpenConns(10)

	return func(c *gin.Context) {
		c.Set(dbName, db)
		c.Next()
	}
}
