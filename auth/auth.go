package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Token struct {
	Value string
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := Token{Value: c.Request.Header.Get("Authorization")}
		if token.Value == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		if !token.valid() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (token *Token) valid() bool {
	return token.Value == os.Getenv("ADMIN_TOKEN")
}
