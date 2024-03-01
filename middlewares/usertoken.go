package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JwtMiddleware() gin.HandlerFunc {
	jwtSecret := "aa3cd62b7b442634bc91c8df818ee5633152154c9208f2301ccd6b19f0e8b675"

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenstring := strings.Split(authHeader, " ")[1]
		token,err := jwt.Parse(tokenstring,func (token *jwt.Token) (interface{},error){
			return []byte(jwtSecret),nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from token"})
			c.Abort()
			return
		}

		userId,ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from token"})
			c.Abort()
			return
		}
		c.Set("user_id", userId)
		c.Next()
	}
}