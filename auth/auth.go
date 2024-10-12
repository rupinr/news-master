package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Token struct {
	Value string
}

func (token *Token) valid() bool {
	return token.Value == os.Getenv("ADMIN_TOKEN")
}

func (token *Token) validSubscriberToken() bool {
	return true
}

func ValidateAdminToken(token Token) bool {
	return token.valid()
}

func ValidateSubscriberToken(token Token) bool {
	return token.validSubscriberToken()
}

func AuthMiddleware(validateToken func(Token) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := Token{Value: c.Request.Header.Get("Authorization")}
		if token.Value == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		if !validateToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
