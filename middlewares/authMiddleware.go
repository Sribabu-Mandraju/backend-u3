package middleware

import (
	helper "backend/helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("error one"),
			})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error  occoureded ",
			})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("company", claims.Company)
		c.Set("user_id", claims.User_id)
		c.Set("contact", claims.Contact)
		c.Next()
	}
}
