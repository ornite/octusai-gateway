package middleware

import (
	"gateway/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
    return func(c *gin.Context) {
		// Retrieve the Authorization header from the incoming request
        authKey := c.GetHeader("authorization")
		claims,err := utils.VerifyToken(authKey)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "StatusUnauthorized",
				"message": err.Error(),
			}})
			c.Abort()
			return
		}

		c.Set("userId",claims.UserID)
        c.Next()
    }
}