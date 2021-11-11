package middlewares

import (
	"fmt"
	"net/http"
	"zonst/logging"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	SecretKey = "88mkva4N{3m47/sm"
)

// Claims is JWT schema of the data it will store
type Claims struct {
	UserID      int32  `json:"user_id"`      // 工号 zonst_user_id
	UserName    string `json:"user_name"`    // 用户姓名
	GroupID     int32  `json:"group_id"`     // 部门ID
	IsSuperUser bool   `json:"is_superuser"` // 是否是超管
	jwt.StandardClaims
}

// GetToken 获取 validate token
func GetToken(ctx *gin.Context) Claims {
	c, exist := ctx.Get("token")
	if exist {
		return c.(Claims)
	}
	return c.(Claims)
}

// ValidateToken 验证 token
func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.Request.Header.Get("x-xq5-jwt")
		// logging.Debugf("jwtToken: %v \n", jwtToken)
		token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte(SecretKey), nil
		})
		if err != nil {
			logging.Errorln(err)
			c.JSON(http.StatusUnauthorized, gin.H{"result": "unauthorized"})
			c.Abort()
			return
		} else if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("token", *claims)
			// before request
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"result": "unauthorized"})
			c.Abort()
			return
		}
	}
}

// MakeToken 生成token
func MakeToken(userID int32) (string, error) {
	claimObj := &Claims{
		UserID:   userID,
		UserName: "默认",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimObj)
	ss, err := token.SignedString([]byte(SecretKey))
	return ss, err
}
