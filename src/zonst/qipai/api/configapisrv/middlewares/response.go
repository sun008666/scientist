package middlewares

import "github.com/gin-gonic/gin"

// AddResponseHeader
func AddResponseHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Pragma", "no-cache")
		c.Next()
	}
}
