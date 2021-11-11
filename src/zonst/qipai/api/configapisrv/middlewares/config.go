package middlewares

import (
	"zonst/qipai-sports/api/configapisrv/config"
	"zonst/qipai-sports/api/configapisrv/constants"

	"github.com/gin-gonic/gin"
)

func Config(tomlConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constants.Config, tomlConfig)
		c.Next()
	}
}
