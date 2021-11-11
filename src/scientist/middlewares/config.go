package middlewares

import (
	"zonst/qipai/api/newdzpapisrv/config"

	"github.com/gin-gonic/gin"
)

func Config(tomlConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", tomlConfig)
		c.Next()
	}
}
