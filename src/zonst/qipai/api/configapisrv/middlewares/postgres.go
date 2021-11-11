package middlewares

import (
	"fmt"

	"zonst/qipai/api/configapisrv/config"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Postgres 连接Postgres数据库
func Postgres(dbName string, tomlConfig *config.Config) gin.HandlerFunc {
	// 读取配置
	dbConfig, ok := tomlConfig.DBServerConf(dbName)
	if !ok {
		panic(fmt.Sprintf("Postgres: %v no set.", dbName))
	}

	// 链接数据库
	db, err := sqlx.Open("postgres", dbConfig.ConnectString())
	if err != nil {
		panic(fmt.Sprintf("sqlx.Open: err:%v", err))
	}

	return func(c *gin.Context) {
		c.Set(dbName, db)
		c.Next()
	}
}

// PostgresUseConfigName 连接Postgres数据库
func PostgresUseConfigName(dbName, configName string, tomlConfig *config.Config) gin.HandlerFunc {
	// 读取配置
	dbConfig, ok := tomlConfig.DBServerConf(configName)
	if !ok {
		panic(fmt.Sprintf("Postgres: %v no set.", configName))
	}

	// 链接数据库
	db, err := sqlx.Open("postgres", dbConfig.ConnectString())
	if err != nil {
		panic(fmt.Sprintf("sqlx.Open: err:%v", err))
	}

	return func(c *gin.Context) {
		c.Set(dbName, db)
		c.Next()
	}
}
