package middlewares

import (
	"github.com/fwhezfwhez/errorx"
	"github.com/gin-gonic/gin"
	"time"
	"zonst/qipai/api/configapisrv/dependency/errs"
)

// gin timeout report
func AlertTimeout(c *gin.Context) {
	var timeoutSecond float64 = 10

	start := time.Now()
	c.Next()

	sub := time.Now().Sub(start).Seconds()

	if sub > timeoutSecond {
		errs.SaveError(errorx.NewFromString("触发了一次超时"), map[string]interface{}{
			"fullpath": c.FullPath(),
		})
	}
}
