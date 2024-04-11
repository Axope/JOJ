package middleware

import (
	"github.com/Axope/JOJ/common/jwt"
	"github.com/Axope/JOJ/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuth(checkAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		j := jwt.GetJWT()

		token, err := j.GetToken(c)
		if err != nil || token == "" {
			c.JSON(http.StatusOK, response.FailMsg("no authority"))
			c.Abort()
			return
		}

		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, response.FailMsg(err.Error()))
			c.Abort()
			return
		}

		if checkAdmin && claims.Admin == 0 {
			c.JSON(http.StatusOK, response.FailMsg("requires administrator privileges"))
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
